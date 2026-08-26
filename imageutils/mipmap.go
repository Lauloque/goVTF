/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

import (
	"github.com/Lauloque/goVTF/texture"
)

// GenerateMipmaps returns DXT1-compressed mip levels, ordered from smallets to
// largest (VTF spec unlike usual common DDS).
// Ignoring top level already included in highres resource.
func GenerateMipmaps(tex *texture.Texture) ([][]byte, error) {
	var levels []*texture.Texture

	currentTex := tex
	for currentTex.Width > 1 && currentTex.Height > 1 {
		var err error
		currentTex, err = downSampleHalf(currentTex)
		if err != nil {
			return nil, err
		}
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

	return mipMaps, nil
}

func CountMipmaps(width, height int) int {
	count := 1 // level 0 provided by fullres resource already
	for w, h := width, height; w > 1 && h > 1; {
		count++
		w >>= 1
		h >>= 1
	}
	return count
}
