package handlers

import (
	"errors"
	"shopifyx/api/responses"
	"shopifyx/internal/database/functions"
	"shopifyx/internal/database/interfaces"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	Database *functions.User
}

func validateUser(req struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}) error {
	lenUsername := len(req.Username)
	lenPassword := len(req.Password)
	lenName := len(req.Name)

	if lenUsername == 0 || lenPassword == 0 || lenName == 0 {
		return errors.New("username and password are required")
	}

	if lenUsername < 5 || lenPassword < 5 || lenName < 5 {
		return errors.New("username and password length must be at least 5 characters")
	}

	if lenUsername > 15 || lenPassword > 15 || lenName > 15 {
		return errors.New("username and password length cannot exceed 15 characters")
	}

	return nil
}

func (u *User) Register(ctx *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	// Validate request body
	if err := validateUser(req); err != nil {
		status, response := responses.ErrorBadRequests(err.Error())
		return ctx.Status(status).JSON(response)
	}

	// Create user object
	usr := interfaces.User{
		Username: req.Username,
		Name:     req.Name,
		Password: req.Password,
	}

	// Register user
	err := u.Database.Register(ctx.UserContext(), usr)
	if err != nil {
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func (u *User) Login(ctx *fiber.Ctx) error {
	// Parse request body
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := ctx.BodyParser(&req); err != nil {
		return err
	}

	return nil
}
