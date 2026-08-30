/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

import (
	"fmt"
	"io"

	"github.com/Lauloque/goVTF/imageutils"
	"github.com/Lauloque/goVTF/texture"
)

type WriteOptions struct {
	AlphaFormat int32
}

func Write(w io.Writer, tex *texture.Texture) error {
	return WriteWithOptions(w, tex, WriteOptions{AlphaFormat: ImageFormatDXT1})
}

func WriteWithOptions(w io.Writer, tex *texture.Texture, options WriteOptions) error {
	// ============VALIDATION===========================================
	if err := imageutils.ValidateDimensions(uint(tex.Width), uint(tex.Height)); err != nil {
		return err
	}
	if options.AlphaFormat != ImageFormatDXT1 && options.AlphaFormat != ImageFormatDXT5 {
		return fmt.Errorf("unsupported alpha format: %d", options.AlphaFormat)
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

	// ==== CALCULATE HEADER SIZE DYNAMICALLY ====
	headerSize := uint32(HeaderSize) + uint32(numResources*8) // 80 + 16 = 96

	// ============ PREPARE LOW RES DATA (DXT1) ========================
	lowResData, err := imageutils.CompressThumbnailDXT1(tex.Pixels, tex.Width, tex.Height, lowResWidth, lowResHeight)
	if err != nil {
		return fmt.Errorf("failed to compress thumbnail: %w", err)
	}

	// ============ CALCULATE MIPMAP COUNT =============================
	mipmapCount := imageutils.CountMipmaps(tex.Width, tex.Height)
	fmt.Printf("Mipmap count: %d\n", mipmapCount)

	// ============ PREPARE HIGH RES DATA ===============================
	highResData := imageutils.CompressDXT1(tex)
	flags := uint32(SprayFlags)
	if options.AlphaFormat == ImageFormatDXT5 {
		highResData = imageutils.CompressDXT5(tex)
		flags |= TextureFlagEightBitAlpha
	}
	// should implement error handling later

	// ============CONSTRUCT HEADER=====================================

	// Create header data and then write it
	header := VTFHeader{
		Signature:     [4]byte{'V', 'T', 'F', 0},
		Version:       [2]uint32{SignatureVersionMajor, SignatureVersionMinor},
		HeaderSize:    headerSize,
		Width:         uint16(tex.Width),
		Height:        uint16(tex.Height),
		Flags:         flags,
		Frames:        1,
		FirstFrame:    0,
		Reflectivity:  [3]float32{1.0, 1.0, 1.0},
		BumpmapScale:  1.0,
		HighResFormat: options.AlphaFormat,
		MipmapCount:   uint8(mipmapCount),
		LowResFormat:  int32(ImageFormatDXT1),
		LowResWidth:   uint8(lowResWidth),
		LowResHeight:  uint8(lowResHeight),
		Depth:         1,
		Padding2:      [3]byte{},
		NumResources:  numResources,
		Padding3:      [8]byte{},
	}

	if err := WriteHeader(w, header); err != nil {
		return err
	}

	// ============ RESOURCES ==========================================
	// Calculate offsets
	// Header (80) + LowResEntry (8) + LowResData + HighResEntry (8)
	lowResOffset := headerSize
	highResOffset := lowResOffset + uint32(len(lowResData))
	var mipMaps [][]byte
	if options.AlphaFormat == ImageFormatDXT5 {
		mipMaps, err = imageutils.GenerateMipmapsDXT5(tex)
	} else {
		mipMaps, err = imageutils.GenerateMipmaps(tex)
	}
	if err != nil {
		return err
	}

	// --- Write resource dictionaries ---
	if err := WriteResourceEntry(w, TagLORES, 0, lowResOffset); err != nil {
		return err
	}
	if err := WriteResourceEntry(w, TagHIRES, 0, highResOffset); err != nil {
		return err
	}

	// --- Write resource data ---
	// Order is important: lowResData -> mipmaps(lowest->largets)->highRes
	if _, err := w.Write(lowResData); err != nil {
		return err
	}

	for _, level := range mipMaps {
		if _, err := w.Write(level); err != nil {
			return err
		}
	}

	if _, err := w.Write(highResData); err != nil {
		return err
	}

	return nil
}
