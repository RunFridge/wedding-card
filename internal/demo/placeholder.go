package demo

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"strconv"

	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const (
	placeholderWidth       = 800
	placeholderHeight      = 1000
	placeholderJPEGQuality = 85
)

var pastelPalette = []color.RGBA{
	{R: 244, G: 200, B: 200, A: 255},
	{R: 247, G: 220, B: 180, A: 255},
	{R: 246, G: 240, B: 190, A: 255},
	{R: 205, G: 235, B: 200, A: 255},
	{R: 195, G: 225, B: 240, A: 255},
	{R: 205, G: 205, B: 240, A: 255},
	{R: 235, G: 205, B: 240, A: 255},
	{R: 240, G: 215, B: 225, A: 255},
}

func placeholderJPEG(idx int) ([]byte, error) {
	n := len(pastelPalette)
	bg := pastelPalette[idx%n]
	fg := pastelPalette[(idx+3+idx/n)%n]

	img := image.NewRGBA(image.Rect(0, 0, placeholderWidth, placeholderHeight))
	for y := 0; y < placeholderHeight; y++ {
		for x := 0; x < placeholderWidth; x++ {
			if patternHit(idx, x, y) {
				img.SetRGBA(x, y, fg)
			} else {
				img.SetRGBA(x, y, bg)
			}
		}
	}

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: placeholderJPEGQuality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

var digitColor = color.RGBA{R: 43, G: 29, B: 14, A: 255}

func numberJPEG(n int) ([]byte, error) {
	const smallW, smallH = 16, 20
	bg := pastelPalette[n%len(pastelPalette)]

	small := image.NewRGBA(image.Rect(0, 0, smallW, smallH))
	draw.Draw(small, small.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)

	face := basicfont.Face7x13
	label := strconv.Itoa(n)
	labelWidth := font.MeasureString(face, label).Ceil()
	drawer := &font.Drawer{
		Dst:  small,
		Src:  image.NewUniform(digitColor),
		Face: face,
		Dot:  fixed.P((smallW-labelWidth)/2, (smallH-face.Height)/2+face.Ascent),
	}
	drawer.DrawString(label)

	img := image.NewRGBA(image.Rect(0, 0, placeholderWidth, placeholderHeight))
	xdraw.NearestNeighbor.Scale(img, img.Bounds(), small, small.Bounds(), xdraw.Src, nil)

	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: placeholderJPEGQuality}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func patternHit(idx, x, y int) bool {
	cx := x - placeholderWidth/2
	cy := y - placeholderHeight/2
	outerRadius := placeholderWidth / 3
	innerRadius := placeholderWidth / 5
	stripeWidth := placeholderWidth / 10
	checkerSize := placeholderWidth / 8

	switch idx % 4 {
	case 0:
		return cx*cx+cy*cy < outerRadius*outerRadius
	case 1:
		return (x+y)/stripeWidth%2 == 0
	case 2:
		return (x/checkerSize+y/checkerSize)%2 == 0
	default:
		d := cx*cx + cy*cy
		return d > innerRadius*innerRadius && d < outerRadius*outerRadius
	}
}
