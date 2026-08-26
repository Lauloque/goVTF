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
	thumbPixels, err := Resize(srcPixels, srcW, srcH, dstW, dstH)
	if err != nil {
		return nil, err
	}

	// Compress
	thumbTex := texture.NewTexture(dstW, dstH, texture.PixelFormatRGBA8888, thumbPixels)
	return CompressDXT1(thumbTex), nil
}
