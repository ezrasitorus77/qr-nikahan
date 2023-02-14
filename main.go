package main

import (
	"net/http"

	"qr-nikahan/domain"
	"qr-nikahan/internal/helper"
	"qr-nikahan/internal/middleware"

	_ "qr-nikahan/config"

	qrService "qr-nikahan/application/qr/service"
	sheetService "qr-nikahan/application/sheets/service"
	waService "qr-nikahan/application/whatsapp/service"

	qrController "qr-nikahan/application/qr/controller"
	waController "qr-nikahan/application/whatsapp/controller"

	"github.com/julienschmidt/httprouter"
)

func main() {
	var (
		qrServ     domain.QRService          = qrService.NewQRService()
		sheetsServ domain.SpreadsheetService = sheetService.NewSpreadsheetService()
		waServ     domain.WhatsAppService    = waService.NewWhatsAppService()

		waCont waController.Controller = waController.NewWhatsAppController(waServ, sheetsServ, qrServ)
		qrCont qrController.Controller = qrController.NewQRController(sheetsServ, qrServ)

		router        *httprouter.Router = httprouter.New()
		logMiddleware middleware.LogMiddleware
		server        http.Server

		err error
	)

	router.POST("/blast", waCont.Blast)
	router.GET("/check/:key", qrCont.Check)

	logMiddleware.Handler = router

	server.Addr = "localhost:8080"
	server.Handler = &logMiddleware

	helper.INFO("Running...")
	err = server.ListenAndServe()
	if err != nil {
		helper.PANIC(err.Error())
	}
}
