package render

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"unicode/utf8"
)

var brailleDots = [4][2]byte{
	{1 << 0, 1 << 3},
	{1 << 1, 1 << 4},
	{1 << 2, 1 << 5},
	{1 << 6, 1 << 7},
}

var bayer = [4][4]uint64{
	{0, 8, 2, 10},
	{12, 4, 14, 6},
	{3, 11, 1, 9},
	{15, 7, 13, 5},
}

func BrailleImage(output io.Writer, data []byte, columns, rows int) error {
	lines, err := BrailleLines(data, columns, rows)
	if err != nil {
		return err
	}
	for _, line := range lines {
		if _, err := fmt.Fprintln(output, line); err != nil {
			return fmt.Errorf("render Braille image: %w", err)
		}
	}
	return nil
}

func BrailleLines(data []byte, columns, rows int) ([]string, error) {
	if columns < 1 || rows < 1 {
		return nil, fmt.Errorf("render Braille image: invalid dimensions %dx%d", columns, rows)
	}

	image, err := png.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("render Braille image: %w", err)
	}

	pixelWidth, pixelHeight := columns*2, rows*4
	bounds := image.Bounds()
	values := make([]uint64, pixelWidth*pixelHeight)
	var brightest uint64
	for y := range pixelHeight {
		y0 := bounds.Min.Y + y*bounds.Dy()/pixelHeight
		y1 := bounds.Min.Y + (y+1)*bounds.Dy()/pixelHeight
		for x := range pixelWidth {
			x0 := bounds.Min.X + x*bounds.Dx()/pixelWidth
			x1 := bounds.Min.X + (x+1)*bounds.Dx()/pixelWidth

			var luminance, pixels uint64
			for sourceY := y0; sourceY < y1; sourceY++ {
				for sourceX := x0; sourceX < x1; sourceX++ {
					r, g, b, _ := image.At(sourceX, sourceY).RGBA()
					luminance += (299*uint64(r) + 587*uint64(g) + 114*uint64(b)) / 1000
					pixels++
				}
			}
			if pixels > 0 {
				values[y*pixelWidth+x] = luminance / pixels
			}
			brightest = max(brightest, values[y*pixelWidth+x])
		}
	}

	lines := make([]string, 0, rows)
	for row := range rows {
		line := make([]byte, 0, columns*3)
		for column := range columns {
			var dots byte
			for y := range 4 {
				for x := range 2 {
					pixelX, pixelY := column*2+x, row*4+y
					value := values[pixelY*pixelWidth+pixelX]
					if brightest > 0 && value > 0 && value*16 > bayer[pixelY%4][pixelX%4]*brightest {
						dots |= brailleDots[y][x]
					}
				}
			}
			if dots == 0 {
				line = append(line, ' ')
			} else {
				line = utf8.AppendRune(line, rune(0x2800)+rune(dots))
			}
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}
