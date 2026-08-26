/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

import (
	"github.com/Lauloque/goVTF/texture"
)

// CompressDXT1 compresses RGBA texture to DXT1 (BC1) format
// Returns compressed bytes (8 bytes per 4×4 block)
func CompressDXT1(tex *texture.Texture) []byte {
	width := tex.Width
	height := tex.Height

	// DXT1: 8 bytes per 4×4 block
	tilesW := (width + 3) / 4
	tilesH := (height + 3) / 4
	compressed := make([]byte, tilesW*tilesH*8)

	// Process each 4×4 block
	for ty := 0; ty < tilesH; ty++ {
		for tx := 0; tx < tilesW; tx++ {
			blockIndex := (ty*tilesW + tx) * 8

			// Extract 4×4 pixels from source
			var pixels [16]uint32 // RGB888 packed in uint32
			for py := 0; py < 4; py++ {
				for px := 0; px < 4; px++ {
					x := tx*4 + px
					y := ty*4 + py

					if x < width && y < height {
						idx := (y*width + x) * 4
						r := uint32(tex.Pixels[idx])
						g := uint32(tex.Pixels[idx+1])
						b := uint32(tex.Pixels[idx+2])

						// Pack RGB into uint32
						pixels[py*4+px] = (r << 16) | (g << 8) | b
					} else {
						pixels[py*4+px] = 0
					}
				}
			}

			// Compress this 4×4 block
			compressedBlock := compressBlockDXT1(pixels)
			copy(compressed[blockIndex:blockIndex+8], compressedBlock[:])
		}
	}

	return compressed
}

// compressBlockDXT1 compresses a single 4×4 pixel block to DXT1 (8 bytes)
func compressBlockDXT1(pixels [16]uint32) [8]byte {
	// Step 1: Find min and max colors (by luminance or simple sum)
	var minColor, maxColor uint32
	minColor = 0xFFFFFFFF
	maxColor = 0x00000000

	for i := 0; i < 16; i++ {
		c := pixels[i]
		r, g, b := (c>>16)&0xFF, (c>>8)&0xFF, c&0xFF

		// Simple brightness check (sum of RGB)
		brightness := uint32(r + g + b)
		minBright := uint32((minColor>>16)&0xFF + (minColor>>8)&0xFF + minColor&0xFF)
		maxBright := uint32((maxColor>>16)&0xFF + (maxColor>>8)&0xFF + maxColor&0xFF)

		if brightness < minBright {
			minColor = c
		}
		if brightness > maxBright {
			maxColor = c
		}
	}

	// Step 2: Convert to 5:6:5 format
	c0 := rgbTo565(maxColor)
	c1 := rgbTo565(minColor)

	// Ensure c0 > c1 (required for DXT1 gradient direction)
	if c0 <= c1 {
		c0, c1 = c1, c0
	}

	// Step 3: Build 4-color palette
	var palette [4]uint32
	palette[0] = expand565ToRGB(c0)
	palette[1] = expand565ToRGB(c1)

	// Interpolated colors (DXT1 spec)
	// Color 2 = (2*c0 + c1) / 3
	// Color 3 = (c0 + 2*c1) / 3
	r0, g0, b0 := (palette[0]>>16)&0xFF, (palette[0]>>8)&0xFF, palette[0]&0xFF
	r1, g1, b1 := (palette[1]>>16)&0xFF, (palette[1]>>8)&0xFF, palette[1]&0xFF

	palette[2] = ((r0*2+r1)/3)<<16 | ((g0*2+g1)/3)<<8 | ((b0*2 + b1) / 3)
	palette[3] = ((r0+r1*2)/3)<<16 | ((g0+g1*2)/3)<<8 | ((b0 + b1*2) / 3)

	// Step 4: Assign 2-bit indices to each pixel
	var indices [16]byte
	for i := 0; i < 16; i++ {
		p := pixels[i]
		pr, pg, pb := (p>>16)&0xFF, (p>>8)&0xFF, p&0xFF

		bestDist := uint32(0xFFFFFFFF)
		bestIdx := byte(0)

		for j := 0; j < 4; j++ {
			pr2, pg2, pb2 := (palette[j]>>16)&0xFF, (palette[j]>>8)&0xFF, palette[j]&0xFF

			// Manhattan distance
			dr := uint32(pr) - pr2
			dg := uint32(pg) - pg2
			db := uint32(pb) - pb2

			dist := dr + dg + db
			if dist < bestDist {
				bestDist = dist
				bestIdx = byte(j)
			}
		}
		indices[i] = bestIdx
	}

	// Step 5: Pack into 8 bytes
	var out [8]byte

	// Bytes 0-1: c0 (big-endian 16-bit)
	out[0] = byte(c0)
	out[1] = byte(c0 >> 8)

	// Bytes 2-3: c1 (big-endian 16-bit)
	out[2] = byte(c1)
	out[3] = byte(c1 >> 8)

	// Bytes 4-7: 16 × 2-bit indices (packed)
	// Group 4 indices per byte
	for i := 0; i < 4; i++ {
		base := i * 4
		val := (indices[base+0] << 0) |
			(indices[base+1] << 2) |
			(indices[base+2] << 4) |
			(indices[base+3] << 6)
		out[4+i] = byte(val)
	}

	return out
}

func expand565ToRGB(c uint16) uint32 {
	r := (c >> 11) & 0x1F
	g := (c >> 5) & 0x3F
	b := c & 0x1F

	// Expand 5/6 bits to 8 bits
	r8 := (r << 3) | (r >> 2)
	g8 := (g << 2) | (g >> 4)
	b8 := (b << 3) | (b >> 2)

	return (uint32(r8) << 16) | (uint32(g8) << 8) | uint32(b8)
}

// rgbTo565 converts RGB888 to RGB565
func rgbTo565(rgb uint32) uint16 {
	r := (rgb >> 16) & 0xFF
	g := (rgb >> 8) & 0xFF
	b := rgb & 0xFF

	r5 := uint16((r * 31) / 255)
	g6 := uint16((g * 63) / 255)
	b5 := uint16((b * 31) / 255)

	return (r5 << 11) | (g6 << 5) | b5
}
