/* SPDX-License-Identifier: GPL-3.0-or-later */
package texture

import (
	"image"
	"image/color"
)

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

func (t *Texture) ToImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, t.Width, t.Height))
	copy(img.Pix, t.Pixels)
	return img
}

func FromImage(img *image.RGBA, format PixelFormat) *Texture {
	return NewTexture(
		img.Rect.Dx(),
		img.Rect.Dy(),
		format,
		img.Pix,
	)
}

// ColorAt returns the color at pixel (x, y)
func (t *Texture) ColorAt(x, y int) color.Color {
	if x < 0 || y < 0 || x >= t.Width || y >= t.Height {
		return color.RGBA{0, 0, 0, 0}
	}
	idx := (y*t.Width + x) * 4
	r, g, b, a := t.Pixels[idx], t.Pixels[idx+1], t.Pixels[idx+2], t.Pixels[idx+3]
	return color.RGBA{r, g, b, a}
}
