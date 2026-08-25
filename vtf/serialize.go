/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

import (
	"encoding/binary"
	"io"
)

// packedHeader mirrors VTFHeader with exact byte layout (no Go padding).
// Fields are ordered to match the VTF 7.3 spec exactly.
type packedHeader struct {
	Signature     [4]byte
	Version       [2]uint32
	HeaderSize    uint32
	Width         uint16
	Height        uint16
	Flags         uint32
	Frames        uint16
	FirstFrame    uint16
	Padding0      [4]byte
	Reflectivity  [3]float32
	Padding1      [4]byte
	BumpmapScale  float32
	HighResFormat int32
	MipmapCount   uint8
	LowResFormat  int32
	LowResWidth   uint8
	LowResHeight  uint8
	Depth         uint16
	Padding2      [3]byte
	NumResources  uint32
	Padding3      [8]byte
}

// WriteHeader serializes a VTFHeader to the writer.
func WriteHeader(w io.Writer, header VTFHeader) error {
	packed := packedHeader{
		Signature:     header.Signature,
		Version:       header.Version,
		HeaderSize:    header.HeaderSize,
		Width:         header.Width,
		Height:        header.Height,
		Flags:         header.Flags,
		Frames:        header.Frames,
		FirstFrame:    header.FirstFrame,
		Padding0:      header.Padding0,
		Reflectivity:  header.Reflectivity,
		Padding1:      header.Padding1,
		BumpmapScale:  header.BumpmapScale,
		HighResFormat: header.HighResFormat,
		MipmapCount:   header.MipmapCount,
		LowResFormat:  header.LowResFormat,
		LowResWidth:   header.LowResWidth,
		LowResHeight:  header.LowResHeight,
		Depth:         header.Depth,
		Padding2:      header.Padding2,
		NumResources:  header.NumResources,
		Padding3:      header.Padding3,
	}
	return binary.Write(w, binary.LittleEndian, &packed)
}

// WriteResourceEntry writes a single resource entry.
func WriteResourceEntry(w io.Writer, tag [3]byte, flags uint8, offset uint32) error {
	if err := binary.Write(w, binary.LittleEndian, tag); err != nil {
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, flags); err != nil {
		return err
	}
	return binary.Write(w, binary.LittleEndian, offset)
}
