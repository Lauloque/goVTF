/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

import (
	"testing"

	"github.com/Lauloque/goVTF/texture"
)

func TestResizeDimensions(t *testing.T) {
	src := make([]byte, 8*8*4)
	out, err := Resize(src, 8, 8, 4, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 4*4*4 {
		t.Errorf("output length = %d, want %d", len(out), 4*4*4)
	}
}

func TestResizeRejectsInvalidDimensions(t *testing.T) {
	tests := []struct {
		name                   string
		srcW, srcH, dstW, dstH int
	}{
		{"zero dst width", 8, 8, 0, 4},
		{"negative src height", 8, -1, 4, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Resize(make([]byte, 8*8*4), tt.srcW, tt.srcH, tt.dstW, tt.dstH); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestResizeRejectsUndersizedBuffer(t *testing.T) {
	tooSmall := make([]byte, 4) // way less than 8*8*4
	if _, err := Resize(tooSmall, 8, 8, 4, 4); err == nil {
		t.Error("expected error for undersized pixel buffer, got nil")
	}
}

func TestDownSampleHalf(t *testing.T) {
	// 2x2 image, single 2x2 block: (0,0,0),(10,10,10),(20,20,20),(30,30,30) -> avg (15,15,15)
	pixels := []byte{
		0, 0, 0, 255,
		10, 10, 10, 255,
		20, 20, 20, 255,
		30, 30, 30, 255,
	}
	tex := texture.NewTexture(2, 2, texture.PixelFormatRGBA8888, pixels)
	out, err := downSampleHalf(tex)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Width != 1 || out.Height != 1 {
		t.Fatalf("expected 1x1 output, got %dx%d", out.Width, out.Height)
	}
	want := byte(15)
	if out.Pixels[0] != want || out.Pixels[1] != want || out.Pixels[2] != want {
		t.Errorf("averaged pixel = %v, want (%d,%d,%d,_)", out.Pixels[:4], want, want, want)
	}
}

func TestDownSampleHalfRejectsOddDimensions(t *testing.T) {
	tex := texture.NewTexture(3, 4, texture.PixelFormatRGBA8888, make([]byte, 3*4*4))
	if _, err := downSampleHalf(tex); err == nil {
		t.Error("expected error for odd width, got nil")
	}
}
