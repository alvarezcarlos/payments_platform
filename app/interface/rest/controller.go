package rest

import (
	"github.com/alvarezcarlos/payment/app/application"
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/interface/rest/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MerchantController struct {
	useCase application.MerchantUseCaseInterface
}
type PaymentController struct {
	useCase application.PaymentUseCaseInterface
}

func NewMerchantController(e *echo.Echo, useCase application.MerchantUseCaseInterface) *MerchantController {
	g := e.Group("/api/merchants")
	m := &MerchantController{useCase: useCase}
	g.POST("/create", m.Create)
	return m
}

func NewPaymentController(e *echo.Echo, useCase application.PaymentUseCaseInterface) *PaymentController {
	g := e.Group("/api/payments")
	m := &PaymentController{useCase: useCase}
	g.POST("/create", m.Create)
	//g.PUT("/update", m.Update)
	return m
}

func (m *MerchantController) Create(c echo.Context) error {
	merch := models.Merchant{}
	if err := c.Bind(&merch); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	em := entity.Merchant{
		Name:          merch.Name,
		AccountNumber: merch.Account,
	}

	if err := m.useCase.Create(&em); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "created"})
}

func (p *PaymentController) Create(c echo.Context) error {
	pay := models.PaymentCreateReq{}
	if err := c.Bind(&pay); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	payment := &entity.Payment{
		MerchantID: pay.MerchantID,
		Amount:     pay.Amount,
	}

	if err := p.useCase.Create(payment); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "created"})
}
