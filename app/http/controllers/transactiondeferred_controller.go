package controllers

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	webpay "github.com/ppastene/transbank-sdk-go"
)

type TransactionDeferredController struct {
	transaction *webpay.Transaction
}

func NewTransactionDeferredController() *TransactionDeferredController {
	options := &webpay.Options{
		ApiKey:       "597055555540",
		CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	}
	transaction := webpay.NewTransaction(options)
	return &TransactionDeferredController{
		transaction: transaction,
	}
}

func (t *TransactionDeferredController) Index(ctx http.Context) http.Response {
	baseUrl := facades.Config().Env("APP_URL", "").(string)
	port := facades.Config().Env("APP_PORT", "").(string)
	var fullUrl string
	if port != "" {
		fullUrl = fmt.Sprintf("%s:%s", baseUrl, port)
	} else {
		fullUrl = baseUrl
	}
	data := map[string]any{
		"title":      "Webpay Plus Diferido",
		"returnUrl":  fullUrl,
		"isDeferred": true,
	}
	return ctx.Response().View().Make("webpayplus/index.tmpl", data)
}

func (t *TransactionDeferredController) CreatedTransaction(ctx http.Context) http.Response {
	request := ctx.Request().All()
	buyOrder := ctx.Request().Input("buy_order")
	sessionId := ctx.Request().Input("session_id")
	returnUrl := ctx.Request().Input("return_url")
	amountStr := ctx.Request().Input("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	response, err := t.transaction.Create(buyOrder, sessionId, amount, returnUrl)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	data := map[string]any{
		"title":    "Transacción Diferida Creada",
		"request":  request,
		"response": response,
	}

	return ctx.Response().View().Make("webpayplus/transaction_created.tmpl", data)
}

func (t *TransactionDeferredController) CommitedTransaction(ctx http.Context) http.Response {
	// Flujo abortado: Se revisa si en la url TBK_TOKEN existe
	// De ser asi se considera la transaccion como abortada
	// Se retorna a vista
	queryParams := ctx.Request().Queries()
	if queryParams["TBK_TOKEN"] != "" {
		data := map[string]any{
			"title":       "Transacción Diferida Abortada",
			"queryParams": queryParams,
		}
		return ctx.Response().View().Make("webpayplus/malltransaction_aborted.tmpl", data)
	}

	// Si no está presente pues es un flujo normal
	var token string = ctx.Request().Input("token_ws")
	response, err := t.transaction.Commit(token)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	data := map[string]any{
		"title":      "Transacción Diferida Confirmada",
		"request":    ctx.Request().All(),
		"response":   response,
		"isApproved": response.IsApproved(),
	}
	return ctx.Response().View().Make("webpayplus/transaction_committed.tmpl", data)
}

func (t *TransactionDeferredController) GetTransactionStatus(ctx http.Context) http.Response {
	request := ctx.Request().All()
	response, err := t.transaction.Status(ctx.Request().Input("token"))
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	data := map[string]any{
		"title":    "Estado de la Transacción Diferida",
		"request":  request,
		"response": response,
	}
	return ctx.Response().View().Make("webpayplus/status.tmpl", data)
}

func (t *TransactionDeferredController) RefundTransaction(ctx http.Context) http.Response {
	request := ctx.Request().All()
	token := ctx.Request().Input("token")
	amountStr := ctx.Request().Input("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	response, err := t.transaction.Refund(token, amount)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	data := map[string]any{
		"title":    "Transacción Diferida Reembolsada",
		"request":  request,
		"response": response,
	}
	return ctx.Response().View().Make("webpayplus/refund_success.tmpl", data)
}

func (t *TransactionDeferredController) CaptureTransaction(ctx http.Context) http.Response {
	request := ctx.Request().All()
	token := ctx.Request().Input("token")
	buyOrder := ctx.Request().Input("buy_order")
	authCode := ctx.Request().Input("authorization_code")
	amountStr := ctx.Request().Input("capture_amount")
	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	response, err := t.transaction.Capture(token, buyOrder, authCode, amount)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	data := map[string]any{
		"title":    "Transacción Diferida Capturada",
		"request":  request,
		"response": response,
	}
	return ctx.Response().View().Make("webpayplus/transaction_captured.tmpl", data)
}
