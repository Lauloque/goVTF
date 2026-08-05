/* SPDX-License-Identifier: GPL-3.0-or-later */
package texture

type PixelFormat int

const (
	PixelFormatRGBA8888 PixelFormat = iota // 32b/px, 8b red, 8b green, 8b blue, 8b alpha
	PixelFormatDXT1                 = 14   // no alpha ot 1 bit, 4b/px, most diffuse textures, equivalent to modern BC1
	PixelFormatDXT5                 = 15   // full 8b alpha, 8b/px, textures with smooth transparency (leaves, decals, UI...), equivalent to modern BC3
)

type Texture struct {
	Width  int
	Height int
	Format PixelFormat
	Pixels []byte
}

func NewTexture(width, height int, format PixelFormat, pixels []byte) *Texture {
	return &Texture{
		Width:  width,
		Height: height,
		Format: format,
		Pixels: pixels,
	}
}
