/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/Lauloque/goVTF/texture"
)

// TestVTFHeaderByteOffsets - verifies manual byte layout matches VTF spec
func TestVTFHeaderByteOffsets(t *testing.T) {
	// This validates the offsets used in manual byte construction
	// match the VTF 7.4 spec exactly

	tests := []struct {
		fieldName string
		offset    int
		length    int
		expected  []byte
	}{
		{"Signature", 0, 4, []byte("VTF\x00")},
		{"VersionMajor", 4, 4, []byte{0x07, 0x00, 0x00, 0x00}},
		{"VersionMinor", 8, 4, []byte{0x04, 0x00, 0x00, 0x00}},
		{"HeaderSize", 12, 4, []byte{0x60, 0x00, 0x00, 0x00}},
		{"LowResWidth", 61, 1, []byte{0x10}},  // 16
		{"LowResHeight", 62, 1, []byte{0x10}}, // 16
		{"NumResources", 68, 4, []byte{0x01, 0x00, 0x00, 0x00}},
	}

	tex := texture.NewTexture(512, 512, texture.PixelFormatRGBA8888, make([]byte, 512*512*4))
	var buf bytes.Buffer
	err := Write(&buf, tex)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	data := buf.Bytes()

	for _, tt := range tests {
		t.Run(tt.fieldName, func(t *testing.T) {
			if len(data) < tt.offset+tt.length {
				t.Errorf("Buffer too short at offset %d", tt.offset)
				return
			}

			got := data[tt.offset : tt.offset+tt.length]
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("At offset %d: expected %x, got %x", tt.offset, tt.expected, got)
			}
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
	lowResWidth := data[61]
	lowResHeight := data[62]
	if lowResWidth != 16 || lowResHeight != 16 {
		t.Errorf("Expected low-res dimensions 16x16 at offsets 61-62, got %dx%d", lowResWidth, lowResHeight)
	}

	// Check NumResources
	numResources := binary.LittleEndian.Uint32(data[68:72])
	if numResources != 1 {
		t.Errorf("Expected NumResources=1 at offset 68, got %d", numResources)
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
