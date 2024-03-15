package handlers

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"shopifyx/api/responses"
	"shopifyx/db/entity"
	"shopifyx/db/functions"
	"strconv"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/gofiber/fiber/v2"
)

type (
	Product struct {
		Database     *functions.Product
		UserDatabase *functions.User
	}

	ProductPayload struct {
		Name           string   `json:"name"`
		Price          int      `json:"price"`
		ImageURL       string   `json:"imageUrl"`
		Stock          int      `json:"stock"`
		Condition      string   `json:"condition"`
		Tags           []string `json:"tags"`
		IsPurchaseable bool     `json:"isPurchaseable"`
	}

	ProductResponse struct {
		ID             string   `json:"id"`
		UserID         string   `json:"userId"`
		Name           string   `json:"name"`
		Price          int      `json:"price"`
		ImageURL       string   `json:"imageUrl"`
		Stock          int      `json:"stock"`
		Condition      string   `json:"condition"`
		Tags           []string `json:"tags"`
		IsPurchaseable bool     `json:"isPurchaseable"`
	}
)

func (app ProductPayload) Validate() error {
	return validation.ValidateStruct(&app,
		// Name cannot be empty, and the length must be between 5 and 60.
		validation.Field(&app.Name, validation.Required, validation.Length(5, 60)),
		// Price cannot be empty, and should be greater than 0.
		validation.Field(&app.Price, validation.NotNil, validation.Min(0)),
		// ImageURL cannot be empty and should be in a valid URL format.
		validation.Field(&app.ImageURL, validation.Required, is.URL),
		// Stock cannot be empty, and should be greater than 0.
		validation.Field(&app.Stock, validation.NotNil, validation.Min(0)),
		// Condition cannot be empty, and should be either "new" or "second".
		validation.Field(&app.Condition, validation.Required, validation.In("new", "second")),
		// Tags cannot be empty, and should have at least 0 items.
		validation.Field(&app.Tags, validation.Required),
		// IsPurchaseable cannot be empty.
		validation.Field(&app.IsPurchaseable, validation.NotNil),
	)
}

