package handlers

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
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

func isHTTPURL(value string) bool {
	return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://")
}
