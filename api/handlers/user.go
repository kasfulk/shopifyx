package handlers

import (
	"errors"
	"shopifyx/internal/database/functions"
	"shopifyx/internal/database/interfaces"

	"github.com/gofiber/fiber/v2"
)

type User struct {
	Database functions.User
}

func validateUser(req struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Password string `json:"password"`
}) error {
	// Check for empty fields
	if req.Username == "" || req.Name == "" || req.Password == "" {
		return errors.New("all fields are required")
	}

	// Check length constraints
	if len(req.Username) < 5 || len(req.Username) > 15 {
		return errors.New("username length must be between 5 and 15 characters")
	}
	if len(req.Name) < 5 || len(req.Name) > 50 {
		return errors.New("name length must be between 5 and 50 characters")
	}
	if len(req.Password) < 5 || len(req.Password) > 15 {
		return errors.New("password length must be between 5 and 15 characters")
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
		return err
	}

	// Create user object
	usr := interfaces.User{
		Username: req.Username,
		Name:     req.Name,
		Password: req.Password,
	}

	// Register user
	err := u.Database.Register(usr)
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
