/* SPDX-License-Identifier: GPL-3.0-or-later */
package vtf

// CompressThumbnailDXT1 performs a basic DXT1 compression inline.
// This implementation encodes 4x4 blocks into 8 bytes using a simple
// average-based endpoint selection. It is sufficient for TF2 sprays.
func CompressThumbnailDXT1(srcPixels []byte, srcW, srcH, dstW, dstH int) ([]byte, error) {
	// 1. Subsample (Nearest Neighbor) to thumbnail size
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

	// 2. DXT1 Compression
	// DXT1 encodes 4x4 blocks into 8 bytes.
	// Dimensions must be multiples of 4.
	blockW := dstW / 4
	blockH := dstH / 4
	totalBlocks := blockW * blockH
	output := make([]byte, totalBlocks*8)

	blockIdx := 0
	for by := 0; by < blockH; by++ {
		for bx := 0; bx < blockW; bx++ {
			// Extract 4x4 block pixels
			var r, g, b, a [16]byte
			for py := 0; py < 4; py++ {
				for px := 0; px < 4; px++ {
					x := bx*4 + px
					y := by*4 + py
					idx := (y*dstW + x) * 4
					r[py*4+px] = thumbPixels[idx]
					g[py*4+px] = thumbPixels[idx+1]
					b[py*4+px] = thumbPixels[idx+2]
					a[py*4+px] = thumbPixels[idx+3]
				}
			}

			// Calculate average color for endpoints
			var sumR, sumG, sumB uint32
			for i := 0; i < 16; i++ {
				sumR += uint32(r[i])
				sumG += uint32(g[i])
				sumB += uint32(b[i])
			}
			avgR := byte(sumR / 16)
			avgG := byte(sumG / 16)
			avgB := byte(sumB / 16)

			// Find min and max colors to use as endpoints (simplified)
			// In a real encoder, we'd find the two most distinct colors.
			// Here we use Max and Min of the block for simplicity.
			minR, maxR := r[0], r[0]
			minG, maxG := g[0], g[0]
			minB, maxB := b[0], b[0]
			for i := 1; i < 16; i++ {
				if r[i] < minR {
					minR = r[i]
				}
				if r[i] > maxR {
					maxR = r[i]
				}
				if g[i] < minG {
					minG = g[i]
				}
				if g[i] > maxG {
					maxG = g[i]
				}
				if b[i] < minB {
					minB = b[i]
				}
				if b[i] > maxB {
					maxB = b[i]
				}
			}

			// Pack endpoints (5-6-5 RGB)
			// c0 = max, c1 = min (DXT1 convention)
			c0 := (uint16(maxR) >> 3) | (uint16(maxG)>>2)<<5 | (uint16(maxB)>>3)<<11
			c1 := (uint16(minR) >> 3) | (uint16(minG)>>2)<<5 | (uint16(minB)>>3)<<11

			// Write c0, c1 (little-endian)
			output[blockIdx*8] = byte(c0)
			output[blockIdx*8+1] = byte(c0 >> 8)
			output[blockIdx*8+2] = byte(c1)
			output[blockIdx*8+3] = byte(c1 >> 8)

			// Generate index map (2 bits per pixel)
			// 00 = c0, 01 = c1, 10 = interpolated (c0+c1)/2, 11 = interpolated (c0+2*c1)/3?
			// Standard DXT1 interpolation:
			// 00 -> c0
			// 01 -> c1
			// 10 -> (2*c0 + c1) / 3
			// 11 -> (c0 + 2*c1) / 3
			// If c0 <= c1, then 11 is transparent (handled by alpha, but we ignore alpha here for simplicity)

			var indices uint32
			for i := 0; i < 16; i++ {
				// Simple distance to endpoints
				dR := int(r[i]) - int(avgR)
				dG := int(g[i]) - int(avgG)
				dB := int(b[i]) - int(avgB)
				dist := dR*dR + dG*dG + dB*dB

				// Map to 0-3 based on distance to avg (rough heuristic)
				// This is a simplification. A real encoder checks distance to c0 and c1.
				idx := uint8(0)
				if dist > 2000 {
					idx = 1 // Far from avg -> closer to max?
				} else if dist > 500 {
					idx = 2 // Mid
				} else {
					idx = 0 // Close to avg -> closer to min?
				}

				// Better heuristic: Distance to c0 vs c1
				// c0 is max, c1 is min.
				// Let's just use a simple linear interpolation check.
				// This is a very rough approximation but valid for loading.

				// Re-calculate distance to c0 and c1
				// c0 color
				c0R := byte(((maxR >> 3) * 8) | ((maxR >> 3) >> 5)) // Rough expansion
				c0G := byte(((maxG >> 2) * 4) | ((maxG >> 2) >> 6))
				c0B := byte(((maxB >> 3) * 8) | ((maxB >> 3) >> 5))

				// c1 color
				c1R := byte(((minR >> 3) * 8) | ((minR >> 3) >> 5))
				c1G := byte(((minG >> 2) * 4) | ((minG >> 2) >> 6))
				c1B := byte(((minB >> 3) * 8) | ((minB >> 3) >> 5))

				d0 := (int(r[i])-int(c0R))*(int(r[i])-int(c0R)) + (int(g[i])-int(c0G))*(int(g[i])-int(c0G)) + (int(b[i])-int(c0B))*(int(b[i])-int(c0B))
				d1 := (int(r[i])-int(c1R))*(int(r[i])-int(c1R)) + (int(g[i])-int(c1G))*(int(g[i])-int(c1G)) + (int(b[i])-int(c1B))*(int(b[i])-int(c1B))

				if d0 < d1 {
					idx = 0
				} else {
					// Check midpoints
					midR := byte((int(c0R) + int(c1R)) / 2)
					midG := byte((int(c0G) + int(c1G)) / 2)
					midB := byte((int(c0B) + int(c1B)) / 2)
					dMid := (int(r[i])-int(midR))*(int(r[i])-int(midR)) + (int(g[i])-int(midG))*(int(g[i])-int(midG)) + (int(b[i])-int(midB))*(int(b[i])-int(midB))

					if dMid < d1 {
						idx = 2
					} else {
						idx = 1
					}
				}

				indices |= uint32(idx) << (i * 2)
			}

			// Write indices (4 bytes)
			output[blockIdx*8+4] = byte(indices)
			output[blockIdx*8+5] = byte(indices >> 8)
			output[blockIdx*8+6] = byte(indices >> 16)
			output[blockIdx*8+7] = byte(indices >> 24)

			blockIdx++
		}
	}

	return output, nil
}
