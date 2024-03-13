package handlers

import (
	"errors"
	"shopifyx/api/responses"
	"shopifyx/db/entity"
	"shopifyx/db/functions"
	"shopifyx/internal/utils"

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
	usr := entity.User{
		Username: req.Username,
		Name:     req.Name,
		Password: req.Password,
	}

	// Register user
	result, err := u.Database.Register(ctx.UserContext(), usr)
	if err != nil {
		status, response := responses.ErrorServers(err.Error())
		return ctx.Status(status).JSON(response)
	}

	if err.Error() == "EXISTING_USERNAME" {
		status, response := responses.ErrorConflict(err.Error())
		return ctx.Status(status).JSON(response)
	}

	// generate access token
	accessToken, err := utils.GenerateAccessToken(result.Username)
	if err != nil {
		status, response := responses.ErrorServers(err.Error())
		return ctx.Status(status).JSON(response)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User registered successfully",
		"data": fiber.Map{
			"name":        result.Name,
			"username":    result.Username,
			"accessToken": accessToken,
		},
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

	// Validate request body
	if req.Username == "" || req.Password == "" {
		return errors.New("username and password are required")
	}

	// login user
	result, err := u.Database.Login(ctx.UserContext(), req.Username, req.Password)
	if err != nil {
		status, response := responses.ErrorServers(err.Error())
		return ctx.Status(status).JSON(response)
	}

	if err.Error() == "USER_NOT_FOUND" {
		status, response := responses.ErrorNotFound(err.Error())
		return ctx.Status(status).JSON(response)
	}

	if err.Error() == "INVALID_PASSWORD" {
		status, response := responses.ErrorBadRequests(err.Error())
		return ctx.Status(status).JSON(response)
	}

	// generate access token
	accessToken, err := utils.GenerateAccessToken(result.Username)
	if err != nil {
		status, response := responses.ErrorServers(err.Error())
		return ctx.Status(status).JSON(response)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User logged successfully",
		"data": fiber.Map{
			"name":        result.Name,
			"username":    result.Username,
			"accessToken": accessToken,
		},
	})
}
