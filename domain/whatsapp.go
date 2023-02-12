package domain

type (
	WhatsAppService interface {
		SendMessage(name string, phone int, qrImage []byte) (err error)
	}
)
