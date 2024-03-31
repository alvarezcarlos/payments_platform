# Payment Platform Service Endpoints Documentation

## Introduction
Welcome to the Payment Platform Service Endpoints Documentation! This document provides detailed instructions on how to interact with the various endpoints exposed by the Payment Platform Service.

## Create Merchant Endpoint

### Description
This endpoint is used to create a new merchant in the system.

### Endpoint
```bash
curl --request POST \
  --url http://localhost:8080/api/merchants/create \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/8.6.1' \
  --data '{
    "name": "enterprise",
    "password": "12345"
  }'
```
### Example Response
```json
{
	"id": 1,
	"name": "enterpice",
	"balance": 11089.54,
	"payments": null,
	"CreatedAt": "2024-03-31T11:33:37.474472-03:00",
	"UpdatedAt": "2024-03-31T11:33:37.474472-03:00"
}
```
# Merchant Login Endpoint

## Description
This endpoint is used for merchants to log in to the system.

## Endpoint
```bash
curl --request POST \
  --url http://localhost:8080/api/merchants/login \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/8.6.1' \
  --data '{
	"name": "enterprise",
	"password": "1234"
}'
```
### Example Response
```json
{
	"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE5ODI1MTAsIm1lcmNoYW50SWQiOiIyIiwibWVyY2hhbnROYW1lIjoiZW50ZXJwaWNlMiJ9.FAHBkYJOfh4z2Mavk-Sfpn8Z_76oQFx91hXd7SX8F5E"
}
```
# Create Payment Endpoint

## Description
This endpoint is used to create a new payment transaction.

## Endpoint
```bash
curl --request POST \
  --url http://localhost:8080/api/payments/create \
  --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE5ODI1MTAsIm1lcmNoYW50SWQiOiIyIiwibWVyY2hhbnROYW1lIjoiZW50ZXJwaWNlMiJ9.FAHBkYJOfh4z2Mavk-Sfpn8Z_76oQFx91hXd7SX8F5E' \
  --header 'Content-Type: application/json' \
  --data '{
	"merchant_id": 2,
	"amount": 1000.00
}'
```
### Example Response
```json
{
	"id": "8f724474-1cc0-43ac-aa5d-2ffe2edc1e81",
	"message": "created"
}
```
# Process Payment Endpoint

## Description
This endpoint allows customer to complete a payment transaction.

## Endpoint
```bash
curl --request POST \
  --url http://localhost:8080/api/payments/process \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/8.6.1' \
  --data '{
    "payment_id": "8f724474-1cc0-43ac-aa5d-2ffe2edc1e81",
    "card": {
        "number": "1111111111111111",
        "code": "111",
        "month": 12,
        "year": 28
    },
    "customer": {
        "personal_id": 111111111,
        "name": "Jhon Doe"
    }
}'
```
### Example Response
```json
{
	"id": "8f724474-1cc0-43ac-aa5d-2ffe2edc1e81"
}
```
# Refund Payment Endpoint

## Description
This endpoint is used to process a refund for a payment transaction. Also is used by merchant as the creation endpoint.

## Endpoint
```bash
curl --request POST \
  --url http://localhost:8080/api/payments/refund \
  --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE5ODE4MzQsIm1lcmNoYW50SWQiOiIxIiwibWVyY2hhbnROYW1lIjoiZW50ZXJwaWNlIn0.mkAfyLPGD_nnvNmI2wbJ60O3o-hHO7eTialfiyK6PAE' \
  --header 'Content-Type: application/json' \
  --header 'User-Agent: insomnia/8.6.1' \
  --data '{
	"payment_id": "52860a9d-b3b3-4d92-a1ed-357196edad85"
}'
```
### Example Response
```json
{
	"id": "8f724474-1cc0-43ac-aa5d-2ffe2edc1e81"
}
```
# Merchant Details Endpoint

## Description
This endpoint is used to retrieve details of a merchant.

## Endpoint
```bash
curl --request GET \
  --url http://localhost:8080/api/merchants/details \
  --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE5NzcyNTMsIm1lcmNoYW50TmFtZSI6ImVudGVycGljZSJ9.4kPT66HD54BNZ6HtMGn6c_0J0IWIWZ7RalewjQuV-68' \
  --header 'User-Agent: insomnia/8.6.1'
```
### Example Response
```json
{
	"id": 1,
	"name": "enterpice",
	"balance": 94046.62,
	"payments": [
		{
			"id": "52860a9d-b3b3-4d92-a1ed-357196edad85",
			"amount": 1000,
			"merchant_id": 1,
			"states": null,
			"created_at": "2024-03-31T11:31:02.068368-03:00",
			"updated_at": "2024-03-31T11:32:40.365093-03:00"
		}
	],
	"CreatedAt": "2024-03-31T11:29:38.597164-03:00",
	"UpdatedAt": "2024-03-31T11:32:40.361271-03:00"
}

```
# Payment Details Endpoint

## Description
This endpoint is used to retrieve details of a specific payment transaction.

## Endpoint
```bash
curl --request GET \
  --url http://localhost:8080/api/payments/8f724474-1cc0-43ac-aa5d-2ffe2edc1e81 \
  --header 'Authorization: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTE5ODIzODgsIm1lcmNoYW50SWQiOiIyIiwibWVyY2hhbnROYW1lIjoiZW50ZXJwaWNlMiJ9.a8EKFUUmIn4VQJxI9vVUvWxyOP2VrJn_ZQtP3QR1tzk' \
  --header 'User-Agent: insomnia/8.6.1'
```
### Example Response
```json
{
	"id": "8f724474-1cc0-43ac-aa5d-2ffe2edc1e81",
	"amount": 1000,
	"merchant_id": 2,
	"states": [
		{
			"id": 1,
			"name": "Pending"
		},
		{
			"id": 3,
			"name": "Succeeded"
		}
	],
	"created_at": "2024-03-31T11:43:30.955633-03:00",
	"updated_at": "2024-03-31T11:43:58.788663-03:00"
}
```