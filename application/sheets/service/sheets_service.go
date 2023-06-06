package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"qr-nikahan/config"
	"qr-nikahan/domain"
	"qr-nikahan/internal/consts"
	log "qr-nikahan/internal/helper"
	"time"

	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type service struct {
	sheetsService *sheets.Service
}

func NewSpreadsheetService() (obj domain.SpreadsheetService) {
	var (
		srv *sheets.Service
		err error
	)

	srv, err = sheets.NewService(context.Background(), option.WithCredentialsFile("./service_account.json"), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		log.PANIC(err.Error())

		return
	}

	obj = &service{
		sheetsService: srv,
	}

	return
}

func (obj *service) SentInvitation(row int, key string) (err error) {
	var time string = time.Now().Format("2006-01-02 15:04:05")

	if _, err = obj.sheetsService.Spreadsheets.Values.Update(config.SheetsID, fmt.Sprintf("%s!%s%d", consts.SheetName, consts.SentAtColumn, row), &sheets.ValueRange{
		Values: [][]interface{}{{
			time,
		}},
	}).ValueInputOption("USER_ENTERED").Do(); err != nil {
		log.ERROR(err.Error())

		return
	}

	if _, err = obj.sheetsService.Spreadsheets.Values.Update(config.SheetsID, fmt.Sprintf("%s!%s%d", consts.SheetName, consts.KeyColumn, row), &sheets.ValueRange{
		Values: [][]interface{}{{
			key,
		}},
	}).ValueInputOption("USER_ENTERED").Do(); err != nil {
		log.ERROR(err.Error())

		return
	}

	return
}

func (obj *service) ScannedQR(row int) (err error) {
	var time string = time.Now().Format("2006-01-02 15:04:05")

	if _, err = obj.sheetsService.Spreadsheets.Values.Update(config.SheetsID, fmt.Sprintf("%s!%s%d", consts.SheetName, consts.ScannedAtColumn, row), &sheets.ValueRange{
		Values: [][]interface{}{{
			time,
		}},
	}).ValueInputOption("USER_ENTERED").Do(); err != nil {
		log.ERROR(err.Error())

		return
	}

	return
}

func (obj *service) GetAllData() (err error, data []domain.GETSheet) {
	var (
		resp     *http.Response
		dataByte []byte
	)

	resp, err = http.Get(config.GetAPI)
	if err != nil {
		log.PANIC(err.Error())

		return
	}

	defer resp.Body.Close()

	dataByte, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.PANIC(err.Error())

		return
	}

	if err = json.Unmarshal([]byte(dataByte), &data); err != nil {
		log.PANIC(err.Error())

		return
	}

	return
}
