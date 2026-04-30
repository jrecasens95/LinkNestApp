package handlers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jrecasens95/link-nest/backend/internal/models"
	"github.com/jrecasens95/link-nest/backend/internal/services"
)

type LinkHandler struct {
	baseURL string
	service *services.LinkService
}

type createLinkRequest struct {
	OriginalURL string  `json:"original_url"`
	Title       *string `json:"title"`
}

type updateLinkRequest struct {
	Title    *string `json:"title"`
	IsActive *bool   `json:"is_active"`
}

type linkResponse struct {
	ID          uint    `json:"id"`
	Code        string  `json:"code"`
	OriginalURL string  `json:"original_url"`
	Title       *string `json:"title,omitempty"`
	ClicksCount uint    `json:"clicks_count"`
	IsActive    bool    `json:"is_active"`
	ShortURL    string  `json:"short_url"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

func NewLinkHandler(baseURL string, service *services.LinkService) *LinkHandler {
	return &LinkHandler{
		baseURL: baseURL,
		service: service,
	}
}

func (h *LinkHandler) Create(c *fiber.Ctx) error {
	var payload createLinkRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	payload.OriginalURL = strings.TrimSpace(payload.OriginalURL)
	if !isHTTPURL(payload.OriginalURL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "original_url must start with http:// or https://",
		})
	}

	if payload.Title != nil {
		title := strings.TrimSpace(*payload.Title)
		if title == "" {
			payload.Title = nil
		} else {
			payload.Title = &title
		}
	}

	link, err := h.service.Create(payload.OriginalURL, payload.Title)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not create short link",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"code":      link.Code,
		"short_url": h.baseURL + "/" + link.Code,
	})
}

func (h *LinkHandler) List(c *fiber.Ctx) error {
	links, err := h.service.List()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not list links"})
	}

	response := make([]linkResponse, 0, len(links))
	for _, link := range links {
		response = append(response, h.toLinkResponse(link))
	}

	return c.JSON(fiber.Map{"links": response})
}

func (h *LinkHandler) Get(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid link id"})
	}

	link, err := h.service.GetByID(id)
	if err != nil {
		return h.handleLinkError(c, err, "could not get link")
	}

	return c.JSON(h.toLinkResponse(*link))
}

func (h *LinkHandler) Update(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid link id"})
	}

	var payload updateLinkRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	var titleValue *string
	var titleUpdate **string
	if payload.Title != nil {
		title := strings.TrimSpace(*payload.Title)
		if title != "" {
			titleValue = &title
		}
		titleUpdate = &titleValue
	}

	link, err := h.service.Update(id, services.UpdateLinkInput{
		Title:    titleUpdate,
		IsActive: payload.IsActive,
	})
	if err != nil {
		return h.handleLinkError(c, err, "could not update link")
	}

	return c.JSON(h.toLinkResponse(*link))
}

func (h *LinkHandler) Delete(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid link id"})
	}

	if err := h.service.Delete(id); err != nil {
		return h.handleLinkError(c, err, "could not delete link")
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *LinkHandler) Redirect(c *fiber.Ctx) error {
	code := c.Params("code")
	link, err := h.service.Resolve(code)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrLinkNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "link not found"})
		case errors.Is(err, services.ErrLinkInactive), errors.Is(err, services.ErrLinkExpired):
			return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "link unavailable"})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not resolve link"})
		}
	}

	return c.Redirect(link.OriginalURL, fiber.StatusFound)
}

func (h *LinkHandler) handleLinkError(c *fiber.Ctx, err error, fallback string) error {
	if errors.Is(err, services.ErrLinkNotFound) {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "link not found"})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": fallback})
}

func (h *LinkHandler) toLinkResponse(link models.ShortLink) linkResponse {
	return linkResponse{
		ID:          link.ID,
		Code:        link.Code,
		OriginalURL: link.OriginalURL,
		Title:       link.Title,
		ClicksCount: link.ClicksCount,
		IsActive:    link.IsActive,
		ShortURL:    h.baseURL + "/" + link.Code,
		CreatedAt:   link.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   link.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func parseID(c *fiber.Ctx) (uint, error) {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return 0, errors.New("invalid id")
	}

	return uint(id), nil
}

func isHTTPURL(value string) bool {
	return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
}
