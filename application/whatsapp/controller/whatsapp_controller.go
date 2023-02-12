package controller

import (
	"fmt"
	"net/http"
	"qr-nikahan/domain"
	"qr-nikahan/internal/helper"

	"github.com/julienschmidt/httprouter"
)

type Controller struct {
	waService     domain.WhatsAppService
	sheetsService domain.SpreadsheetService
	qrService     domain.QRService
}

func NewWhatsAppController(service domain.WhatsAppService, sheetsService domain.SpreadsheetService, qrService domain.QRService) Controller {
	return Controller{
		waService:     service,
		sheetsService: sheetsService,
		qrService:     qrService,
	}
}

func (obj *Controller) Blast(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var (
		resp      domain.Response
		sheetData []domain.GETSheet
		err       error
	)

	err, sheetData = obj.sheetsService.GetAllData()
	if err != nil {
		helper.ERROR(err.Error())

		resp.Error.Err = err
		resp.Error.Desc = "Failed getting all sheet data"

		helper.Response(w, resp, http.StatusInternalServerError)

		return
	}

	for _, data := range sheetData {
		var (
			qrImage []byte
			key     string
			errDesc string = fmt.Sprintf("Failed generating QR for %d/%s ", data.Phone, data.Name)
		)

		err, qrImage, key = obj.qrService.Generate(data)
		if err != nil {
			helper.ERROR(err.Error())

			resp.Error.Err = err
			resp.Error.Desc = errDesc

			helper.Response(w, resp, http.StatusInternalServerError)

			return
		}

		if err = obj.waService.SendMessage(data.Name, data.Phone, qrImage); err != nil {
			helper.ERROR(err.Error())

			resp.Error.Err = err
			resp.Error.Desc = errDesc

			helper.Response(w, resp, http.StatusInternalServerError)

			return
		}

		if err = obj.sheetsService.SentInvitation(data.ID+1, key); err != nil {
			helper.ERROR(err.Error())

			resp.Error.Err = err
			resp.Error.Desc = errDesc

			helper.Response(w, resp, http.StatusInternalServerError)

			return
		}
	}

	helper.INFO("Succeed in blasting WA and updating sheets")
	resp.Data = "OK"

	helper.Response(w, resp, http.StatusOK)

	return
}
