package render

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/png"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestBrailleImageDimensions(t *testing.T) {
	data := encodePNG(t, image.NewUniform(color.White), image.Rect(0, 0, 96, 96))
	var output bytes.Buffer
	if err := BrailleImage(&output, data, 48, 24); err != nil {
		t.Fatal(err)
	}

	lines := strings.Split(strings.TrimSuffix(output.String(), "\n"), "\n")
	if len(lines) != 24 {
		t.Fatalf("line count = %d, want 24", len(lines))
	}
	for _, line := range lines {
		if utf8.RuneCountInString(line) != 48 {
			t.Fatalf("line width = %d, want 48", utf8.RuneCountInString(line))
		}
		if !utf8.ValidString(line) || strings.ContainsRune(line, '\x1b') {
			t.Fatalf("invalid terminal output: %q", line)
		}
	}
}

func TestBrailleImageTransparencyAndDots(t *testing.T) {
	source := image.NewNRGBA(image.Rect(0, 0, 2, 4))
	source.SetNRGBA(0, 0, color.NRGBA{R: 255, G: 255, B: 255, A: 255})
	data := encodePNG(t, source, source.Bounds())

	var output bytes.Buffer
	if err := BrailleImage(&output, data, 1, 1); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), "⠁\n"; got != want {
		t.Fatalf("BrailleImage() = %q, want %q", got, want)
	}

	transparent := encodePNG(t, image.NewNRGBA(image.Rect(0, 0, 2, 4)), image.Rect(0, 0, 2, 4))
	output.Reset()
	if err := BrailleImage(&output, transparent, 1, 1); err != nil {
		t.Fatal(err)
	}
	if got, want := output.String(), " \n"; got != want {
		t.Fatalf("transparent BrailleImage() = %q, want %q", got, want)
	}
}

func TestBrailleImageErrors(t *testing.T) {
	if err := BrailleImage(&bytes.Buffer{}, []byte("not a PNG"), 48, 24); err == nil {
		t.Fatal("BrailleImage() accepted an invalid PNG")
	}
	if err := BrailleImage(&bytes.Buffer{}, nil, 0, 24); err == nil {
		t.Fatal("BrailleImage() accepted invalid dimensions")
	}

	data := encodePNG(t, image.NewUniform(color.White), image.Rect(0, 0, 2, 4))
	if err := BrailleImage(&errorWriter{}, data, 1, 1); !errors.Is(err, errWrite) {
		t.Fatalf("BrailleImage() error = %v, want %v", err, errWrite)
	}
}

func encodePNG(t *testing.T, source image.Image, bounds image.Rectangle) []byte {
	t.Helper()
	var output bytes.Buffer
	if err := png.Encode(&output, sourceWithBounds{Image: source, bounds: bounds}); err != nil {
		t.Fatal(err)
	}
	return output.Bytes()
}

type sourceWithBounds struct {
	image.Image
	bounds image.Rectangle
}

func (source sourceWithBounds) Bounds() image.Rectangle { return source.bounds }

var errWrite = errors.New("write failed")

type errorWriter struct{}

func (*errorWriter) Write([]byte) (int, error) { return 0, errWrite }
