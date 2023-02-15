package service

import (
	"context"
	"fmt"
	"os"
	"qr-nikahan/config"
	"qr-nikahan/domain"

	"qr-nikahan/internal/helper"

	_ "github.com/lib/pq"
	"github.com/mdp/qrterminal"
	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/store"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

type service struct {
	client *whatsmeow.Client
}

func NewWhatsAppService() (obj domain.WhatsAppService) {
	var (
		container   *sqlstore.Container
		deviceStore *store.Device
		client      *whatsmeow.Client
		qrChan      <-chan whatsmeow.QRChannelItem
		dbAdress    string = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", config.DBUsername, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
		err         error
	)

	container, err = sqlstore.New(config.DBDialect, dbAdress, waLog.Noop)
	if err != nil {
		helper.PANIC(err.Error())

		return
	}

	deviceStore, err = container.GetFirstDevice()
	if err != nil {
		helper.PANIC(err.Error())

		return
	}

	client = whatsmeow.NewClient(deviceStore, waLog.Noop)
	if client.Store.ID == nil {
		qrChan, _ = client.GetQRChannel(context.Background())
		err = client.Connect()
		if err != nil {
			helper.PANIC(err.Error())

			return
		}

		for evt := range qrChan {
			if evt.Event == "code" {
				qrterminal.GenerateHalfBlock(evt.Code, qrterminal.L, os.Stdout)
			} else {
				helper.INFO(fmt.Sprintf("Login event: %s", evt.Event))
			}
		}
	} else {
		err = client.Connect()
		if err != nil {
			helper.PANIC(err.Error())

			return
		}
	}

	obj = &service{
		client: client,
	}

	return
}

func (obj *service) SendMessage(name, phone string, qrImage []byte) (err error) {
	var (
		sendResp   whatsmeow.SendResponse
		uploadResp whatsmeow.UploadResponse
		imageMsg   waProto.ImageMessage
	)

	uploadResp, err = obj.client.Upload(context.Background(), qrImage, whatsmeow.MediaImage)
	if err != nil {
		helper.ERROR("Failed send message to " + phone + "/" + name)

		return
	}

	imageMsg.Caption = proto.String(helper.CreateMessage(name))
	imageMsg.Mimetype = proto.String("image/jpeg")
	imageMsg.Url = &uploadResp.URL
	imageMsg.DirectPath = &uploadResp.DirectPath
	imageMsg.MediaKey = uploadResp.MediaKey
	imageMsg.FileEncSha256 = uploadResp.FileEncSHA256
	imageMsg.FileSha256 = uploadResp.FileSHA256
	imageMsg.FileLength = &uploadResp.FileLength

	sendResp, err = obj.client.SendMessage(context.Background(), types.JID{
		User:   phone,
		Server: types.DefaultUserServer,
	}, &waProto.Message{
		ImageMessage: &imageMsg,
	})

	if err != nil {
		helper.ERROR("Failed send message to " + phone + "/" + name)

		return
	}

	helper.INFO("Succeed sending to " + phone + "/" + name)
	helper.INFO(sendResp)

	return
}
