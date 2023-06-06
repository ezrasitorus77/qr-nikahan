package domain

type (
	QRService interface {
		Generate(data GETSheet) (err error, qrImage []byte, key string)
	}
)
