/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/Lauloque/goVTF/texture"
)

const testNumResources = 2 // low-res + high-res, the only two this writer emits

// dxt1Size computes the exact DXT1 byte size for a w×h image, independent of
// production code — this is what a compliant DXT1 encoder must output.
func dxt1Size(w, h int) int {
	tw := (w + 3) / 4
	th := (h + 3) / 4
	return tw * th * 8
}

func encodeUint32(v uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, v)
	return b
}

func TestHeaderSize(t *testing.T) {
	var p packedHeader
	if got := binary.Size(p); got != HeaderSize {
		t.Errorf("packedHeader size is %d, expected %d", got, HeaderSize)
	}
}

func TestImageFormatIsDXT1(t *testing.T) {
	if ImageFormatDXT1 != 13 {
		t.Errorf("ImageFormatDXT1 = %d, want 13 (per Source engine imageformat.h)", ImageFormatDXT1)
	}
}

func TestVTFHeaderByteOffsets(t *testing.T) {
	tex := texture.NewTexture(512, 512, texture.PixelFormatRGBA8888, make([]byte, 512*512*4))
	var buf bytes.Buffer
	if err := Write(&buf, tex); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	data := buf.Bytes()
	dynamicHeaderSize := uint32(HeaderSize) + uint32(testNumResources*8)

	checks := []struct {
		name     string
		offset   int
		length   int
		expected []byte
	}{
		{"Signature", 0, 4, []byte("VTF\x00")},
		{"VersionMajor", 4, 4, encodeUint32(SignatureVersionMajor)},
		{"VersionMinor", 8, 4, encodeUint32(SignatureVersionMinor)},
		{"HeaderSize", 12, 4, encodeUint32(dynamicHeaderSize)},
		{"HighResFormat", 52, 4, encodeUint32(uint32(ImageFormatDXT1))},
		{"LowResFormat", 57, 4, encodeUint32(uint32(ImageFormatDXT1))},
		{"LowResWidth", 61, 1, []byte{16}},
		{"LowResHeight", 62, 1, []byte{16}},
		{"NumResources", 68, 4, encodeUint32(testNumResources)},
	}

	for _, tt := range checks {
		t.Run(tt.name, func(t *testing.T) {
			got := data[tt.offset : tt.offset+tt.length]
			if !bytes.Equal(got, tt.expected) {
				t.Errorf("At offset %d (%s): expected %x, got %x", tt.offset, tt.name, tt.expected, got)
			}
		})
	}
}

// TestResourceOffsetsMatchData verifies each resource dictionary entry's
// stored Offset points to where that resource's data actually begins.
func TestResourceOffsetsMatchData(t *testing.T) {
	tex := texture.NewTexture(512, 512, texture.PixelFormatRGBA8888, make([]byte, 512*512*4))
	var buf bytes.Buffer
	if err := Write(&buf, tex); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	data := buf.Bytes()
	dynamicHeaderSize := uint32(HeaderSize) + uint32(testNumResources*8)

	loResOffset := binary.LittleEndian.Uint32(data[HeaderSize+4 : HeaderSize+8])
	hiResOffset := binary.LittleEndian.Uint32(data[HeaderSize+12 : HeaderSize+16])

	if loResOffset != dynamicHeaderSize {
		t.Errorf("LowRes offset = %d, want %d", loResOffset, dynamicHeaderSize)
	}

	wantHiResOffset := dynamicHeaderSize + uint32(dxt1Size(16, 16))
	if hiResOffset != wantHiResOffset {
		t.Errorf("HighRes offset = %d, want %d", hiResOffset, wantHiResOffset)
	}
	if int(hiResOffset) >= len(data) {
		t.Fatalf("HighRes offset %d is beyond file length %d", hiResOffset, len(data))
	}
}

func TestVTFWriteRoundTrip(t *testing.T) {
	tex := texture.NewTexture(512, 512, texture.PixelFormatRGBA8888, make([]byte, 512*512*4))
	var buf bytes.Buffer
	if err := Write(&buf, tex); err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	data := buf.Bytes()

	dictSize := uint32(HeaderSize) + uint32(testNumResources*8)
	lowResSize := dxt1Size(16, 16)
	mipSizes := 0
	for size := 256; size >= 1; size /= 2 {
		mipSizes += dxt1Size(size, size)
	}
	highResSize := dxt1Size(512, 512)
	expectedTotal := int(dictSize) + lowResSize + mipSizes + highResSize

	if len(data) != expectedTotal {
		t.Errorf("Expected total size %d, got %d", expectedTotal, len(data))
	}
}

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
			tex := texture.NewTexture(tt.width, tt.height, texture.PixelFormatRGBA8888, make([]byte, tt.width*tt.height*4))
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
