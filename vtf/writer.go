/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Lauloque/goVTF/texture"
)

func Write(w io.Writer, tex *texture.Texture) error {
	// Validate dimensions first
	if tex.Width%4 != 0 || tex.Height%4 != 0 {
		return fmt.Errorf("dimensions must be multiples of 4, got %dx%d", tex.Width, tex.Height)
	}

	// Calculate low-res thumbnail size (max 16px on either dimension, DXT1 blocks)
	lowResWidth := uint8(tex.Width)
	lowResHeight := uint8(tex.Height)
	if lowResWidth > 16 {
		lowResWidth = 16
	}
	if lowResHeight > 16 {
		lowResHeight = 16
	}
	// Round up to nearest 4 for DXT1 compression alignment
	lowResWidth = ((lowResWidth + 3) / 4) * 4
	lowResHeight = ((lowResHeight + 3) / 4) * 4

	header := VTFHeader{
		Signature:     [4]byte{'V', 'T', 'F', 0},
		Version:       [2]uint32{SignatureVersionMajor, SignatureVersionMinor},
		HeaderSize:    HeaderSize,
		Width:         uint16(tex.Width),
		Height:        uint16(tex.Height),
		Flags:         0,
		Frames:        1,
		FirstFrame:    0,
		Reflectivity:  [3]float32{1.0, 1.0, 1.0},
		BumpmapScale:  1.0,
		HighResFormat: ImageFormatRGBA8888,
		MipmapCount:   1,
		LowResFormat:  ImageFormatDXT1,
		LowResWidth:   lowResWidth,
		LowResHeight:  lowResHeight,
		Depth:         1,
		NumResources:  1, // Only high-res image, no thumbnail for minimal v1
	}

	// Write header
	if err := binary.Write(w, binary.LittleEndian, header); err != nil {
		return err
	}

	// Write resource entry for high-res image
	// Offset = header size (96) + resource entry size (8)
	resource := ResourceEntry{
		Tag:    TagHIRES,
		Flags:  0,
		Offset: HeaderSize + 8,
	}
	if err := binary.Write(w, binary.LittleEndian, resource); err != nil {
		return err
	}

	// Write raw RGBA8888 pixel data
	if _, err := w.Write(tex.Pixels); err != nil {
		return err
	}

	return nil
}
