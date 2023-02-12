package domain

type (
	GETSheet struct {
		ID        int    `json:"No"`
		Name      string `json:"Nama"`
		Group     string `json:"Grup"`
		Phone     int    `json:"Phone"`
		Key       string `json:"Key"`
		SentAt    string `json:"SentAt"`
		ScannedAt string `json:"ScannedAt"`
	}

	SpreadsheetService interface {
		GetAllData() (err error, data []GETSheet)
		SentInvitation(row int, key string) (err error)
		ScannedQR(row int) (err error)
		GetRowByNameAndPhone(name string, phone int, data []GETSheet) (err error, row int)
	}
)
