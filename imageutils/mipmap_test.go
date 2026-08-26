/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

import (
	"testing"

	"github.com/Lauloque/goVTF/texture"
)

func TestCountMipmaps(t *testing.T) {
	tests := []struct {
		w, h, want int
	}{
		{512, 512, 10}, // 512,256,128,64,32,16,8,4,2,1
		{64, 64, 7},    // 64,32,16,8,4,2,1
		{4, 4, 3},      // 4,2,1
	}
	for _, tt := range tests {
		if got := CountMipmaps(tt.w, tt.h); got != tt.want {
			t.Errorf("CountMipmaps(%d,%d) = %d, want %d", tt.w, tt.h, got, tt.want)
		}
	}
}

// TestCountMipmapsMatchesGeneratedLevels locks in the invariant that
// header.MipmapCount and the number of levels actually written to the file
// must agree — this exact mismatch previously caused an infinite loop and
// (before that) misread texture data.
func TestCountMipmapsMatchesGeneratedLevels(t *testing.T) {
	tex := texture.NewTexture(512, 512, texture.PixelFormatRGBA8888, make([]byte, 512*512*4))
	levels, err := GenerateMipmaps(tex)
	if err != nil {
		t.Fatalf("GenerateMipmaps failed: %v", err)
	}

	want := CountMipmaps(tex.Width, tex.Height)
	got := len(levels) + 1 // +1 for the full-res level, written separately
	if got != want {
		t.Errorf("CountMipmaps() = %d, but len(GenerateMipmaps())+1 = %d", want, got)
	}
}

func TestGenerateMipmapsOrdering(t *testing.T) {
	tex := texture.NewTexture(64, 64, texture.PixelFormatRGBA8888, make([]byte, 64*64*4))
	levels, err := GenerateMipmaps(tex)
	if err != nil {
		t.Fatalf("GenerateMipmaps failed: %v", err)
	}
	if len(levels) < 2 {
		t.Fatalf("expected multiple mip levels, got %d", len(levels))
	}
	for i := 1; i < len(levels); i++ {
		if len(levels[i]) < len(levels[i-1]) {
			t.Errorf("level %d (%d bytes) is smaller than level %d (%d bytes); expected smallest-to-largest order",
				i, len(levels[i]), i-1, len(levels[i-1]))
		}
	}
}

func TestGenerateMipmapsRejectsOddDimensions(t *testing.T) {
	// Guards the downSampleHalf error path even though ValidateDimensions
	// currently prevents this from occurring via Write.
	tex := texture.NewTexture(3, 3, texture.PixelFormatRGBA8888, make([]byte, 3*3*4))
	if _, err := GenerateMipmaps(tex); err == nil {
		t.Error("expected error for odd dimensions, got nil")
	}
}
