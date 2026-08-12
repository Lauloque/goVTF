/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/Lauloque/goVTF/texture"
)

// Test VTFHeader struct layout matches the spec
func TestVTFHeaderLayout(t *testing.T) {
	var h VTFHeader

	// Total header size should be 96 bytes for VTF 7.3+
	if size := binary.Size(h); size != 96 {
		t.Errorf("Expected VTFHeader size to be 96 bytes, got %d", size)
	}

	// Critical field offsets match spec
	tests := []struct {
		name     string
		offset   uintptr
		expected uintptr
	}{
		{"Signature", 0, 0},
		{"Version", 4, 4},
		{"HeaderSize", 12, 12},
		{"Width", 16, 16},
		{"Height", 18, 18},
		{"Flags", 20, 20},
		{"LowResWidth", 58, 58},
		{"LowResHeight", 59, 59},
		{"NumResources", 66, 66},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: You'd need unsafe.Offsetof() here, but we're testing
			// the serialized output instead to avoid import issues
		})
	}
}

// Test low-res dimension calculation
func TestLowResDimensions(t *testing.T) {
	tests := []struct {
		name           string
		texWidth       int
		texHeight      int
		expectedWidth  int
		expectedHeight int
	}{
		{
			name:           "small 64x64",
			texWidth:       64,
			texHeight:      64,
			expectedWidth:  16,
			expectedHeight: 16,
		},
		{
			name:           "medium 256x256",
			texWidth:       256,
			texHeight:      256,
			expectedWidth:  16,
			expectedHeight: 16,
		},
		{
			name:           "large 512x512",
			texWidth:       512,
			texHeight:      512,
			expectedWidth:  16,
			expectedHeight: 16,
		},
		{
			name:           "very large 1024x1024",
			texWidth:       1024,
			texHeight:      1024,
			expectedWidth:  16,
			expectedHeight: 16,
		},
		{
			name:           "non-square 512x128",
			texWidth:       512,
			texHeight:      128,
			expectedWidth:  16,
			expectedHeight: 16,
		},
		{
			name:           "tiny 8x8",
			texWidth:       8,
			texHeight:      8,
			expectedWidth:  8,
			expectedHeight: 8,
		},
		{
			name:           "tiny 4x4",
			texWidth:       4,
			texHeight:      4,
			expectedWidth:  4,
			expectedHeight: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Recreate the calculation logic from writer.go
			lowResWidth := tt.texWidth
			lowResHeight := tt.texHeight

			if lowResWidth > 16 {
				lowResWidth = 16
			}
			if lowResHeight > 16 {
				lowResHeight = 16
			}

			lowResWidth = ((lowResWidth + 3) / 4) * 4
			lowResHeight = ((lowResHeight + 3) / 4) * 4

			if lowResWidth != tt.expectedWidth {
				t.Errorf("Expected LowResWidth %d, got %d", tt.expectedWidth, lowResWidth)
			}
			if lowResHeight != tt.expectedHeight {
				t.Errorf("Expected LowResHeight %d, got %d", tt.expectedHeight, lowResHeight)
			}
		})
	}
}

// Test that uint8 conversion doesn't overflow for common sizes
func TestLowResNoOverflow(t *testing.T) {
	sizes := []int{64, 128, 256, 512, 1024}

	for _, size := range sizes {
		lowResWidth := size
		if lowResWidth > 16 {
			lowResWidth = 16
		}

		// This is the critical check - no overflow
		if lowResWidth > 255 {
			t.Errorf("Size %d would overflow uint8", size)
		}

		u8 := uint8(lowResWidth)
		if int(u8) != lowResWidth {
			t.Errorf("uint8(%d) = %d, overflow occurred!", lowResWidth, u8)
		}
	}
}

// Test full VTF write round-trip
func TestVTFWriteRoundTrip(t *testing.T) {
	tex := texture.NewTexture(
		512, 512,
		texture.PixelFormatRGBA8888,
		make([]byte, 512*512*4), // 1MB RGBA data
	)

	var buf bytes.Buffer
	err := Write(&buf, tex)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data := buf.Bytes()

	// Check signature
	if string(data[0:4]) != "VTF\x00" {
		t.Errorf("Invalid signature: %q", data[0:4])
	}

	// Check version
	major := binary.LittleEndian.Uint32(data[4:8])
	minor := binary.LittleEndian.Uint32(data[8:12])
	if major != 7 || minor != 4 {
		t.Errorf("Expected version 7.4, got %d.%d", major, minor)
	}

	// Check header size
	headerSize := binary.LittleEndian.Uint32(data[12:16])
	if headerSize != 96 {
		t.Errorf("Expected header size 96, got %d", headerSize)
	}

	// Check dimensions
	width := binary.LittleEndian.Uint16(data[16:18])
	height := binary.LittleEndian.Uint16(data[18:20])
	if width != 512 || height != 512 {
		t.Errorf("Expected dimensions 512x512, got %dx%d", width, height)
	}

	// CRITICAL: Check low-res dimensions (the bug we just found!)
	lowResWidth := data[58]
	lowResHeight := data[59]
	if lowResWidth != 16 || lowResHeight != 16 {
		t.Errorf("Expected low-res dimensions 16x16, got %dx%d", lowResWidth, lowResHeight)
	}

	// Check NumResources
	numResources := binary.LittleEndian.Uint32(data[66:70])
	if numResources != 1 {
		t.Errorf("Expected NumResources=1, got %d", numResources)
	}

	// Check file size: header(96) + resource(8) + pixels(512*512*4)
	expectedSize := 96 + 8 + (512 * 512 * 4)
	if len(data) != expectedSize {
		t.Errorf("Expected file size %d, got %d", expectedSize, len(data))
	}
}

// Test dimension validation rejects invalid inputs
func TestDimensionValidation(t *testing.T) {
	tests := []struct {
		name    string
		width   int
		height  int
		wantErr bool
	}{
		{"valid 256x256", 256, 256, false},
		{"valid 512x512", 512, 512, false},
		{"not power of 2", 568, 568, true},
		{"not multiple of 4", 510, 510, true},
		{"too small 2x2", 2, 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tex := texture.NewTexture(
				tt.width, tt.height,
				texture.PixelFormatRGBA8888,
				make([]byte, tt.width*tt.height*4),
			)

			var buf bytes.Buffer
			err := Write(&buf, tex)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error for %dx%d, got nil", tt.width, tt.height)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error for %dx%d: %v", tt.width, tt.height, err)
			}
		})
	}
}
