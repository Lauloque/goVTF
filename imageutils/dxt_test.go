/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

import (
	"testing"

	"github.com/Lauloque/goVTF/texture"
)

func decodeDXT5Alpha(block []byte, pixel int) byte {
	palette := [8]byte{block[0], block[1]}
	if palette[0] > palette[1] {
		for i := 1; i <= 6; i++ {
			palette[i+1] = byte(((7-i)*int(palette[0]) + i*int(palette[1])) / 7)
		}
	} else {
		for i := 1; i <= 4; i++ {
			palette[i+1] = byte(((5-i)*int(palette[0]) + i*int(palette[1])) / 5)
		}
		palette[6], palette[7] = 0, 255
	}

	var indices uint64
	for i := 0; i < 6; i++ {
		indices |= uint64(block[2+i]) << (8 * i)
	}
	return palette[(indices>>uint(3*pixel))&7]
}

func TestCompressDXT5PreservesAlpha(t *testing.T) {
	pixels := make([]byte, 4*4*4)
	for i := 0; i < 16; i++ {
		pixels[i*4] = 255
		pixels[i*4+3] = byte(i * 17)
	}
	tex := texture.NewTexture(4, 4, texture.PixelFormatRGBA8888, pixels)
	compressed := CompressDXT5(tex)

	if len(compressed) != 16 {
		t.Fatalf("DXT5 block length = %d, want 16", len(compressed))
	}
	for i := 0; i < 16; i++ {
		got := decodeDXT5Alpha(compressed[:8], i)
		want := pixels[i*4+3]
		difference := int(got) - int(want)
		if difference < 0 {
			difference = -difference
		}
		if difference > 18 {
			t.Errorf("pixel %d alpha = %d, want approximately %d", i, got, want)
		}
	}
}

func TestCompressDXT5ConstantAlpha(t *testing.T) {
	for _, alpha := range []byte{0, 127, 255} {
		pixels := make([]byte, 4*4*4)
		for i := 0; i < 16; i++ {
			pixels[i*4+3] = alpha
		}
		compressed := CompressDXT5(texture.NewTexture(4, 4, texture.PixelFormatRGBA8888, pixels))
		for i := 0; i < 16; i++ {
			if got := decodeDXT5Alpha(compressed[:8], i); got != alpha {
				t.Errorf("alpha %d, pixel %d decoded as %d", alpha, i, got)
			}
		}
	}
}
