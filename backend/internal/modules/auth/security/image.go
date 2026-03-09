package security

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

// renderCaptchaImage generates a PNG captcha image and returns it as a data URI.
func renderCaptchaImage(code string) (string, error) {
	const (
		width  = 120
		height = 42
	)

	// Create image
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// Fill background with light gray
	draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{240, 240, 240, 255}}, image.Point{}, draw.Src)

	// Add noise dots
	for i := 0; i < 100; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		img.Set(x, y, randomGray())
	}

	// Add noise lines
	for i := 0; i < 4; i++ {
		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)
		drawLine(img, x1, y1, x2, y2, randomGray())
	}

	// Draw characters
	f := basicfont.Face7x13
	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.Black,
		Face: f,
		Dot:  fixed.Point26_6{},
	}

	charWidth := width / len(code)
	for i, c := range code {
		// Random vertical offset
		yOffset := rand.Intn(8) + 14
		drawer.Dot = fixed.Point26_6{
			X: fixed.Int26_6((i*charWidth + 8) * 64),
			Y: fixed.Int26_6(yOffset * 64),
		}
		drawer.DrawString(string(c))
	}

	// Encode to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return "", err
	}

	// Return as data URI
	return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func randomGray() color.Color {
	v := rand.Intn(128) + 64
	return color.RGBA{uint8(v), uint8(v), uint8(v), 255}
}

func drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx, sy := 1, 1
	if x1 > x2 {
		sx = -1
	}
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy

	for {
		img.Set(x1, y1, c)
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
