package domain

type (
	GETSheet struct {
		ID        string `json:"No"`
		Name      string `json:"Nama"`
		Group     string `json:"Grup"`
		Phone     string `json:"Phone"`
		Key       string `json:"Key"`
		SentAt    string `json:"SentAt"`
		ScannedAt string `json:"ScannedAt"`
	}

	GETSheetMarhusip struct {
		ID     string `json:"No"`
		Phone  string `json:"NoHP"`
		Name   string `json:"Nama"`
		SentAt string `json:"SentAt"`
	}

	SpreadsheetService interface {
		GetAllData() (err error, data []GETSheet)
		SentInvitation(row int, key string) (err error)
		ScannedQR(row int) (err error)

		GetAllMarhusipData() (err error, data []GETSheetMarhusip)
		SentMarhusipInvitation(row int, key string) (err error)
	}
)
