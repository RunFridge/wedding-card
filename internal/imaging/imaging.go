package imaging

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"io"

	_ "image/gif"
	_ "image/png"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

const (
	fullMaxDim      = 2048
	fullJPEGQuality = 85
	thumbMaxDim      = 400
	thumbJPEGQuality = 85
)

type ProcessedImage struct {
	Optimized []byte
	Thumbnail string
}

type ProcessedAsset struct {
	Optimized []byte
	Thumbnail []byte
}

func ProcessUpload(r io.Reader) (*ProcessedImage, error) {
	src, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	optimized := resizeToFit(src, fullMaxDim)
	var fullBuf bytes.Buffer
	if err := jpeg.Encode(&fullBuf, optimized, &jpeg.Options{Quality: fullJPEGQuality}); err != nil {
		return nil, fmt.Errorf("encode optimized: %w", err)
	}

	thumb := resizeToFit(src, thumbMaxDim)
	var thumbBuf bytes.Buffer
	if err := jpeg.Encode(&thumbBuf, thumb, &jpeg.Options{Quality: thumbJPEGQuality}); err != nil {
		return nil, fmt.Errorf("encode thumbnail: %w", err)
	}

	thumbDataURL := "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(thumbBuf.Bytes())

	return &ProcessedImage{
		Optimized: fullBuf.Bytes(),
		Thumbnail: thumbDataURL,
	}, nil
}

func ProcessAssetUpload(r io.Reader) (*ProcessedAsset, error) {
	src, _, err := image.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode image: %w", err)
	}

	optimized := resizeToFit(src, fullMaxDim)
	var fullBuf bytes.Buffer
	if err := jpeg.Encode(&fullBuf, optimized, &jpeg.Options{Quality: fullJPEGQuality}); err != nil {
		return nil, fmt.Errorf("encode optimized: %w", err)
	}

	thumb := resizeToFit(src, thumbMaxDim)
	var thumbBuf bytes.Buffer
	if err := jpeg.Encode(&thumbBuf, thumb, &jpeg.Options{Quality: thumbJPEGQuality}); err != nil {
		return nil, fmt.Errorf("encode thumbnail: %w", err)
	}

	return &ProcessedAsset{
		Optimized: fullBuf.Bytes(),
		Thumbnail: thumbBuf.Bytes(),
	}, nil
}

func resizeToFit(src image.Image, maxDim int) image.Image {
	bounds := src.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w <= maxDim && h <= maxDim {
		return src
	}

	var newW, newH int
	if w > h {
		newW = maxDim
		newH = h * maxDim / w
	} else {
		newH = maxDim
		newW = w * maxDim / h
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), src, bounds, draw.Over, nil)
	return dst
}
