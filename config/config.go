package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"qr-nikahan/internal/helper"
)

var (
	DBHost               string
	DBPort               string
	DBUsername           string
	DBPassword           string
	DBName               string
	DBDialect            string
	GetAPI               string
	SheetsID             string
	AbleToScanAfer       string
	BaseURL              string
	ServiceAccount       string
	ServiceAccountConfig ServiceAccountJSON
	err                  error
)

type ServiceAccountJSON struct {
	Type         string `json:"type"`
	ProjectID    string `json:"project_id"`
	PrivateKeyID string `json:"private_key_id"`
	PrivateKey   string `json:"private_key"`
	ClientEmail  string `json:"client_email"`
	ClientID     string `json:"client_id"`
	AuthURI      string `json:"auth_uri"`
	TokenURI     string `json:"token_uri"`
	AuthProvicer string `json:"auth_provider_x509_cert_url"`
	ClientCert   string `json:"client_x509_cert_url"`
}

func init() {
	DBHost = os.Getenv("DB_HOSTNAME")
	if DBHost == "" {
		helper.PANIC("DB_HOSTNAME not found")
	}

	DBPort = os.Getenv("DATABASE_PORT")
	if DBPort == "" {
		helper.PANIC("DATABASE_PORT not found")
	}

	DBUsername = os.Getenv("DATABASE_USERNAME")
	if DBUsername == "" {
		helper.PANIC("DATABASE_USERNAME not found")
	}

	DBPassword = os.Getenv("DATABASE_PASSWORD")
	if DBPassword == "" {
		helper.PANIC("DATABASE_PASSWORD not found")
	}

	DBName = os.Getenv("DB_NAME")
	if DBName == "" {
		helper.PANIC("DB_NAME not found")
	}

	DBDialect = os.Getenv("DB_DIALECT")
	if DBDialect == "" {
		helper.PANIC("DB_DIALECT not found")
	}

	GetAPI = os.Getenv("GET_API")
	if GetAPI == "" {
		helper.PANIC("GET_API not found")
	}

	SheetsID = os.Getenv("SHEETS_ID")
	if SheetsID == "" {
		helper.PANIC("SHEETS_ID not found")
	}

	AbleToScanAfer = os.Getenv("ABLE_TO_SCAN_AFTER")
	if AbleToScanAfer == "" {
		helper.PANIC("ABLE_TO_SCAN_AFTER not found")
	}

	BaseURL = os.Getenv("BASE_URL")
	if BaseURL == "" {
		helper.PANIC("BASE_URL not found")
	}

	ServiceAccount = os.Getenv("SERVICE_ACCOUNT_CONFIG")
	if ServiceAccount == "" {
		helper.PANIC("SERVICE_ACCOUNT_CONFIG not found")
	} else {
		var file []byte

		if err := json.Unmarshal([]byte(ServiceAccount), &ServiceAccountConfig); err != nil {
			helper.PANIC(err)
		}

		file, err = json.MarshalIndent(ServiceAccountConfig, "", " ")
		if err != nil {
			helper.PANIC(err)
		}

		if err = ioutil.WriteFile("service_account.json", file, 0644); err != nil {
			helper.PANIC(err)
		}
	}
}
