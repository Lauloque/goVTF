/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

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
