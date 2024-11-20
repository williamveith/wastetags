package qrcodegen

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/color"
	"os"

	"github.com/skip2/go-qrcode"
)

var defaultConfigs = &QRCodeSettings{
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

func generateQRCode(content string, settings *QRCodeSettings) (*qrcode.QRCode, error) {
	settings = applyDefaults(settings)
	qrCode, err := qrcode.New(content, settings.RecoveryLevel)
	qrCode.DisableBorder = settings.DisableBorder
	qrCode.BackgroundColor = settings.BackgroundColor
	qrCode.ForegroundColor = settings.ForegroundColor

	return qrCode, wrapError(err, "failed to generate QR code")
}

func applyDefaults(settings *QRCodeSettings) *QRCodeSettings {
	if settings == nil {
		return defaultConfigs
	}
	return settings
}

func wrapError(err error, message string) error {
	if err != nil {
		return fmt.Errorf("%s: %w", message, err)
	}
	return nil
}

func normalizeColorComponent(component uint32) float64 {
	return float64(component) / 255.0
}

func GetMostContrastingColor(providedColor color.Color) color.Color {
	r, g, b, _ := providedColor.RGBA()
	red := normalizeColorComponent(r)
	green := normalizeColorComponent(g)
	blue := normalizeColorComponent(b)

	luminance := 0.299*red + 0.587*green + 0.114*blue

	if luminance > 0.5 {
		return color.Black
	}
	return color.White
}

func Image(content string, settings *QRCodeSettings) []byte {
	settings = applyDefaults(settings)
	qrCode, _ := generateQRCode(content, settings)
	var buffer bytes.Buffer
	qrCode.Write(settings.Size, &buffer)

	return buffer.Bytes()
}

func ImageFile(content, fileName string, settings *QRCodeSettings) error {
	settings = applyDefaults(settings)
	imageBytes := Image(content, settings)
	return os.WriteFile(fileName, imageBytes, 0644)
}

func Base64(content string, settings *QRCodeSettings) string {
	settings = applyDefaults(settings)
	imageBytes := Image(content, settings)
	return base64.StdEncoding.EncodeToString(imageBytes)
}

func DataURI(content string, settings *QRCodeSettings) string {
	settings = applyDefaults(settings)
	base64String := Base64(content, settings)
	return fmt.Sprintf("data:image/png;base64,%s", base64String)
}

func SVG(content string, settings *QRCodeSettings) string {
	settings = applyDefaults(settings)
	qrCode, _ := generateQRCode(content, settings)
	return qrCode.ToSmallString(true)
}

func SVGFile(content, fileName string, settings *QRCodeSettings) error {
	settings = applyDefaults(settings)
	svgString := SVG(content, settings)
	return os.WriteFile(fileName, []byte(svgString), 0644)
}
