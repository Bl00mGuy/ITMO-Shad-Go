package main

import (
	"image"
	"image/color"
	"strings"
)

func GenerateClockImage(request *ClockRequest) *image.RGBA {
	const (
		charWidth  = 8
		charHeight = 12
		colonWidth = 4
	)

	imgWidth := (6*charWidth + 2*colonWidth) * request.Scale
	imgHeight := charHeight * request.Scale
	img := image.NewRGBA(image.Rect(0, 0, imgWidth, imgHeight))

	offset := 0
	for _, char := range request.Time {
		charBitmap := getBitmapForChar(char)
		drawCharOnImage(img, charBitmap, offset, request.Scale)
		if string(char) == ":" {
			offset += colonWidth
		} else {
			offset += charWidth
		}
	}

	return img
}

func drawCharOnImage(img *image.RGBA, bitmap string, xOffset, scale int) {
	lines := strings.Split(bitmap, "\n")
	for y, line := range lines {
		for x, pixel := range line {
			colorParameter := color.RGBA{R: 255, G: 255, B: 255, A: 255}
			if pixel == '1' {
				colorParameter = Cyan
			}
			fillScaledPixel(img, xOffset+x, y, scale, colorParameter)
		}
	}
}

func fillScaledPixel(img *image.RGBA, x, y, scale int, fillColor color.Color) {
	for i := 0; i < scale; i++ {
		for j := 0; j < scale; j++ {
			img.Set(x*scale+i, y*scale+j, fillColor)
		}
	}
}

func getBitmapForChar(char rune) string {
	switch char {
	case '0':
		return Zero
	case '1':
		return One
	case '2':
		return Two
	case '3':
		return Three
	case '4':
		return Four
	case '5':
		return Five
	case '6':
		return Six
	case '7':
		return Seven
	case '8':
		return Eight
	case '9':
		return Nine
	default:
		return Colon
	}
}
