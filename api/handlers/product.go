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
		BankDatabase *functions.Bank
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

	QueryFilterGetProducts struct {
		UserOnly       bool     `json:"userOnly"`
		Limit          int      `json:"limit"`
		Offset         int      `json:"offset"`
		Tags           []string `json:"tags"`
		Condition      string   `json:"condition"`
		ShowEmptyStock bool     `json:"showEmptyStock"`
		MaxPrice       int      `json:"maxPrice"`
		MinPrice       int      `json:"minPrice"`
		SortBy         string   `json:"sortBy"`
		OrderBy        string   `json:"orderBy"`
		Search         string   `json:"search"`
	}

	ProductResponse struct {
		ProductId      string   `json:"productId"`
		Name           string   `json:"name"`
		Price          int      `json:"price"`
		ImageUrl       string   `json:"imageUrl"`
		Stock          int      `json:"stock"`
		Condition      string   `json:"condition"`
		Tags           []string `json:"tags"`
		IsPurchaseable bool     `json:"isPurchaseable"`
		PurchaseCount  int      `json:"purchaseCount"`
	}

	Meta struct {
		Limit  int `json:"limit"`
		Offset int `json:"offset"`
		Total  int `json:"total"`
	}

	Bank struct {
		BankAccountId     string `json:"bankAccountId"`
		BankName          string `json:"bankName"`
		BankAccountName   string `json:"bankAccountName"`
		BankAccountNumber string `json:"bankAccountNumber"`
	}

	SellerData struct {
		Name             string `json:"name"`
		ProductSoldTotal int    `json:"productSoldTotal"`
		BankAccounts     []Bank `json:"bankAccounts"`
	}

	GetProductsResponse struct {
		Data []ProductResponse `json:"data"`
		Meta Meta              `json:"meta"`
	}

	GetProductDetailResponse struct {
		Product    ProductResponse `json:"product"`
		SellerData SellerData      `json:"seller"`
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

func (app QueryFilterGetProducts) Validate() error {
	return validation.ValidateStruct(&app,
		// Limit should be greater than 0.
		validation.Field(&app.Limit, validation.Min(0)),
		// Offset cannot should be greater than 0.
		validation.Field(&app.Offset, validation.Min(0)),
		// Condition should be either "new" or "second".
		validation.Field(&app.Condition, validation.In("new", "second")),
		// MaxPrice should be greater than 0.
		validation.Field(&app.MaxPrice, validation.Min(0)),
		// MinPrice should be greater than 0.
		validation.Field(&app.MinPrice, validation.Min(0)),
		// SortBy should be either "price" or "date".
		validation.Field(&app.SortBy, validation.In("price", "date")),
		// OrderBy should be either "asc" or "dsc".
		validation.Field(&app.OrderBy, validation.In("asc", "dsc")),
	)
}

func (p *Product) convertProductEntityToResponse(product entity.Product) ProductResponse {
	return ProductResponse{
		ProductId:      strconv.Itoa(product.ID),
		Name:           product.Name,
		Price:          product.Price,
		ImageUrl:       product.ImageUrl,
		Stock:          product.Stock,
		Condition:      product.Condition,
		Tags:           product.Tags,
		IsPurchaseable: product.IsPurchaseable,
		PurchaseCount:  product.PurchaseCount,
	}
}

func (p *Product) convertQueryFilterToEntity(filter QueryFilterGetProducts) entity.FilterGetProducts {
	return entity.FilterGetProducts{
		UserOnly:       filter.UserOnly,
		Limit:          filter.Limit,
		Offset:         filter.Offset,
		Tags:           filter.Tags,
		Condition:      filter.Condition,
		ShowEmptyStock: filter.ShowEmptyStock,
		MaxPrice:       filter.MaxPrice,
		MinPrice:       filter.MinPrice,
		SortBy:         filter.SortBy,
		OrderBy:        filter.OrderBy,
		Search:         filter.Search,
	}
}

func (p *Product) convertProductsToGetProductsResponse(
	products []entity.Product,
	limit, offset, total int,
) GetProductsResponse {
	var result []ProductResponse
	for _, product := range products {
		result = append(result, p.convertProductEntityToResponse(product))
	}

	return GetProductsResponse{
		Data: result,
		Meta: Meta{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	}
}

func (p *Product) convertProductToProductDetailResponse(
	product entity.Product,
	seller entity.User,
	productSoldTotal int,
	bankAccounts []entity.Bank,
) GetProductDetailResponse {
	bankAccountsResponse := []Bank{}
	for _, bank := range bankAccounts {
		bankAccountsResponse = append(bankAccountsResponse, Bank{
			BankAccountId:     bank.Id,
			BankName:          bank.BankName,
			BankAccountName:   bank.BankAccountName,
			BankAccountNumber: bank.BankAccountNumber,
		})
	}

	return GetProductDetailResponse{
		Product: p.convertProductEntityToResponse(product),
		SellerData: SellerData{
			Name:             seller.Name,
			ProductSoldTotal: productSoldTotal,
			BankAccounts:     bankAccountsResponse,
		},
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
}

func (p *Product) GetProducts(c *fiber.Ctx) error {
	var (
		userID int
		err    error
	)

	var filter QueryFilterGetProducts
	if err := c.QueryParser(&filter); err != nil {
		return p.handleError(c, errors.New(fmt.Sprintf("failed to parse query params: %v", err.Error())))
	}

	err = filter.Validate()
	if err != nil {
		return p.handleError(c, err)
	}

	if c.Locals("user_id") != nil {
		userIDClaim := c.Locals("user_id").(string)
		userID, err = strconv.Atoi(userIDClaim)
		if err != nil {
			return p.handleError(c, errors.New(fmt.Sprintf("failed parse user id: %v", err.Error())))
		}
	}

	filterDB := p.convertQueryFilterToEntity(filter)
	products, err := p.Database.FindAll(c.UserContext(), filterDB, userID)
	if err != nil {
		return p.handleError(c, err)
	}

	total, err := p.Database.Count(c.UserContext(), filterDB, userID)
	if err != nil {
		return p.handleError(c, err)
	}

	result := p.convertProductsToGetProductsResponse(products, filter.Limit, filter.Offset, total)

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"message": "ok",
		"data":    result,
	})
}

func (p *Product) GetProductDetail(c *fiber.Ctx) error {
	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return p.handleError(c, errors.New("failed parse product id"))
	}

	product, err := p.Database.FindByID(c.UserContext(), productID)
	if err != nil {
		return p.handleError(c, err)
	}

	user, err := p.UserDatabase.GetUserById(c.UserContext(), strconv.Itoa(product.UserID))
	if err != nil {
		return p.handleError(c, err)
	}

	productSoldTotal, err := p.Database.SumPurchaseCountByUserID(c.UserContext(), product.UserID)
	if err != nil {
		return p.handleError(c, err)
	}

	bankAccounts, err := p.BankDatabase.Get(c.UserContext(), strconv.Itoa(product.UserID))
	if err != nil {
		return p.handleError(c, err)
	}

	result := p.convertProductToProductDetailResponse(product, user, productSoldTotal, bankAccounts)

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"message": "ok",
		"data":    result,
	})
}

