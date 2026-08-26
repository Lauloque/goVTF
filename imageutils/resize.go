/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

import "github.com/Lauloque/goVTF/texture"

// Resize resizes an RGBA8888 pixel buffer using nearest-neighbor sampling.
//
// pixels must contain srcW * srcH * 4 bytes.
// All dimensions must be positive.
func Resize(srcPixels []byte, srcW, srcH, dstW, dstH int) []byte {
	newPixels := make([]byte, dstW*dstH*4)
	stepX := srcW / dstW
	stepY := srcH / dstH

	for y := 0; y < dstH; y++ {
		for x := 0; x < dstW; x++ {
			srcX := x * stepX
			srcY := y * stepY
			srcIdx := (srcY*srcW + srcX) * 4
			dstIdx := (y*dstW + x) * 4
			copy(newPixels[dstIdx:], srcPixels[srcIdx:srcIdx+4])
		}
	}
	return newPixels
}

// downSampleHalf outputs a texture half the width and height of the tex.using a
// 2x2 box filter. Used for mipmapping.
// Requires both dimensions to be divisible by 2
func downSampleHalf(tex *texture.Texture) *texture.Texture {
	w, h := tex.Width/2, tex.Height/2
	pixels := make([]byte, w*h*4)

	for py := 0; py < h; py++ {
		for px := 0; px < w; px++ {
			sx, sy := px*2, py*2
			var r, g, b, a uint32
			for dy := 0; dy < 2; dy++ {
				for dx := 0; dx < 2; dx++ {
					srcIdx := ((sy+dy)*tex.Width + (sx + dx)) * 4
					r += uint32(tex.Pixels[srcIdx])
					g += uint32(tex.Pixels[srcIdx+1])
					b += uint32(tex.Pixels[srcIdx+2])
					a += uint32(tex.Pixels[srcIdx+3])
				}
			}
			dstIdx := (py*w + px) * 4
			pixels[dstIdx] = byte(r / 4)
			pixels[dstIdx+1] = byte(g / 4)
			pixels[dstIdx+2] = byte(b / 4)
			pixels[dstIdx+3] = byte(a / 4)
		}
	}
	return texture.NewTexture(w, h, texture.PixelFormatRGBA8888, pixels)
}
