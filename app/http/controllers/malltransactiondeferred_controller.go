package controllers

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/ppastene/transbank-sdk-go"
)

type MallTransactionDeferredController struct {
	mallTransaction *transbank.WebpayPlusMallTransaction
}

func NewMallTransactionDeferredController() *MallTransactionDeferredController {
	options := &transbank.Options{
		ApiKey:       "597055555581",
		CommerceCode: "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	}
	transaction := transbank.NewMallTransaction(options)
	return &MallTransactionDeferredController{
		mallTransaction: transaction,
	}
}

func (m *MallTransactionDeferredController) Index(ctx http.Context) http.Response {
	baseUrl := facades.Config().Env("APP_URL", "").(string)
	port := facades.Config().Env("APP_PORT", "").(string)
	var fullUrl string
	if port != "" {
		fullUrl = fmt.Sprintf("%s:%s", baseUrl, port)
	} else {
		fullUrl = baseUrl
	}
	data := map[string]any{
		"title":      "Webpay Plus Mall Diferido",
		"returnUrl":  fullUrl,
		"isDeferred": true,
	}
	return ctx.Response().View().Make("webpayplusmall/index.tmpl", data)
}

func (m *MallTransactionDeferredController) CreatedTransaction(ctx http.Context) http.Response {
	request := ctx.Request().All()
	var details []transbank.WebpayPlusMallDetails
	for i := 0; ; i++ {
		amountKey := fmt.Sprintf("detail[%d][amount]", i)
		val, exists := request[amountKey]

		if !exists {
			break
		}

		enabledKey := fmt.Sprintf("detail[%d][enabled]", i)
		if _, enabled := request[enabledKey]; !enabled {
			continue
		}

		amountStr := val.(string)
		commCodeKey := fmt.Sprintf("detail[%d][commerce_code]", i)
		commCode, _ := request[commCodeKey].(string)

		buyOrdKey := fmt.Sprintf("detail[%d][buy_order]", i)
		buyOrd, _ := request[buyOrdKey].(string)

		amount, _ := strconv.ParseFloat(amountStr, 64)

		details = append(details, transbank.WebpayPlusMallDetails{
			Amount:       amount,
			CommerceCode: commCode,
			BuyOrder:     buyOrd,
		})
	}
	parentBuyOrder := request["buy_order"].(string)
	sessionId := request["session_id"].(string)
	returnUrl := request["return_url"].(string)
	response, err := m.mallTransaction.Create(parentBuyOrder, sessionId, returnUrl, details)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}

	data := map[string]any{
		"title":    "Transacción Mall Diferida Creada",
		"request":  request,
		"response": response,
	}

	return ctx.Response().View().Make("webpayplusmall/malltransaction_created.tmpl", data)
}

func (m *MallTransactionDeferredController) CommitedTransaction(ctx http.Context) http.Response {
	// Flujo abortado: Se revisa si en la url TBK_TOKEN existe
	// De ser asi se considera la transaccion como abortada
	queryParams := ctx.Request().Queries()
	if queryParams["TBK_TOKEN"] != "" {
		data := map[string]any{
			"title":       "Transacción Mall Diferida Abortada",
			"queryParams": queryParams,
		}
		return ctx.Response().View().Make("webpayplusmall/malltransaction_aborted.tmpl", data)
	}

	// Si no está presente pues es un flujo normal
	request := ctx.Request().All()
	var token string = ctx.Request().Input("token_ws")
	response, err := m.mallTransaction.Commit(token)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}

	data := map[string]any{
		"title":      "Transacción Mall Diferida Confirmada",
		"request":    request,
		"response":   response,
		"isApproved": response.IsApproved(),
	}
	return ctx.Response().View().Make("webpayplusmall/malltransaction_committed.tmpl", data)
}

func (m *MallTransactionDeferredController) GetTransactionStatus(ctx http.Context) http.Response {
	request := ctx.Request().All()
	response, err := m.mallTransaction.Status(request["token"].(string))
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	data := map[string]any{
		"title":    "Status Transacción Mall Diferida",
		"request":  request,
		"response": response,
	}
	return ctx.Response().View().Make("webpayplusmall/mall_status.tmpl", data)
}

func (m *MallTransactionDeferredController) RefundTransaction(ctx http.Context) http.Response {
	request := ctx.Request().All()
	token := ctx.Request().Input("token")
	childBuyOrder := ctx.Request().Input("child_buy_order")
	childCommerceCode := ctx.Request().Input("child_commerce_code")
	amountStr := ctx.Request().Input("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)

	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}

	response, err := m.mallTransaction.Refund(token, childBuyOrder, childCommerceCode, amount)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	data := map[string]any{
		"title":    "Transacción Mall Diferida Reembolsada",
		"request":  request,
		"response": response,
	}
	return ctx.Response().View().Make("webpayplusmall/malltransaction_refunded.tmpl", data)
}

func (t *MallTransactionDeferredController) CaptureTransaction(ctx http.Context) http.Response {
	request := ctx.Request().All()
	token := ctx.Request().Input("token")
	childBuyOrder := ctx.Request().Input("child_buy_order")
	childCommerceCode := ctx.Request().Input("child_commerce_code")
	childAuthCode := ctx.Request().Input("child_authorization_code")
	amountStr := ctx.Request().Input("capture_amount")
	amount, err := strconv.ParseFloat(amountStr, 64)

	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}

	response, err := t.mallTransaction.Capture(token, childCommerceCode, childBuyOrder, childAuthCode, amount)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	data := map[string]any{
		"title":    "Transacción Mall Diferida Capturada",
		"request":  request,
		"response": response,
	}
	return ctx.Response().View().Make("webpayplusmall/malltransaction_captured.tmpl", data)
}