func (p *Product) BuyProduct(c *fiber.Ctx) error {
	var payload struct {
		BankAccountId        string `json:"bankAccountId"`
		PaymentProofImageUrl string `json:"paymentProofImageUrl"`
		Qty                  int    `json:"quantity"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.
			Status(http.StatusBadRequest).
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

	if err != nil {
		if errors.Is(err, functions.ErrNoRow) {
			return c.Status(http.StatusNotFound).JSON(err.Error())
		} else if errors.Is(err, functions.ErrInsuficientQty) {
			return c.Status(http.StatusBadRequest).JSON(err.Error())
		}

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
		return c.SendStatus(http.StatusBadRequest)
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
		return c.SendStatus(http.StatusBadRequest)
	}

	err = payload.Validate()
	if err != nil {
		return p.handleError(c, err)
	}

	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return p.handleError(c, errors.New("failed parse product id"))
	}

	product, err := p.Database.FindByIDUser(c.UserContext(), productID, userID)
	if err != nil {
		if err == functions.ErrNoRow {
			return p.handleError(c, fiber.ErrForbidden)
		}
		return p.handleError(c, err)
	}

	product.Name = payload.Name
	product.Price = payload.Price
	product.ImageUrl = payload.ImageURL
	product.Stock = payload.Stock
	product.Condition = payload.Condition
	product.Tags = payload.Tags
	product.IsPurchaseable = payload.IsPurchaseable

	err = p.Database.Update(c.UserContext(), product)
	if err != nil {
		return p.handleError(c, err)
	}

	result := p.convertProductEntityToResponse(product)

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"message": "product updated successfully",
		"data":    result,
	})
}

func (p *Product) UpdateStock(c *fiber.Ctx) error {
	userIDClaim := c.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return p.handleError(c, err)
	}

	var requestBody struct {
		Stock int `json:"stock"`
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

	productCheck, err := p.Database.FindByIDUser(c.UserContext(), productID, userID)

	if err != nil {
		return c.Status(http.StatusNotFound).SendString("Product not found")
	}

	productCheck.Stock = requestBody.Stock

	// Call UpdateStock method of the database
	product, err := p.Database.UpdateStock(c.UserContext(), productCheck, userID)

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

func (p *Product) DeleteProduct(c *fiber.Ctx) error {
	userIDClaim := c.Locals("user_id").(string)
	userID, err := strconv.Atoi(userIDClaim)
	if err != nil {
		return p.handleError(c, errors.New(fmt.Sprintf("failed parse user id: %v", err.Error())))
	}

	_, err = p.UserDatabase.GetUserById(c.UserContext(), userIDClaim)
	if err != nil {
		return p.handleError(c, fiber.ErrUnauthorized)
	}

	productID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return p.handleError(c, errors.New("failed parse product id"))
	}

	_, err = p.Database.FindByIDUser(c.UserContext(), productID, userID)
	if err != nil {
		if err == functions.ErrNoRow {
			return p.handleError(c, fiber.ErrForbidden)
		}
		return p.handleError(c, err)
	}

	err = p.Database.DeleteByID(c.UserContext(), productID)
	if err != nil {
		return p.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(map[string]interface{}{
		"message": "product deleted successfully",
	})
}
