package rest

import (
	"net/http"
	"strconv"

	"github.com/alvarezcarlos/payment/app/interface/rest/middelware"

	"github.com/alvarezcarlos/payment/app/application"
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/interface/rest/models"
	"github.com/alvarezcarlos/payment/app/interface/rest/validation"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PaymentController struct {
	useCase         application.PaymentUseCaseInterface
	customValidator validation.Validator
	middleware      middelware.Middleware
}

func NewPaymentController(e *echo.Echo, useCase application.PaymentUseCaseInterface,
	customValidator validation.Validator,
	middleware middelware.Middleware) *PaymentController {
	g := e.Group("/api/payments")
	p := &PaymentController{useCase: useCase, customValidator: customValidator}
	g.POST("/create", p.Create, middleware.JwtMiddleware)
	g.GET("/:id", p.GetByID)
	g.POST("/process", p.Process)
	g.POST("/refund", p.Refund, middleware.JwtMiddleware)
	return p
}
func (p *PaymentController) Create(c echo.Context) error {
	pay := models.PaymentCreateReq{}
	if err := c.Bind(&pay); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	merchId, err := getMerchantAttrFromToken(c.Request().Header.Get("Authorization"), merchantIdAttr)
	if err != nil {
		return err
	}

	if strconv.Itoa(int(pay.MerchantID)) != merchId {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid merchant"})
	}

	if err := p.customValidator.ValidateStruct(pay); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	payment := &entity.Payment{
		MerchantID: pay.MerchantID,
		Amount:     pay.Amount,
	}

	payment, err = p.useCase.Create(payment)
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
	return c.JSON(http.StatusOK, map[string]interface{}{"id": payment.ID})
}

func (p *PaymentController) Refund(c echo.Context) error {
	refundReq := models.RefundPaymentReq{}
	if err := c.Bind(&refundReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := p.customValidator.ValidateStruct(refundReq); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	uid, _ := uuid.Parse(refundReq.PaymentID)

	merchId, err := getMerchantAttrFromToken(c.Request().Header.Get("Authorization"), merchantIdAttr)
	if err != nil {
		return err
	}

	err = p.useCase.ProcessRefund(uid, merchId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"id": refundReq.PaymentID})
}
