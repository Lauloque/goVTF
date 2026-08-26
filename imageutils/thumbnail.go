/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

import (
	"github.com/Lauloque/goVTF/texture"
)

// CompressThumbnailDXT1 performs a basic DXT1 compression inline.
// This implementation encodes 4x4 blocks into 8 bytes using a simple
// average-based endpoint selection. It is sufficient for TF2 sprays.
func CompressThumbnailDXT1(srcPixels []byte, srcW, srcH, dstW, dstH int) ([]byte, error) {
	// Subsample (Nearest Neighbor) to thumbnail size
	thumbPixels := make([]byte, dstW*dstH*4)
	stepX := srcW / dstW
	stepY := srcH / dstH

	for y := 0; y < dstH; y++ {
		for x := 0; x < dstW; x++ {
			srcX := x * stepX
			srcY := y * stepY
			srcIdx := (srcY*srcW + srcX) * 4
			dstIdx := (y*dstW + x) * 4
			copy(thumbPixels[dstIdx:], srcPixels[srcIdx:srcIdx+4])
		}
	}

	// Compress
	thumbTex := texture.NewTexture(dstW, dstH, texture.PixelFormatRGBA8888, thumbPixels)
	return CompressDXT1(thumbTex), nil
}
