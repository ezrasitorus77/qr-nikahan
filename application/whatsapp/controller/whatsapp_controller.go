package controller

import (
	"fmt"
	"net/http"
	"qr-nikahan/domain"
	"qr-nikahan/internal/helper"
	"strconv"
	"strings"

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
			phone   string = strings.Replace(data.Phone, " ", "", -1)
			errDesc string = fmt.Sprintf("Failed generating QR for %s/%s ", data.Phone, data.Name)
			row     int
		)

		if phone != "" && data.ID != "" && data.Name != "" && data.SentAt == "" && data.ScannedAt == "" {
			err, qrImage, key = obj.qrService.Generate(data)
			if err != nil {
				helper.ERROR(err.Error())

				resp.Error.Err = err
				resp.Error.Desc = errDesc

				helper.Response(w, resp, http.StatusInternalServerError)

				return
			}

			if err = obj.waService.SendMessage(data.Name, phone, qrImage); err != nil {
				helper.ERROR(err.Error())

				resp.Error.Err = err
				resp.Error.Desc = errDesc

				helper.Response(w, resp, http.StatusInternalServerError)

				return
			}

			row, err = strconv.Atoi(data.ID)
			if err != nil {
				helper.ERROR(err.Error())

				resp.Error.Err = err
				resp.Error.Desc = errDesc

				helper.Response(w, resp, http.StatusInternalServerError)
			}

			if err = obj.sheetsService.SentInvitation(row+1, key); err != nil {
				helper.ERROR(err.Error())

				resp.Error.Err = err
				resp.Error.Desc = errDesc

				helper.Response(w, resp, http.StatusInternalServerError)

				return
			}
		}
	}

	helper.INFO("Succeed in blasting WA and updating sheets")
	resp.Data = "OK"

	helper.Response(w, resp, http.StatusOK)

	return
}
