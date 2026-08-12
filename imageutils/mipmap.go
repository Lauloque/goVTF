/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

import (
	"image"

	"github.com/Lauloque/goVTF/texture"
)

// GenerateMipmaps creates all mipmap levels for a texture
// Uses simple box filtering (avg of 2×2 pixel blocks)
func GenerateMipmaps(src *texture.Texture) [][]byte {
	var mipMaps [][]byte

	img := src.ToImage()
	imgBounds := img.Bounds()

	for w, h := src.Width, src.Height; w >= 1 && h >= 1; {
		mipmap := image.NewRGBA(image.Rect(0, 0, w, h))

		for py := 0; py < h; py++ {
			for px := 0; px < w; px++ {
				sx := px * 2
				sy := py * 2

				var r, g, b, a uint32
				count := uint32(0)

				for dy := 0; dy < 2; dy++ {
					for dx := 0; dx < 2; dx++ {
						srcX := sx + dx
						srcY := sy + dy

						if srcX < imgBounds.Dx() && srcY < imgBounds.Dy() {
							srcIdx := (srcY*imgBounds.Dx() + srcX) * 4
							r += uint32(img.Pix[srcIdx])
							g += uint32(img.Pix[srcIdx+1])
							b += uint32(img.Pix[srcIdx+2])
							a += uint32(img.Pix[srcIdx+3])
							count++
						}
					}
				}

				dstIdx := (py*w + px) * 4
				if count > 0 {
					mipmap.Pix[dstIdx] = uint8(r / count)
					mipmap.Pix[dstIdx+1] = uint8(g / count)
					mipmap.Pix[dstIdx+2] = uint8(b / count)
					mipmap.Pix[dstIdx+3] = uint8(a / count)
				} else {
					// Out of bounds pixel - set to transparent black
					mipmap.Pix[dstIdx] = 0
					mipmap.Pix[dstIdx+1] = 0
					mipmap.Pix[dstIdx+2] = 0
					mipmap.Pix[dstIdx+3] = 0
				}
			}
		}

		mipMaps = append(mipMaps, mipmap.Pix)

		w >>= 1
		h >>= 1
	}

	return mipMaps
}

func CountMipmaps(width, height int) int {
	count := 0
	for w, h := width, height; w >= 1 && h >= 1; {
		count++
		w >>= 1
		h >>= 1
	}
	return count
}
