package rest

import (
	"github.com/alvarezcarlos/payment/app/application"
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/interface/rest/models"
	"github.com/alvarezcarlos/payment/app/interface/rest/validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

type MerchantController struct {
	useCase         application.MerchantUseCaseInterface
	customValidator validation.Validator
}
type PaymentController struct {
	useCase         application.PaymentUseCaseInterface
	customValidator validation.Validator
}

func NewMerchantController(e *echo.Echo, useCase application.MerchantUseCaseInterface, customValidator validation.Validator) *MerchantController {
	g := e.Group("/api/merchants")
	m := &MerchantController{useCase: useCase, customValidator: customValidator}
	g.POST("/create", m.Create)
	return m
}

func NewPaymentController(e *echo.Echo, useCase application.PaymentUseCaseInterface, customValidator validation.Validator) *PaymentController {
	g := e.Group("/api/payments")
	p := &PaymentController{useCase: useCase, customValidator: customValidator}
	g.POST("/create", p.Create)
	g.GET("/:id", p.GetByID)
	g.POST("/process", p.Process)
	return p
}

func (m *MerchantController) Create(c echo.Context) error {
	merch := models.Merchant{}
	if err := c.Bind(&merch); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := m.customValidator.ValidateStruct(merch); err != nil {
		c.Logger().Error(err)
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

	if err := p.customValidator.ValidateStruct(pay); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	payment := &entity.Payment{
		MerchantID: pay.MerchantID,
		Amount:     pay.Amount,
	}

	payment, err := p.useCase.Create(payment)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, map[string]string{"message": "created", "id": payment.ID.String()})
}

func (p *PaymentController) GetByID(c echo.Context) error {
	id := c.Param("id")
	parsedUUID, err := uuid.Parse(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	payment, err := p.useCase.GetByID(parsedUUID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, payment)
}

func (p *PaymentController) Process(c echo.Context) error {
	processReq := models.ProcessPaymentReq{}
	if err := c.Bind(&processReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := p.customValidator.ValidateStruct(processReq); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}
	id, _ := uuid.Parse(processReq.PaymentID)
	processPay := &entity.Payment{
		ID: id,
	}
	customer := &entity.Card{
		HolderID:   processReq.Customer.PersonalID,
		HolderName: processReq.Customer.Name,
		Number:     processReq.Card.Number,
		Code:       processReq.Card.Code,
		Month:      processReq.Card.Month,
		Year:       processReq.Card.Year,
	}

	payment, err := p.useCase.ProcessPayment(processPay, customer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	//crear payment response paar no devolver datos sensibles
	return c.JSON(http.StatusOK, payment)
}
