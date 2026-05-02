package handlers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/jrecasens95/link-nest/backend/internal/auth"
	"github.com/jrecasens95/link-nest/backend/internal/models"
)

type AuthHandler struct {
	service    *auth.Service
	inviteOnly bool
	apiKey     string
}

type authRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string       `json:"token"`
	User  userResponse `json:"user"`
}

type userResponse struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

func NewAuthHandler(service *auth.Service, inviteOnly bool, apiKey string) *AuthHandler {
	return &AuthHandler{service: service, inviteOnly: inviteOnly, apiKey: apiKey}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	if h.inviteOnly && (h.apiKey == "" || c.Get("X-API-Key") != h.apiKey) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": auth.ErrInviteRequired.Error()})
	}

	var payload authRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, token, err := h.service.Register(payload.Name, payload.Email, payload.Password)
	if err != nil {
		return h.authError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(authResponse{Token: token, User: toUserResponse(*user)})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var payload authRequest
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
	}

	user, token, err := h.service.Login(payload.Email, payload.Password)
	if err != nil {
		return h.authError(c, err)
	}

	return c.JSON(authResponse{Token: token, User: toUserResponse(*user)})
}

func (h *AuthHandler) Me(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	user, err := h.service.GetUser(userID)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	return c.JSON(toUserResponse(*user))
}

func (h *AuthHandler) authError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, auth.ErrEmailTaken):
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	case errors.Is(err, auth.ErrInvalidCredentials):
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	default:
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "authentication failed"})
	}
}

func toUserResponse(user models.User) userResponse {
	return userResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
