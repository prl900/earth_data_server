package main

import (
	"flag"
	"image"
	"image/png"
	"os"
)

const (
	xSize    = 21600
	pixDeg   = xSize / 360
	fileName = "world.topo.bathy.200412.3x21600x10800.png"
)

func Region(lat, lon float64) image.Image {
	// i & j contain the pixel position of the input coordinates
	i := int(.5+(lon+180)) * pixDeg
	j := int(.5+(90-lat)) * pixDeg
	rect := image.Rect(i-200, j-200, i+200, j+200)
	data, _ := os.Open(fileName)
	img, _ := png.Decode(data)
	tile := img.(*image.RGBA).SubImage(rect)

	return tile
}

func main() {
	lat := flag.Float64("lat", 0, "Input latitude [-90, 90]")
	lon := flag.Float64("lon", 0, "Input longitude [-180, 180]")
	flag.Parse()

	im := Region(*lat, *lon)
	f, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}

	png.Encode(f, im.(*image.RGBA))
}
