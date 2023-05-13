package main

import (
	"encoding/json"
	"net/http"

	"qr-nikahan/config"
	"qr-nikahan/domain"
	"qr-nikahan/internal/helper"
	"qr-nikahan/internal/middleware"

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

	router.GET("/", func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		var (
			about       = make(map[string]string)
			jsonEncoder *json.Encoder
		)

		about["Application Name"] = "QR Nikahan"
		about["Version"] = "1.0"

		w.WriteHeader(http.StatusOK)

		jsonEncoder = json.NewEncoder(w)
		jsonEncoder.Encode(about)
	})

	router.POST("/blast", waCont.Blast)
	router.GET("/check/:key", qrCont.Check)

	logMiddleware.Handler = router

	server.Addr = ":" + config.IPPort
	server.Handler = &logMiddleware

	helper.INFO("Running...")
	err = server.ListenAndServe()
	if err != nil {
		helper.PANIC(err.Error())
	}
}
