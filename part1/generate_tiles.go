package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

const (
	xSize    = 21600
	ySize    = 10800
	tileSize = 400
	tileName = "world.topo.bathy.200412.3x400x400.%02d.%02d.png"
)

func GenerateTiles(img image.Image) {
	for i := 0; i < xSize/tileSize; i++ {
		for j := 0; j < ySize/tileSize; j++ {
			rect := image.Rect(i*tileSize, j*tileSize,
				(i+1)*tileSize, (j+1)*tileSize)
			tile := img.(*image.RGBA).SubImage(rect)
			out, _ := os.Create(fmt.Sprintf(tileName, i, j))
			png.Encode(out, tile)
		}
	}
}

func main() {
	data, err := os.Open("world.topo.bathy.200412.3x21600x10800.png")
	if err != nil {
		panic(err)
	}
	img, _ := png.Decode(data)
	GenerateTiles(img)
}
