package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"qr-nikahan/domain"
	"qr-nikahan/internal/consts"
	log "qr-nikahan/internal/helper"
	"strconv"
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

	srv, err = sheets.NewService(context.Background(), option.WithCredentialsFile(consts.ClientSecretPath), option.WithScopes(sheets.SpreadsheetsScope))
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

	if _, err = obj.sheetsService.Spreadsheets.Values.Update(consts.SheetsID, fmt.Sprintf("%s!%s%d", consts.SheetName, consts.SentAtColumn, row), &sheets.ValueRange{
		Values: [][]interface{}{{
			time,
		}},
	}).ValueInputOption("USER_ENTERED").Do(); err != nil {
		log.ERROR(err.Error())

		return
	}

	if _, err = obj.sheetsService.Spreadsheets.Values.Update(consts.SheetsID, fmt.Sprintf("%s!%s%d", consts.SheetName, consts.KeyColumn, row), &sheets.ValueRange{
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

	if _, err = obj.sheetsService.Spreadsheets.Values.Update(consts.SheetsID, fmt.Sprintf("%s!%s%d", consts.SheetName, consts.ScannedAtColumn, row), &sheets.ValueRange{
		Values: [][]interface{}{{
			time,
		}},
	}).ValueInputOption("USER_ENTERED").Do(); err != nil {
		log.ERROR(err.Error())

		return
	}

	return
}

func (obj *service) GetRowByNameAndPhone(name string, phone int, data []domain.GETSheet) (err error, row int) {
	for idx, v := range data {
		if v.Name == name && v.Phone == phone {
			row = idx + 2

			return
		}
	}

	err = errors.New("No data found with Nama: " + name + " and Phone: " + strconv.Itoa(phone))
	log.ERROR(err.Error())

	return
}

func (obj *service) GetAllData() (err error, data []domain.GETSheet) {
	var (
		resp     *http.Response
		dataByte []byte
	)

	resp, err = http.Get(consts.GetAPI)
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
