package controller

import (
	"fmt"
	"net/http"
	"qr-nikahan/config"
	"qr-nikahan/domain"
	"qr-nikahan/internal/helper"
	"strconv"
	"strings"
	"time"

	"html/template"

	"github.com/julienschmidt/httprouter"
)

type Controller struct {
	sheetsService domain.SpreadsheetService
	qrService     domain.QRService
}

func NewQRController(sheetsService domain.SpreadsheetService, qrService domain.QRService) Controller {
	return Controller{
		sheetsService: sheetsService,
		qrService:     qrService,
	}
}

func (obj *Controller) Check(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var (
		key            string = params.ByName("key")
		sheetData      []domain.GETSheet
		tmpl           *template.Template
		now            time.Time = time.Now()
		ableToScanAfer time.Time
		timeLayout     string = "2006-01-02T15:04:05"
		err            error
	)

	ableToScanAfer, err = time.Parse(timeLayout, config.AbleToScanAfer)
	if err != nil {
		helper.ERROR("Failed parsing time")
		obj.internalServerError(w)

		return
	}

	if now.Before(ableToScanAfer) {
		helper.INFO("Trying to scan before the date")

		tmpl, err = template.ParseFiles("./assets/web/not_before.html")
		if err != nil {
			helper.ERROR("Failed parsing not_before.html")

			return
		}

		if err = tmpl.Execute(w, nil); err != nil {
			helper.ERROR("Failed rendering not_before.html")
			obj.internalServerError(w)

			return
		}

		return
	}

	err, sheetData = obj.sheetsService.GetAllData()
	if err != nil {
		helper.ERROR("Failed getting all sheet data")
		obj.internalServerError(w)

		return
	}

	for _, data := range sheetData {
		if key == data.Key {
			var row int

			if data.ScannedAt != "" {
				var scannedAt time.Time

				scannedAt, err = time.Parse(timeLayout, strings.Replace(data.ScannedAt, "Z", "", -1))
				if err != nil {
					helper.ERROR(fmt.Sprintf("Failed parsing time for %s/%s;err: %s", data.Phone, data.Name, err.Error()))
					obj.internalServerError(w)

					return
				}

				if now.After(scannedAt) {
					tmpl, err = template.ParseFiles("./assets/web/redundant.html")
					if err != nil {
						helper.ERROR(fmt.Sprintf("Failed parsing redundant.html for %s/%s;err: %s", data.Phone, data.Name, err.Error()))

						return
					}

					if err = tmpl.Execute(w, data); err != nil {
						helper.ERROR(fmt.Sprintf("Failed rendering redundant.html time for %s/%s;err: %s", data.Phone, data.Name, err.Error()))
						obj.internalServerError(w)

						return
					}

					return
				}
			}

			helper.INFO(fmt.Sprintf("QRValid for %s/%s", data.Phone, data.Name))

			tmpl, err = template.ParseFiles("./assets/web/qr_valid.html")
			if err != nil {
				helper.ERROR(fmt.Sprintf("Failed parsing qr_valid.html for %s/%s;err: %s", data.Phone, data.Name, err.Error()))
				obj.internalServerError(w)

				return
			}

			if err = tmpl.Execute(w, data); err != nil {
				helper.ERROR(fmt.Sprintf("Failed rendering qr_valid.html for %s/%s;err: %s", data.Phone, data.Name, err.Error()))
				obj.internalServerError(w)

				return
			}

			row, err = strconv.Atoi(data.ID)
			if err != nil {
				helper.ERROR(fmt.Sprintf("Failed converting id to int for %s/%s;err: %s", data.Phone, data.Name, err.Error()))
				obj.internalServerError(w)

				return
			}

			if err = obj.sheetsService.ScannedQR(row + 1); err != nil {
				helper.ERROR(fmt.Sprintf("Failed updating sheet for %s/%s;err: %s", data.Phone, data.Name, err.Error()))
				obj.internalServerError(w)

				return
			}

			return
		}
	}

	tmpl, err = template.ParseFiles("./assets/web/qr_invalid.html")
	if err != nil {
		helper.ERROR("Failed parsing qr_invalid.html")
		obj.internalServerError(w)

		return
	}

	if err = tmpl.Execute(w, nil); err != nil {
		helper.ERROR("Failed rendering qr_invalid.html")
		obj.internalServerError(w)

		return
	}

	return
}

func (obj *Controller) internalServerError(w http.ResponseWriter) {
	var (
		tmpl *template.Template
		err  error
	)

	tmpl, err = template.ParseFiles("./assets/web/internal_server_error.html")
	if err != nil {
		helper.ERROR("Failed parsing internal_server_error.html")
	}

	if err = tmpl.Execute(w, nil); err != nil {
		helper.ERROR("Failed rendering internal_server_error.html")
	}
}
