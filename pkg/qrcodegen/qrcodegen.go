package qrcodegen

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"os"

	"github.com/skip2/go-qrcode"
)

var configs = QRCodeSettings{
	Size:            256,
	RecoveryLevel:   qrcode.Medium,
	DisableBorder:   true,
	BackgroundColor: color.Transparent,
	ForegroundColor: color.Black,
}

type QRCodeSettings struct {
	Size            int
	RecoveryLevel   qrcode.RecoveryLevel
	DisableBorder   bool
	BackgroundColor color.Color
	ForegroundColor color.Color
}

func generateQRCode(content string) (*qrcode.QRCode, error) {
	qrCode, err := qrcode.New(content, configs.RecoveryLevel)
	if err != nil {
		return nil, err
	}

	qrCode.DisableBorder = configs.DisableBorder
	qrCode.BackgroundColor = configs.BackgroundColor
	qrCode.ForegroundColor = configs.ForegroundColor

	return qrCode, nil
}

func Image(content string) ([]byte, error) {
	qrCode, err := generateQRCode(content)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	err = qrCode.Write(configs.Size, &buffer)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func ImageFile(content, fileName string) error {
	imageBytes, err := Image(content)
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, imageBytes, 0644)
}

func Base64(content string) (string, error) {
	imageBytes, err := Image(content)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(imageBytes), err
}

func DataURI(content string) (string, error) {
	base64String, err := Base64(content)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("data:image/png;base64,%s", base64String), err
}

func SVG(content string) (string, error) {
	qrCode, err := generateQRCode(content)
	if err != nil {
		return "", err
	}

	return qrCode.ToSmallString(true), nil
}

func SVGFile(content, fileName string) error {
	svgString, err := SVG(content)
	if err != nil {
		return err
	}

	return os.WriteFile(fileName, []byte(svgString), 0644)
}
