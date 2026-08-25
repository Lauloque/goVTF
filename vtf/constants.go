/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

// Based on https://developer.valvesoftware.com/wiki/VTF_(Valve_Texture_Format)

// File signatures and format codes
const (
	SignatureVersionMajor = 7
	SignatureVersionMinor = 3
	HeaderSize            = 80
	ImageFormatRGBA8888   = 0
	ImageFormatDXT1       = 14
	ImageFormatDXT5       = 16
	TextureFlagClampS     = 0x00000004
	TextureFlagClampT     = 0x00000008
	TextureFlagNormal     = 0x00000080
)

// pseudo-constants that cannot be constants in Go
var (
	TagHIRES   = [3]byte{0x30, 0x00, 0x00} // High-res image data
	TagLORES   = [3]byte{0x01, 0x00, 0x00} // Low-res thumbnail
	SprayFlags = TextureFlagClampS | TextureFlagClampT
)

// Helper to create resource tag from string (for extension purposes)
func MakeTag(tag string) [3]byte {
	b := [3]byte{}
	copy(b[:], tag)
	return b
}
