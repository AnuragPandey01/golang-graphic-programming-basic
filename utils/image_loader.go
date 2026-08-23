package utils

import (
	"image"
	"image/draw"
	"os"
)

func LoadImageAsset(filePath string) ([]byte, image.Rectangle, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return nil, image.Rectangle{}, err
	}
	defer file.Close()
	
	img, _ , err := image.Decode(file)
	if err != nil {
		return nil, image.Rectangle{}, err
	}
	bounds := img.Bounds()

	if rgba, ok := img.(*image.RGBA); ok {
		return rgba.Pix, bounds, nil
	}

	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, img, bounds.Min, draw.Src)

	return rgbaImg.Pix, bounds, nil	
}