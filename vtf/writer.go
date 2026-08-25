/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

import (
	"fmt"
	"io"

	"github.com/Lauloque/goVTF/imageutils"
	"github.com/Lauloque/goVTF/texture"
)

func isPowerOfTwo(n uint) bool {
	return n > 0 && (n&(n-1)) == 0
}

func Write(w io.Writer, tex *texture.Texture) error {
	// ============VALIDATION===========================================
	if tex.Width <= 0 || tex.Height <= 0 {
		return fmt.Errorf("dimensions must be positive")
	}

	if tex.Width%4 != 0 || tex.Height%4 != 0 {
		return fmt.Errorf("dimensions must be multiples of 4, got %dx%d", tex.Width, tex.Height)
	}

	if !isPowerOfTwo(uint(tex.Width)) || !isPowerOfTwo(uint(tex.Height)) {
		return fmt.Errorf("dimensions must be power of 2, got %dx%d", tex.Width, tex.Height)
	}

	// ============CALCULATE LOW REST DIMENSION ========================
	// (max 16x16)
	lowResWidth := tex.Width
	lowResHeight := tex.Height
	if lowResWidth > 16 {
		lowResWidth = 16
	}
	if lowResHeight > 16 {
		lowResHeight = 16
	}
	if lowResWidth < 4 {
		lowResWidth = 4
	}
	if lowResHeight < 4 {
		lowResHeight = 4
	}
	// Round up to nearest 4 for DXT1 compression alignment
	// Hence multiple-of-4 requirement!
	lowResWidth = ((lowResWidth + 3) / 4) * 4
	lowResHeight = ((lowResHeight + 3) / 4) * 4

	// ============ PREPARE LOW RES DATA (DXT1) ========================
	lowResData, err := CompressThumbnailDXT1(tex.Pixels, tex.Width, tex.Height, lowResWidth, lowResHeight)
	if err != nil {
		return fmt.Errorf("failed to compress thumbnail: %w", err)
	}

	// Calculate mipmap count
	//
	// temporary forced to no mipmap / full res
	mipmapCount := 1
	// mipmapCount := imageutils.CountMipmaps(tex.Width, tex.Height)
	// fmt.Printf("Mipmap count: %d\n", mipmapCount)

	// ============ PREPARE HIGH RES DATA (DXT1)========================
	highResData := imageutils.CompressDXT1(tex)
	// should implement error handling later

	// ============CONSTRUCT HEADER=====================================

	// Create header data and then write it
	header := VTFHeader{
		Signature:     [4]byte{'V', 'T', 'F', 0},
		Version:       [2]uint32{SignatureVersionMajor, SignatureVersionMinor},
		HeaderSize:    HeaderSize,
		Width:         uint16(tex.Width),
		Height:        uint16(tex.Height),
		Flags:         uint32(SprayFlags),
		Frames:        1,
		FirstFrame:    0,
		Reflectivity:  [3]float32{1.0, 1.0, 1.0},
		BumpmapScale:  1.0,
		HighResFormat: int32(ImageFormatDXT1), // Forced DXT1 for v1
		MipmapCount:   1,                      // None for v1
		LowResFormat:  int32(ImageFormatDXT1),
		LowResWidth:   uint8(lowResWidth),
		LowResHeight:  uint8(lowResHeight),
		Depth:         1,
		Padding2:      [3]byte{},
		NumResources:  2,
		Padding3:      [8]byte{},
	}

	if err := WriteHeader(w, header); err != nil {
		return err
	}

	// ============ RESOURCES ==========================================
	// Calculate offsets
	// Header (80) + LowResEntry (8) + LowResData + HighResEntry (8)
	lowResOffset := uint32(HeaderSize) + uint32(8*header.NumResources) // 80 + 16 = 96
	highResOffset := lowResOffset + 8 + uint32(len(lowResData))

	// --- Write resource dictionaries ---
	if err := WriteResourceEntry(w, TagLORES, 0, lowResOffset); err != nil {
		return err
	}
	if err := WriteResourceEntry(w, TagHIRES, 0, highResOffset); err != nil {
		return err
	}

	// --- Write resource data ---
	if _, err := w.Write(lowResData); err != nil {
		return err
	}
	if _, err := w.Write(highResData); err != nil {
		return err
	}

	// -------------------------------------------------------------
	// Write DXT1-compressed pixel data
	// if _, err := w.Write(compressedHiRes); err != nil {
	// 	return err
	// }

	// Write additional mipmaps (smaller than full res)
	// Skip the first one since we already wrote the full res above
	//
	// temporary disabled for testing whether providing a lowres thumbnail is enough
	// toggle back header[56] if needed.
	//
	if mipmapCount > 1 {
		mipMaps := imageutils.GenerateMipmaps(tex)
		for i := 1; i < len(mipMaps); i++ {
			if _, err := w.Write(mipMaps[i]); err != nil {
				return err
			}
		}
	}

	return nil
}
