package handlers

import (
	"errors"
	"net/http"
	"shopifyx/db/entity"
	"shopifyx/db/functions"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type BankHandler struct {
	Bank functions.Bank
}

func (b *BankHandler) Create(c *fiber.Ctx) error {
	userIDClaim := c.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	var payload struct {
		BankName          string `json:"bankName"`
		BankAccountName   string `json:"bankAccountName"`
		BankAccountNumber string `json:"bankAccountNumber"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	var (
		lenBankName      = len(payload.BankName)
		lenAccountName   = len(payload.BankAccountName)
		lenAccountNumber = len(payload.BankAccountNumber)
		isValid          = func(l int) bool {
			return l >= 5 && l <= 15
		}
	)

	if !isValid(lenBankName) || !isValid(lenAccountName) || !isValid(lenAccountNumber) {
		return c.SendStatus(http.StatusBadRequest)
	}

	if err := b.Bank.Create(c.UserContext(), entity.Bank{
		UserId:            userID,
		BankName:          payload.BankName,
		BankAccountName:   payload.BankAccountName,
		BankAccountNumber: payload.BankAccountNumber,
	}); err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

func (b *BankHandler) Get(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)

	accounts, err := b.Bank.Get(c.UserContext(), userId)
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.JSON(map[string]interface{}{
		"message": "success",
		"data":    accounts,
	})
}

func (b *BankHandler) Delete(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(string)
	bankAccountId := c.Params("bankAccountId")

	err := b.Bank.Delete(c.UserContext(), userId, bankAccountId)
	if err != nil {
		if errors.Is(err, functions.ErrNoRow) {
			return c.SendStatus(http.StatusNotFound)
		}

		if errors.Is(err, functions.ErrUnauthorized) {
			return c.SendStatus(http.StatusUnauthorized)
		}

		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}

func (b *BankHandler) Update(c *fiber.Ctx) error {
	userIDClaim := c.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return c.SendStatus(http.StatusUnauthorized)
	}

	var payload struct {
		BankName          string `json:"bankName"`
		BankAccountName   string `json:"bankAccountName"`
		BankAccountNumber string `json:"bankAccountNumber"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	var (
		lenBankName      = len(payload.BankName)
		lenAccountName   = len(payload.BankAccountName)
		lenAccountNumber = len(payload.BankAccountNumber)
		isValid          = func(l int) bool {
			return l >= 5 && l <= 15
		}
	)

	if !isValid(lenBankName) || !isValid(lenAccountName) || !isValid(lenAccountNumber) {
		return c.SendStatus(http.StatusBadRequest)
	}

	err = b.Bank.Update(c.UserContext(), entity.Bank{
		Id:                c.Params("bankAccountId"),
		UserId:            userID,
		BankName:          payload.BankName,
		BankAccountName:   payload.BankAccountName,
		BankAccountNumber: payload.BankAccountNumber,
	})

	if errors.Is(err, functions.ErrNoRow) {
		return c.SendStatus(http.StatusNotFound)
	}

	if errors.Is(err, functions.ErrUnauthorized) {
		return c.SendStatus(http.StatusUnauthorized)
	}

	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	return c.SendStatus(http.StatusOK)
}
