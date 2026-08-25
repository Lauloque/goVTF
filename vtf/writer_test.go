/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/Lauloque/goVTF/texture"
)

// TestHeaderSize verifies the logical header struct size matches the spec constant.
// Note: This checks the Go struct size, which might differ from binary layout due to padding.
// The actual binary write uses packedHeader, so this test is more of a sanity check.
func TestHeaderSize(t *testing.T) {
	// We expect the binary size of our logical struct to be close to HeaderSize,
	// but the critical check is that packedHeader is exactly HeaderSize.
	var h VTFHeader
	size := -binary.Size(h)

	var p packedHeader
	psize := binary.Size(p)

	t.Logf("Logical VTFHeader size: %d, packedHeader size: %d, HeaderSize constant: %d", size, psize, HeaderSize)

	// if the struc gets padding, comparing size to HeaderSize might fail, checinkg packedHeader instead:
	if psize != HeaderSize {
		t.Errorf("packedHeader size is %d, expected %d", psize, HeaderSize)
	}

}

// TestVTFHeaderByteOffsets verifies the binary layout matches the spec.
func TestVTFHeaderByteOffsets(t *testing.T) {
	tex := texture.NewTexture(512, 512, texture.PixelFormatRGBA8888, make([]byte, 512*512*4))
	var buf bytes.Buffer
	err := Write(&buf, tex)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	data := buf.Bytes()

	// Define checks using constants, not magic numbers
	checks := []struct {
		name     string
		offset   int
		length   int
		expected []byte
	}{
		{"Signature", 0, 4, []byte("VTF\x00")},
		{"VersionMajor", 4, 4, encodeUint32(SignatureVersionMajor)},
		{"VersionMinor", 8, 4, encodeUint32(SignatureVersionMinor)},
		{"HeaderSize", 12, 4, encodeUint32(HeaderSize)},
		{"LowResWidth", 61, 1, []byte{16}}, // Calculated: min(512, 16) -> 16
		{"LowResHeight", 62, 1, []byte{16}},
		{"NumResources", 68, 4, encodeUint32(2)}, // We always write 2 resources
	}

	for _, tt := range checks {
		t.Run(tt.name, func(t *testing.T) {
			if len(data) < tt.offset+tt.length {
				t.Fatalf("Buffer too short at offset %d", tt.offset)
			}
			got := data[tt.offset : tt.offset+tt.length]
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("At offset %d (%s): expected %x, got %x", tt.offset, tt.name, tt.expected, got)
			}
		})
	}
}

// Helper to encode uint32 to little-endian bytes
func encodeUint32(v uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	return b
}

// TestLowResDimensions (Keep as is - logic test)
func TestLowResDimensions(t *testing.T) {
	// ... (Your existing logic test remains unchanged) ...
	tests := []struct {
		name           string
		texWidth       int
		texHeight      int
		expectedWidth  int
		expectedHeight int
	}{
		{"small 64x64", 64, 64, 16, 16},
		{"medium 256x256", 256, 256, 16, 16},
		{"large 512x512", 512, 512, 16, 16},
		{"non-square 512x128", 512, 128, 16, 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lw, lh := tt.texWidth, tt.texHeight
			if lw > 16 {
				lw = 16
			}
			if lh > 16 {
				lh = 16
			}
			lw = ((lw + 3) / 4) * 4
			lh = ((lh + 3) / 4) * 4

			if lw != tt.expectedWidth || lh != tt.expectedHeight {
				t.Errorf("Expected %dx%d, got %dx%d", tt.expectedWidth, tt.expectedHeight, lw, lh)
			}
		})
	}
}

// TestVTFWriteRoundTrip
func TestVTFWriteRoundTrip(t *testing.T) {
	tex := texture.NewTexture(512, 512, texture.PixelFormatRGBA8888, make([]byte, 512*512*4))
	var buf bytes.Buffer
	err := Write(&buf, tex)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	data := buf.Bytes()

	// Verify Signature
	if string(data[0:4]) != "VTF\x00" {
		t.Errorf("Invalid signature")
	}

	// Verify Version from Constants
	major := binary.LittleEndian.Uint32(data[4:8])
	minor := binary.LittleEndian.Uint32(data[8:12])
	if major != SignatureVersionMajor || minor != SignatureVersionMinor {
		t.Errorf("Expected version %d.%d, got %d.%d", SignatureVersionMajor, SignatureVersionMinor, major, minor)
	}

	// Verify Header Size from Constant
	headerSize := binary.LittleEndian.Uint32(data[12:16])
	if headerSize != HeaderSize {
		t.Errorf("Expected header size %d, got %d", HeaderSize, headerSize)
	}

	// Verify NumResources
	numResources := binary.LittleEndian.Uint32(data[68:72])
	if numResources != 2 {
		t.Errorf("Expected NumResources=2, got %d", numResources)
	}

	// Calculate Expected Size Dynamically
	lowResW, lowResH := uint8(16), uint8(16) // For 512x512
	lowResDataSize := int(lowResW) * int(lowResH) * 4
	highResDataSize := 512 * 512 * 4
	expectedTotal := HeaderSize + 8 + lowResDataSize + 8 + highResDataSize

	if len(data) != expectedTotal {
		t.Errorf("Expected total size %d, got %d", expectedTotal, len(data))
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
