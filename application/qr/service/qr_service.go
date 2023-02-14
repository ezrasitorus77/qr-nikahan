package service

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"qr-nikahan/config"
	"qr-nikahan/domain"
	log "qr-nikahan/internal/helper"

	"github.com/skip2/go-qrcode"
)

type service struct {
}

func NewQRService() (obj domain.QRService) {
	return &service{}
}

func (obj *service) Generate(data domain.GETSheet) (err error, qrImage []byte, key string) {
	var (
		pnga       []byte
		img        image.Image
		fimg       image.Image
		out        *os.File
		in         *os.File
		uniqueIden string = fmt.Sprintf("%#v", data)
		path       string = fmt.Sprintf("./assets/qrimage/%s%s.jpeg", data.ID, data.Name)
	)

	key = base64.StdEncoding.EncodeToString([]byte(uniqueIden))

	pnga, err = qrcode.Encode(config.BaseURL+"/check/"+key, qrcode.Medium, 256)
	if err != nil {
		log.ERROR(fmt.Sprintf("Error generate QR for %s;id %s;err=%s", data.Name, data.ID, err.Error()))

		return
	}

	img, _, err = image.Decode(bytes.NewReader(pnga))
	if err != nil {
		log.ERROR(fmt.Sprintf("Error generate QR for %s;id %s;err=%s", data.Name, data.ID, err.Error()))

		return
	}

	out, _ = os.Create(path)
	defer out.Close()

	var opts jpeg.Options
	opts.Quality = 1

	err = jpeg.Encode(out, img, &opts)
	if err != nil {
		log.ERROR(fmt.Sprintf("Error generate QR for %s;id %s;err=%s", data.Name, data.ID, err.Error()))

		return
	}

	in, err = os.Open(path)
	if err != nil {
		log.ERROR(fmt.Sprintf("Error generate QR for %s;id %s;err=%s", data.Name, data.ID, err.Error()))

		return
	}

	defer in.Close()

	bimg := new(bytes.Buffer)
	fimg, err = jpeg.Decode(in)
	if err != nil {
		log.ERROR(fmt.Sprintf("Error generate QR for %s;id %s;err=%s", data.Name, data.ID, err.Error()))

		return
	}

	if err = jpeg.Encode(bimg, fimg, nil); err != nil {
		log.ERROR(fmt.Sprintf("Error generate QR for %s;id %s;err=%s", data.Name, data.ID, err.Error()))

		return
	}

	qrImage = bimg.Bytes()

	return
}
