package render

import (
	"encoding/base64"
	"fmt"
	"io"
)

const kittyChunkSize = 3 * 4096 / 4

func KittyImage(output io.Writer, data []byte, columns, rows int) error {
	if len(data) == 0 {
		return fmt.Errorf("render Kitty image: empty PNG")
	}
	if columns < 1 || rows < 1 {
		return fmt.Errorf("render Kitty image: invalid dimensions %dx%d", columns, rows)
	}

	encoded := make([]byte, base64.StdEncoding.EncodedLen(kittyChunkSize))
	for offset := 0; offset < len(data); offset += kittyChunkSize {
		end := min(offset+kittyChunkSize, len(data))
		chunk := data[offset:end]
		more := end < len(data)

		control := fmt.Sprintf("q=2,m=%d", boolInt(more))
		if offset == 0 {
			control = fmt.Sprintf("a=T,f=100,t=d,c=%d,r=%d,%s", columns, rows, control)
		}
		if _, err := fmt.Fprintf(output, "\x1b_G%s;", control); err != nil {
			return fmt.Errorf("render Kitty image: %w", err)
		}

		n := base64.StdEncoding.EncodedLen(len(chunk))
		base64.StdEncoding.Encode(encoded[:n], chunk)
		if _, err := output.Write(encoded[:n]); err != nil {
			return fmt.Errorf("render Kitty image: %w", err)
		}
		if _, err := io.WriteString(output, "\x1b\\"); err != nil {
			return fmt.Errorf("render Kitty image: %w", err)
		}
	}

	if _, err := io.WriteString(output, "\r\n"); err != nil {
		return fmt.Errorf("render Kitty image: %w", err)
	}
	return nil
}

func ITermImage(output io.Writer, data []byte, columns, rows int) error {
	if len(data) == 0 {
		return fmt.Errorf("render iTerm image: empty PNG")
	}
	if columns < 1 || rows < 1 {
		return fmt.Errorf("render iTerm image: invalid dimensions %dx%d", columns, rows)
	}

	if _, err := fmt.Fprintf(output, "\x1b]1337;File=inline=1;size=%d;width=%d;height=%d:%s",
		len(data), columns, rows, base64.StdEncoding.EncodeToString(data)); err != nil {
		return fmt.Errorf("render iTerm image: %w", err)
	}
	if _, err := io.WriteString(output, "\a\r\n"); err != nil {
		return fmt.Errorf("render iTerm image: %w", err)
	}
	return nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}
