package handlers

import (
	"errors"
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jrecasens95/link-nest/backend/internal/models"
	"github.com/jrecasens95/link-nest/backend/internal/security"
	"github.com/jrecasens95/link-nest/backend/internal/services"
)

type LinkHandler struct {
	baseURL   string
	service   *services.LinkService
	validator *security.URLValidator
}

type createLinkRequest struct {
	OriginalURL string  `json:"original_url"`
	Title       *string `json:"title"`
	CustomAlias *string `json:"custom_alias"`
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

type clickEventResponse struct {
	ID        uint   `json:"id"`
	UserAgent string `json:"user_agent"`
	Referer   string `json:"referer"`
	IPAddress string `json:"ip_address"`
	CreatedAt string `json:"created_at"`
}

type refererStatResponse struct {
	Referer string `json:"referer"`
	Count   int64  `json:"count"`
}

type linkStatsResponse struct {
	TotalClicks  uint                  `json:"total_clicks"`
	RecentClicks []clickEventResponse  `json:"recent_clicks"`
	Referers     []refererStatResponse `json:"referers"`
}

func NewLinkHandler(baseURL string, service *services.LinkService, validator *security.URLValidator) *LinkHandler {
	return &LinkHandler{
		baseURL:   baseURL,
		service:   service,
		validator: validator,
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
	if err := h.validator.ValidateCreateURL(payload.OriginalURL); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": validationMessage(err),
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

	if payload.CustomAlias != nil {
		alias := strings.TrimSpace(*payload.CustomAlias)
		if alias == "" {
			payload.CustomAlias = nil
		} else {
			payload.CustomAlias = &alias
		}
	}

	link, err := h.service.Create(payload.OriginalURL, payload.Title, payload.CustomAlias)
	if err != nil {
		if errors.Is(err, services.ErrInvalidAlias) ||
			errors.Is(err, services.ErrReservedAlias) ||
			errors.Is(err, services.ErrAliasExists) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

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

func (h *LinkHandler) Stats(c *fiber.Ctx) error {
	id, err := parseID(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid link id"})
	}

	stats, err := h.service.Stats(id)
	if err != nil {
		return h.handleLinkError(c, err, "could not get link stats")
	}

	recentClicks := make([]clickEventResponse, 0, len(stats.RecentClicks))
	for _, click := range stats.RecentClicks {
		recentClicks = append(recentClicks, clickEventResponse{
			ID:        click.ID,
			UserAgent: click.UserAgent,
			Referer:   click.Referer,
			IPAddress: click.IPAddress,
			CreatedAt: click.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	referers := make([]refererStatResponse, 0, len(stats.Referers))
	for _, referer := range stats.Referers {
		referers = append(referers, refererStatResponse{
			Referer: referer.Referer,
			Count:   referer.Count,
		})
	}

	return c.JSON(linkStatsResponse{
		TotalClicks:  stats.TotalClicks,
		RecentClicks: recentClicks,
		Referers:     referers,
	})
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
	link, err := h.service.Resolve(code, services.ClickInput{
		UserAgent: c.Get("User-Agent"),
		Referer:   c.Get("Referer"),
		IPAddress: anonymizeIP(c.IP()),
	})
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

	if err := h.validator.ValidateRedirectURL(link.OriginalURL); err != nil {
		return c.Status(fiber.StatusGone).JSON(fiber.Map{"error": "link target is no longer allowed"})
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

func anonymizeIP(value string) string {
	ip := net.ParseIP(value)
	if ip == nil {
		return ""
	}

	if ipv4 := ip.To4(); ipv4 != nil {
		return net.IPv4(ipv4[0], ipv4[1], ipv4[2], 0).String()
	}

	ipv6 := ip.To16()
	if ipv6 == nil {
		return ""
	}

	for i := 8; i < len(ipv6); i++ {
		ipv6[i] = 0
	}

	return ipv6.String()
}

func validationMessage(err error) string {
	switch {
	case errors.Is(err, security.ErrURLTooLong):
		return "original_url is too long"
	case errors.Is(err, security.ErrURLUnsupportedScheme):
		return "original_url must start with http:// or https://"
	case errors.Is(err, security.ErrURLBlockedHost), errors.Is(err, security.ErrURLPrivateAddress):
		return "original_url host is not allowed"
	default:
		return "original_url is invalid"
	}
}