func (p *Product) BuyProduct(c *fiber.Ctx) error {
	var payload struct {
		BankAccountId        string `json:"bankAccountId"`
		PaymentProofImageUrl string `json:"paymentProofImageUrl"`
		Qty                  int    `json:"quantity"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(fmt.Sprintf("failed parse payload: %v", err.Error()))
	}

	if c.Params("id") == "" {
		return c.
			Status(http.StatusBadRequest).
			JSON("product id is is required")
	}
	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON("failed parse productId")
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

	bankAccountId, err := strconv.Atoi(payload.BankAccountId)
	if err != nil {
		return c.
			Status(http.StatusBadRequest).
			JSON("failed parse bankAccountId")
	}

	payment, err := p.Database.Buy(c.UserContext(), entity.Payment{
		ProductId:            productID,
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

func (p *Product) AddProduct(c *fiber.Ctx) error {
	userIDClaim := c.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return p.handleError(c, errors.New(fmt.Sprintf("failed parse user id: %v", err.Error())))
	}

	_, err = p.UserDatabase.GetUserById(c.UserContext(), userIDClaim)
	if err != nil {
		return p.handleError(c, fiber.ErrUnauthorized)
	}

	var payload ProductPayload
	if err := c.BodyParser(&payload); err != nil {
		return p.handleError(c, fmt.Errorf("failed parse payload: %w", err))
	}

	err = payload.Validate()
	if err != nil {
		return p.handleError(c, err)
	}

	product, err := p.Database.Add(c.UserContext(), entity.Product{
		UserID:         userID,
		Name:           payload.Name,
		Price:          payload.Price,
		ImageUrl:       payload.ImageURL,
		Stock:          payload.Stock,
		Condition:      payload.Condition,
		Tags:           payload.Tags,
		IsPurchaseable: payload.IsPurchaseable,
	})

	if err != nil {
		return p.handleError(c, err)
	}

	result := p.convertProductEntityToResponse(product)

	return c.Status(http.StatusCreated).JSON(map[string]interface{}{
		"message": "product created successfully",
		"data":    result,
	})

}

func (p *Product) UpdateProduct(c *fiber.Ctx) error {
	userIDClaim := c.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return p.handleError(c, errors.New(fmt.Sprintf("failed parse user id: %v", err.Error())))
	}

	_, err = p.UserDatabase.GetUserById(c.UserContext(), userIDClaim)
	if err != nil {
		return p.handleError(c, fiber.ErrUnauthorized)
	}

	var payload ProductPayload
	if err := c.BodyParser(&payload); err != nil {
		return p.handleError(c, errors.New(fmt.Sprintf("failed parse payload: %v", err.Error())))
	}

	err = payload.Validate()
	if err != nil {
		return p.handleError(c, err)
	}

	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return p.handleError(c, errors.New("failed parse product id"))
	}

	productData, err := p.Database.FindByID(c.UserContext(), productID)
	if err != nil {
		return p.handleError(c, err)
	}

	if productData.UserID != userID {
		return p.handleError(c, fiber.ErrForbidden)
	}

	product, err := p.Database.Update(c.UserContext(), entity.Product{
		ID:             productID,
		UserID:         userID,
		Name:           payload.Name,
		Price:          payload.Price,
		ImageUrl:       payload.ImageURL,
		Stock:          payload.Stock,
		Condition:      payload.Condition,
		Tags:           payload.Tags,
		IsPurchaseable: payload.IsPurchaseable,
	})

	if err != nil {
		return p.handleError(c, err)
	}

	result := p.convertProductEntityToResponse(product)

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"message": "product updated successfully",
		"data":    result,
	})
}

func (p *Product) convertProductEntityToResponse(product entity.Product) ProductResponse {
	return ProductResponse{
		ID:             strconv.Itoa(product.ID),
		UserID:         strconv.Itoa(product.UserID),
		Name:           product.Name,
		Price:          product.Price,
		ImageURL:       product.ImageUrl,
		Stock:          product.Stock,
		Condition:      product.Condition,
		Tags:           product.Tags,
		IsPurchaseable: product.IsPurchaseable,
	}
}

func (p *Product) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, functions.ErrProductNameDuplicate),
		strings.Contains(err.Error(), "failed parse payload"),
		strings.Contains(err.Error(), "failed parse product id"):
		status, response := responses.ErrorBadRequests(err.Error())
		return c.Status(status).JSON(response)
	case errors.Is(err, fiber.ErrUnauthorized):
		return fiber.ErrUnauthorized
	case errors.Is(err, fiber.ErrForbidden):
		return fiber.ErrForbidden
	case errors.Is(err, functions.ErrNoRow):
		status, response := responses.ErrorNotFound("no product found")
		return c.Status(status).JSON(response)
	default:
		validationErrors, ok := err.(validation.Errors)
		if !ok {
			status, response := responses.ErrorServer(err.Error())
			return c.Status(status).JSON(response)
		}

		errMessages := []string{}
		for key, ve := range validationErrors {
			errMessages = append(errMessages, fmt.Sprintf(
				"field %s: %s",
				key,
				ve.Error()))
		}

		status, response := responses.ErrorBadRequests(strings.Join(errMessages, ""))
		return c.Status(status).JSON(response)
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid request body")
	}

	// Validate stock value
	if requestBody.Stock < 0 {
		return c.Status(http.StatusBadRequest).SendString("Stock must be greater than or equal to 0")
	}

	// Retrieve product ID from request parameters
	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid request body")
	}
	dataUpdate := entity.Product{
		ID:    productID,
		Stock: requestBody.Stock,
	}

	// Call UpdateStock method of the database
	product, err := p.Database.UpdateStock(c.UserContext(), dataUpdate, userID)

	if err != nil {
		if err.Error() == "data not found" {
			return c.Status(http.StatusNotFound).SendString(err.Error())
		}
		return c.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"message": "stock updated successfully",
		"data":    product,
	})
}
