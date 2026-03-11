package routes

import (
	"goravel/app/http/controllers"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
)

func Web() {
	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("index.tmpl", map[string]any{
			"version": support.Version,
		})
	})
	// Webpay Plus
	transactionController := controllers.NewTransactionController()
	facades.Route().Get("/webpayplus/", transactionController.Index)
	facades.Route().Post("/webpayplus/create", transactionController.CreatedTransaction)
	facades.Route().Any("/webpayplus/returnUrl", transactionController.CommitedTransaction)
	facades.Route().Post("/webpayplus/status", transactionController.GetTransactionStatus)
	facades.Route().Post("/webpayplus/refund", transactionController.RefundTransaction)

	// Webpay Plus Diferido
	transactionDeferredController := controllers.NewTransactionDeferredController()
	facades.Route().Get("/webpayplusdeferred/", transactionDeferredController.Index)
	facades.Route().Post("/webpayplusdeferred/create", transactionDeferredController.CreatedTransaction)
	facades.Route().Any("/webpayplusdeferred/returnUrl", transactionDeferredController.CommitedTransaction)
	facades.Route().Post("/webpayplusdeferred/status", transactionDeferredController.GetTransactionStatus)
	facades.Route().Post("/webpayplusdeferred/refund", transactionDeferredController.RefundTransaction)
	facades.Route().Post("/webpayplusdeferred/capture", transactionDeferredController.CaptureTransaction)

	// Webpay Plus Mall
	mallTransactionController := controllers.NewMallTransactionController()
	facades.Route().Get("/webpayplusmall", mallTransactionController.Index)
	facades.Route().Post("/webpayplusmallcreate", mallTransactionController.CreatedTransaction)
	facades.Route().Any("/webpayplusmall/returnUrl", mallTransactionController.CommitedTransaction)
	facades.Route().Post("/webpayplusmall/status", mallTransactionController.GetTransactionStatus)
	facades.Route().Post("/webpayplusmall/refund", mallTransactionController.RefundTransaction)

	// Webpay Plus Mall Deferred
	mallTransactionDeferredController := controllers.NewMallTransactionDeferredController()
	facades.Route().Get("/webpayplusmalldeferred/", mallTransactionDeferredController.Index)
	facades.Route().Post("/webpayplusmalldeferred/create", mallTransactionDeferredController.CreatedTransaction)
	facades.Route().Any("/webpayplusmalldeferred/returnUrl", mallTransactionDeferredController.CommitedTransaction)
	facades.Route().Post("/webpayplusmalldeferred/status", mallTransactionDeferredController.GetTransactionStatus)
	facades.Route().Post("/webpayplusmalldeferred/refund", mallTransactionDeferredController.RefundTransaction)
	facades.Route().Post("/webpayplusmalldeferred/capture", mallTransactionDeferredController.CaptureTransaction)
}
