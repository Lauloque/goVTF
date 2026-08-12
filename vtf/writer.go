/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/Lauloque/goVTF/texture"
)

func isPowerOfTwo(n uint) bool {
	return n > 0 && (n&(n-1)) == 0
}

func Write(w io.Writer, tex *texture.Texture) error {
	// Validate dimensions first
	if tex.Width%4 != 0 || tex.Height%4 != 0 {
		return fmt.Errorf("dimensions must be multiples of 4, got %dx%d", tex.Width, tex.Height)
	}

	if !isPowerOfTwo(uint(tex.Width)) || !isPowerOfTwo(uint(tex.Height)) {
		return fmt.Errorf("dimensions must be power of 2, got %dx%d", tex.Width, tex.Height)
	}

	// Calculate low-res thumbnail size (max 16px on either dimension, DXT1 blocks)
	lowResWidth := 16
	lowResHeight := 16
	if tex.Width <= 16 {
		lowResWidth = tex.Width
	}
	if tex.Height <= 16 {
		lowResHeight = tex.Height
	}
	// Round up to nearest 4 for DXT1 compression alignment
	// Hence multiple-of-4 requirement!
	lowResWidth = ((lowResWidth + 3) / 4) * 4
	lowResHeight = ((lowResHeight + 3) / 4) * 4

	// Build header manually (96 bytes exactly per VTF 7.4 spec)
	header := make([]byte, 96)

	// Signature: 4 bytes
	copy(header[0:4], []byte("VTF\x00"))

	// Version: 2 × uint32 (offsets 4-11)
	binary.LittleEndian.PutUint32(header[4:8], SignatureVersionMajor)
	binary.LittleEndian.PutUint32(header[8:12], SignatureVersionMinor)

	// HeaderSize: uint32 (offset 12-15)
	binary.LittleEndian.PutUint32(header[12:16], HeaderSize)

	// Width/Height: uint16 (offsets 16-19)
	binary.LittleEndian.PutUint16(header[16:18], uint16(tex.Width))
	binary.LittleEndian.PutUint16(header[18:20], uint16(tex.Height))

	// Flags: uint32 (offset 20-23)
	binary.LittleEndian.PutUint32(header[20:24], uint32(SprayFlags))

	// Frames: uint16 (offset 24-25)
	binary.LittleEndian.PutUint16(header[24:26], 1)

	// FirstFrame: uint16 (offset 26-27)
	binary.LittleEndian.PutUint16(header[26:28], 0)

	// Padding0: 4 bytes (offset 28-31)
	copy(header[28:32], []byte{0, 0, 0, 0})

	// Reflectivity: 3 × float32 (offset 32-43)
	binary.LittleEndian.PutUint32(header[32:36], 0x3f800000) // 1.0
	binary.LittleEndian.PutUint32(header[36:40], 0x3f800000) // 1.0
	binary.LittleEndian.PutUint32(header[40:44], 0x3f800000) // 1.0

	// Padding1: 4 bytes (offset 44-47)
	copy(header[44:48], []byte{0, 0, 0, 0})

	// BumpmapScale: float32 (offset 48-51)
	binary.LittleEndian.PutUint32(header[48:52], 0x3f800000) // 1.0

	// HighResFormat: int32 (offset 52-55)
	binary.LittleEndian.PutUint32(header[52:56], ImageFormatRGBA8888)

	// MipmapCount: uint8 (offset 56)
	header[56] = 1

	// LowResFormat: int32 (offset 57-60)
	binary.LittleEndian.PutUint32(header[57:61], ImageFormatDXT1)

	// LowResWidth: uint8 (offset 61)
	header[61] = uint8(lowResWidth)

	// LowResHeight: uint8 (offset 62)
	header[62] = uint8(lowResHeight)

	// Depth: uint16 (offset 63-64)
	binary.LittleEndian.PutUint16(header[63:65], 1)

	// Padding2: 3 bytes (offset 65-67)
	copy(header[65:68], []byte{0, 0, 0})

	// NumResources: uint32 (offset 68-71)
	binary.LittleEndian.PutUint32(header[68:72], 1)

	// Padding3: 8 bytes (offset 72-79)
	copy(header[72:80], []byte{0, 0, 0, 0, 0, 0, 0, 0})

	// Write header
	if _, err := w.Write(header); err != nil {
		return err
	}

	// Write resource entry for high-res image (8 bytes)
	// DISABLED: to avoid TF2 saying `*** Error unserializing VTF file... is the file empty?`
	//
	// resource := make([]byte, 8)
	// copy(resource[0:3], TagHIRES[:])                  // Tag: 3 bytes
	// resource[3] = 0                                   // Flags: 1 byte
	// binary.LittleEndian.PutUint32(resource[4:8], 104) // Offset: 96 + 8 = 104

	// if _, err := w.Write(resource); err != nil {
	// 	return err
	// }

	// Write raw RGBA8888 pixel data
	if _, err := w.Write(tex.Pixels); err != nil {
		return err
	}

	return nil
}
