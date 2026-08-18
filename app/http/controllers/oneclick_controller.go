package controllers

import (
	"fmt"
	"strconv"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/ppastene/transbank-sdk-go"
	"github.com/ppastene/transbank-sdk-go/oneclick"
)

type OneclickController struct {
	inscription *oneclick.MallInscription
	transaction *oneclick.MallTransaction
}

func NewOneclickController() (*OneclickController, error) {
	options := transbank.Options{
		Environment:  transbank.Integration,
		CommerceCode: "597055555541",
		ApiKey:       "579B532A7440BB0C9079DED94D31EA1615BACEB56610332264630D42D0A36B1C",
	}

	inscription, err := oneclick.NewMallInscription(options)
	if err != nil {
		return nil, err
	}

	transaction, err := oneclick.NewMallTransaction(options)
	if err != nil {
		return nil, err
	}

	return &OneclickController{
		inscription: inscription,
		transaction: transaction,
	}, nil
}

func (o *OneclickController) Index(ctx http.Context) http.Response {
	baseUrl := facades.Config().Env("APP_URL", "").(string)
	port := facades.Config().Env("APP_PORT", "").(string)
	var fullUrl string
	if port != "" {
		fullUrl = fmt.Sprintf("%s:%s", baseUrl, port)
	} else {
		fullUrl = baseUrl
	}

	data := map[string]any{
		"title":         "Oneclick Mall",
		"responseUrl":   fullUrl,
		"isDeferred":    false,
		"commerceCode1": 597055555542,
		"commerceCode2": 597055555543,
	}
	return ctx.Response().View().Make("oneclick/index.tmpl", data)
}

func (o *OneclickController) StartInscription(ctx http.Context) http.Response {
	request := ctx.Request().All()
	username := ctx.Request().Input("username")
	email := ctx.Request().Input("email")
	responseUrl := ctx.Request().Input("response_url")

	response, err := o.inscription.Start(username, email, responseUrl)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}

	data := map[string]any{
		"title":    "Inscripción Creada",
		"request":  request,
		"response": response,
	}

	return ctx.Response().View().Make("oneclick/mallinscription_start.tmpl", data)
}

func (o *OneclickController) FinishInscription(ctx http.Context) http.Response {
	// Revisar la respuesta entregada por transbank:
	// Si viene TBK_TOKEN en la URL y ResponseCode = 0, la inscripcion fue exitosa
	// Si viene TBK_TOKEN en la URL y ResponseCode = -1, la inscripcion fallo porque el pago no fue exitoso
	// Si viene TBK_TOKEN, TBK_ORDEN_COMPRA, TBK_ID_SESION en la url y ResponseCode = -96, la inscripcion fallo por cancelar el proceso
	// En caso de timeout, tienen que pasar 4 minutos en el formulario de transbank (o 10 en integracion) para que sea considerada como tal
	// Ahí se devuele TBK_TOKEN, revisar que ResponseCode devuelve

	var token string = ctx.Request().Input("TBK_TOKEN")
	response, err := o.inscription.Finish(token)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	data := map[string]any{
		"title":    "Inscripción finalizada",
		"request":  ctx.Request().All(),
		"response": response,
	}
	return ctx.Response().View().Make("oneclick/mallinscription_finish.tmpl", data)
}

func (o *OneclickController) DeleteInscription(ctx http.Context) http.Response {
	request := ctx.Request().All()
	tbkUser := ctx.Request().Input("tbk_user")
	username := ctx.Request().Input("username")
	err := o.inscription.Delete(tbkUser, username)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	data := map[string]any{
		"title":       "Inscripción Eliminada",
		"request":     request,
		"status_code": ctx.Response(),
	}
	return ctx.Response().View().Make("oneclick/mallinscription_delete.tmpl", data)
}

func (o *OneclickController) AuthorizeTransaction(ctx http.Context) http.Response {
	request := ctx.Request().All()
	var details []oneclick.MallDetails
	for i := 0; ; i++ {
		amountKey := fmt.Sprintf("detail[%d][amount]", i)
		amountVal, exists := request[amountKey]

		if !exists {
			break
		}

		installmentsKey := fmt.Sprintf("detail[%d][installments_number]", i)
		installmentsVal, exists := request[installmentsKey]

		if !exists {
			break
		}

		enabledKey := fmt.Sprintf("detail[%d][enabled]", i)
		if _, enabled := request[enabledKey]; !enabled {
			continue
		}

		amountStr := amountVal.(string)
		installString := installmentsVal.(string)
		commCodeKey := fmt.Sprintf("detail[%d][commerce_code]", i)
		commCode, _ := request[commCodeKey].(string)

		buyOrdKey := fmt.Sprintf("detail[%d][buy_order]", i)
		buyOrd, _ := request[buyOrdKey].(string)

		amount, _ := strconv.ParseFloat(amountStr, 64)
		installments, _ := strconv.Atoi(installString)

		fmt.Println(commCode, buyOrd, amount, installments)

		details = append(details, oneclick.MallDetails{
			Amount:             amount,
			CommerceCode:       commCode,
			BuyOrder:           buyOrd,
			InstallmentsNumber: installments,
		})
	}
	username := ctx.Request().Input("username")
	tbkUser := ctx.Request().Input("tbk_user")
	buyOrder := ctx.Request().Input("buy_order")
	response, err := o.transaction.Authorize(username, tbkUser, buyOrder, details)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}

	data := map[string]any{
		"title":    "Transacción Autorizada",
		"request":  request,
		"response": response,
	}

	return ctx.Response().View().Make("oneclick/malltransaction_authorize.tmpl", data)
}

func (o *OneclickController) StatusTransaction(ctx http.Context) http.Response {
	request := ctx.Request().All()
	buyOrder := ctx.Request().Input("buy_order")
	response, err := o.transaction.Status(buyOrder)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}

	data := map[string]any{
		"title":    "Transacción Autorizada",
		"request":  request,
		"response": response,
	}

	return ctx.Response().View().Make("oneclick/malltransaction_status.tmpl", data)
}

func (o *OneclickController) RefundTransaction(ctx http.Context) http.Response {
	request := ctx.Request().All()
	buyOrder := ctx.Request().Input("buy_order")
	commerceCode := ctx.Request().Input("child_commerce_code")
	detailBuyOrder := ctx.Request().Input("child_buy_order")
	amountStr := ctx.Request().Input("amount")
	amount, err := strconv.ParseFloat(amountStr, 64)

	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}

	response, err := o.transaction.Refund(buyOrder, commerceCode, detailBuyOrder, amount)
	if err != nil {
		return ctx.Response().Status(500).Json(map[string]any{
			"error": err,
		})
	}
	data := map[string]any{
		"title":    "Transacción reembolsada",
		"request":  request,
		"response": response,
	}
	return ctx.Response().View().Make("oneclick/malltransaction_refund.tmpl", data)
}
