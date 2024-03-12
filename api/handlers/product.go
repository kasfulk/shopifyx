package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"shopifyx/internal/database/functions"
	"shopifyx/internal/database/interfaces"
	"strconv"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofiber/fiber/v2"
)

type Product struct {
	Database *functions.Product
}

func (p *Product) BuyProduct(c *fiber.Ctx) error {
	var payload struct {
		ProductId            string `json:"productId"`
		BankAccountId        string `json:"bankAccountId"`
		PaymentProofImageUrl string `json:"paymentProofImageUrl"`
		Qty                  int    `json:"quantity"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(fmt.Sprintf("failed parse payload: %v", err.Error()))
	}

	if payload.ProductId == "" || c.Params("id") == "" {
		return c.
			Status(http.StatusBadRequest).
			JSON("product id is is required")
	}

	if c.Params("id") != payload.ProductId {
		return c.
			Status(http.StatusBadRequest).
			JSON("product id from params and body is different")
	}

	if payload.BankAccountId == "" {
		return c.
			Status(http.StatusBadRequest).
			JSON("bank account id is required")
	}

	if validation.Validate(payload.PaymentProofImageUrl, validation.Required, is.URL) != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON("payment proof image url is empty or malformat")
	}

	if payload.Qty < 1 {
		return c.
			Status(http.StatusBadRequest).
			JSON("minimum amount of quantity must be 1")
	}

	productId, err := strconv.Atoi(payload.ProductId)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON("failed parse productId")
	}

	bankAccountId, err := strconv.Atoi(payload.ProductId)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON("failed parse productId")
	}

	payment, err := p.Database.Buy(c.UserContext(), interfaces.Payment{
		ProductId:            productId,
		BankAccountId:        bankAccountId,
		PaymentProofImageUrl: payload.PaymentProofImageUrl,
		Qty:                  payload.Qty,
	})

	if errors.Is(err, functions.ErrNoRow) || errors.Is(err, functions.ErrInsuficientQty) {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if err != nil {
		slog.Error(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(err.Error())
	}

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"message": "payment processed successfully",
		"data":    payment,
	})
}
