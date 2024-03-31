package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/alvarezcarlos/payment/app/config"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	"github.com/alvarezcarlos/payment/app/application"
	"github.com/alvarezcarlos/payment/app/domain/entity"
	"github.com/alvarezcarlos/payment/app/interface/rest/models"
	"github.com/alvarezcarlos/payment/app/interface/rest/validation"
	"github.com/labstack/echo/v4"
)

const (
	merchantNameAttr = "merchantName"
	merchantIdAttr   = "merchantId"
)

type MerchantController struct {
	useCase         application.MerchantUseCaseInterface
	customValidator validation.Validator
}

func NewMerchantController(e *echo.Echo, useCase application.MerchantUseCaseInterface, customValidator validation.Validator) *MerchantController {
	g := e.Group("/api/merchants")
	m := &MerchantController{useCase: useCase, customValidator: customValidator}
	g.POST("/create", m.Create)
	g.POST("/login", m.Login)
	g.GET("/details", m.GetDetails)
	return m
}

func (m *MerchantController) GetDetails(c echo.Context) error {
	merchantName, err := getMerchantAttrFromToken(c.Request().Header.Get("Authorization"), merchantNameAttr)
	if err != nil {
		return err
	}

	merchant, err := m.useCase.GetByName(merchantName)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	return c.JSON(http.StatusOK, merchant)
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

	password, err := hashPassword(merch.Password)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	em := entity.Merchant{
		Name:     merch.Name,
		Password: password,
	}

	merchResult, err := m.useCase.Create(&em)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, merchResult)
}

func (m *MerchantController) Login(c echo.Context) error {
	merch := models.Merchant{}
	if err := c.Bind(&merch); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	if err := m.customValidator.ValidateStruct(merch); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
	}

	merchant, err := m.useCase.GetByName(merch.Name)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(merchant.Password), []byte(merch.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, "Unauthorized")
	}

	token, err := generateToken(merchant)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"token": token})
}

func hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

func generateToken(merchant *entity.Merchant) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims[merchantIdAttr] = strconv.Itoa(int(merchant.ID))
	claims[merchantNameAttr] = merchant.Name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	secretKey := config.Config().SecretKey
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func getMerchantAttrFromToken(tokenString, claimAttr string) (string, error) {
	if tokenString == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Token not provided")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config().SecretKey), nil
	})
	if err != nil {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	if !token.Valid {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", echo.NewHTTPError(http.StatusInternalServerError, "Failed to parse claims")
	}
	return claims[claimAttr].(string), nil
}
