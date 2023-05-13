package domain

type (
	WhatsAppService interface {
		SendMessage(name, phone string, qrImage []byte) (err error)
		SendMarhusipMessage(name, phone string) (err error)
	}
)
