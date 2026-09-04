package render

import (
	"bytes"
	"encoding/base64"
	"strings"
	"testing"
)

func TestKittyImage(t *testing.T) {
	data := bytes.Repeat([]byte{0x89, 0x50, 0x4e, 0x47}, kittyChunkSize/4+2)
	var output bytes.Buffer
	if err := KittyImage(&output, data, 48, 24); err != nil {
		t.Fatal(err)
	}

	commands := strings.Split(output.String(), "\x1b_G")[1:]
	if len(commands) != 2 || !strings.HasPrefix(commands[0], "a=T,f=100,t=d,c=48,r=24,C=1,q=2,m=1;") ||
		!strings.HasPrefix(commands[1], "q=2,m=0;") {
		t.Fatalf("unexpected Kitty commands: %q", commands)
	}

	var decoded []byte
	for _, command := range commands {
		payload := strings.TrimSuffix(strings.SplitN(command, ";", 2)[1], "\x1b\\")
		if len(payload) > 4096 {
			t.Fatalf("payload length = %d, want at most 4096", len(payload))
		}
		chunk, err := base64.StdEncoding.DecodeString(payload)
		if err != nil {
			t.Fatal(err)
		}
		decoded = append(decoded, chunk...)
	}
	if !bytes.Equal(decoded, data) {
		t.Fatal("decoded Kitty payload does not match the source")
	}
}

func TestITermImage(t *testing.T) {
	var output bytes.Buffer
	if err := ITermImage(&output, []byte("png"), 48, 24); err != nil {
		t.Fatal(err)
	}
	want := "\x1b]1337;File=inline=1;doNotMoveCursor=1;size=3;width=48;height=24:cG5n\a"
	if output.String() != want {
		t.Fatalf("ITermImage() = %q, want %q", output.String(), want)
	}
}

func TestRasterImageErrors(t *testing.T) {
	if err := KittyImage(&errorWriter{}, []byte("png"), 48, 24); err == nil {
		t.Fatal("KittyImage() did not propagate a write error")
	}
	if err := ITermImage(&errorWriter{}, []byte("png"), 48, 24); err == nil {
		t.Fatal("ITermImage() did not propagate a write error")
	}
	if err := KittyImage(&bytes.Buffer{}, nil, 48, 24); err == nil {
		t.Fatal("KittyImage() accepted empty data")
	}
}
