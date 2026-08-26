/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

import (
	"github.com/Lauloque/goVTF/texture"
)

// downSampleHalf outputs a texture half the width and height of the tex.using a
// 2x2 box filter. Assuming width and height are equal, for now.
// TBD: support uneven texture
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

// GenerateMipmaps returns DXT1-compressed mip levels, ordered from smallets to
// largest (VTF spec unlike usual common DDS).
// Ignoring top level already included in highres resource.
func GenerateMipmaps(tex *texture.Texture) [][]byte {
	var levels []*texture.Texture

	currentTex := tex
	for currentTex.Width > 1 && currentTex.Height > 1 {
		currentTex = downSampleHalf(currentTex)
		levels = append(levels, currentTex)
	}

	// levels are processed largets->smallest so to keep downsampling from the
	// current level downwards at each level.
	// Therefore, need to reverse order while compressing to get back to
	// smallest->largest order.

	mipMaps := make([][]byte, len(levels))
	for i, lvl := range levels {
		mipMaps[len(levels)-1-i] = CompressDXT1(lvl)
	}

	return mipMaps
}

func CountMipmaps(width, height int) int {
	count := 1 // level 0 provided by fullres resource already
	for w, h := width, height; w >= 1 && h >= 1; {
		count++
		w >>= 1
		h >>= 1
	}
	return count
}
