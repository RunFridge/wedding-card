package imaging

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"strings"
	"testing"
)

func gradientImage(w, h int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{
				R: uint8(x * 255 / w),
				G: uint8(y * 255 / h),
				B: 128,
				A: 255,
			})
		}
	}
	return img
}

func encodeJPEG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("encode JPEG: %v", err)
	}
	return buf.Bytes()
}

func encodePNG(t *testing.T, img image.Image) []byte {
	t.Helper()
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode PNG: %v", err)
	}
	return buf.Bytes()
}

func TestProcessUploadJPEG(t *testing.T) {
	src := gradientImage(4000, 3000)
	data := encodeJPEG(t, src)

	result, err := ProcessUpload(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ProcessUpload: %v", err)
	}

	decoded, err := jpeg.Decode(bytes.NewReader(result.Optimized))
	if err != nil {
		t.Fatalf("decode optimized: %v", err)
	}

	bounds := decoded.Bounds()
	if bounds.Dx() != fullMaxDim {
		t.Errorf("optimized width = %d, want %d", bounds.Dx(), fullMaxDim)
	}
	expectedH := 3000 * fullMaxDim / 4000
	if bounds.Dy() != expectedH {
		t.Errorf("optimized height = %d, want %d", bounds.Dy(), expectedH)
	}

	if !strings.HasPrefix(result.Thumbnail, "data:image/jpeg;base64,") {
		t.Error("thumbnail missing data URL prefix")
	}

	thumbData := result.Thumbnail[len("data:image/jpeg;base64,"):]
	if len(thumbData) == 0 {
		t.Error("thumbnail base64 data is empty")
	}
}

func TestProcessUploadPNG(t *testing.T) {
	src := gradientImage(1000, 800)
	data := encodePNG(t, src)

	result, err := ProcessUpload(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ProcessUpload: %v", err)
	}

	decoded, err := jpeg.Decode(bytes.NewReader(result.Optimized))
	if err != nil {
		t.Fatalf("decode optimized (should be JPEG): %v", err)
	}

	bounds := decoded.Bounds()
	if bounds.Dx() != 1000 || bounds.Dy() != 800 {
		t.Errorf("dimensions = %dx%d, want 1000x800 (no resize needed)", bounds.Dx(), bounds.Dy())
	}
}

func TestProcessUploadSmallImage(t *testing.T) {
	src := gradientImage(200, 150)
	data := encodeJPEG(t, src)

	result, err := ProcessUpload(bytes.NewReader(data))
	if err != nil {
		t.Fatalf("ProcessUpload: %v", err)
	}

	decoded, err := jpeg.Decode(bytes.NewReader(result.Optimized))
	if err != nil {
		t.Fatalf("decode optimized: %v", err)
	}

	bounds := decoded.Bounds()
	if bounds.Dx() != 200 || bounds.Dy() != 150 {
		t.Errorf("dimensions = %dx%d, want 200x150 (no resize)", bounds.Dx(), bounds.Dy())
	}
}

func TestProcessUploadInvalidFormat(t *testing.T) {
	_, err := ProcessUpload(bytes.NewReader([]byte("not an image")))
	if err == nil {
		t.Fatal("expected error for invalid format")
	}
}

func TestResizeToFitLandscape(t *testing.T) {
	src := gradientImage(4000, 2000)
	result := resizeToFit(src, 2048)

	bounds := result.Bounds()
	if bounds.Dx() != 2048 {
		t.Errorf("width = %d, want 2048", bounds.Dx())
	}
	expectedH := 2000 * 2048 / 4000
	if bounds.Dy() != expectedH {
		t.Errorf("height = %d, want %d", bounds.Dy(), expectedH)
	}
}

func TestResizeToFitPortrait(t *testing.T) {
	src := gradientImage(2000, 4000)
	result := resizeToFit(src, 2048)

	bounds := result.Bounds()
	if bounds.Dy() != 2048 {
		t.Errorf("height = %d, want 2048", bounds.Dy())
	}
	expectedW := 2000 * 2048 / 4000
	if bounds.Dx() != expectedW {
		t.Errorf("width = %d, want %d", bounds.Dx(), expectedW)
	}
}

func TestResizeToFitNoOpWhenSmall(t *testing.T) {
	src := gradientImage(500, 400)
	result := resizeToFit(src, 2048)

	if result != src {
		t.Error("expected same image pointer when no resize needed")
	}
}
