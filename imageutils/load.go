/* SPDX-License-Identifier: GPL-3.0-or-later */
package imageutils

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/Lauloque/goVTF/texture"
)

func ReadImage(inputPath string) {
	f, err := os.Open(inputPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("Decoded %s\n", format)
	fmt.Printf("Bounds: %v\n", img.Bounds())
}

func Load(inputPath string) (image.Image, string, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}

func LoadTexture(inputPath string) *texture.Texture {
	img, _, err := Load(inputPath)
	if err != nil {
		return nil
	}

	bounds := img.Bounds()

	// use `bounds.Min` instead of `image.Point{}` since the latter assumes source image starts at (0, 0) while the former allows for non-zero bounds (i.e. cropped sub-images)
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, img, bounds.Min, draw.Src)

	return texture.NewTexture(
		bounds.Dx(),
		bounds.Dy(),
		texture.PixelFormatRGBA8888,
		rgba.Pix,
	)
}
