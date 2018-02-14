package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
)

const (
	xSize    = 21600
	pixDeg   = xSize / 360
	tileSize = 400
	tileName = "world.topo.bathy.200412.3x400x400.%02d.%02d.png"
)

func Mosaic(lat, lon float64) image.Image {
	i := int(.5+(lon+180)) * pixDeg
	j := int(.5+(90-lat)) * pixDeg
	tileC0 := (i - 200) / tileSize
	tileC1 := (i + 199) / tileSize
	tileR0 := (j - 200) / tileSize
	tileR1 := (j + 199) / tileSize
	canvas := image.NewRGBA(image.Rect(0, 0, tileSize, tileSize))
	offYCanvas := 0
	for tileR := tileR0; tileR <= tileR1; tileR++ {
		y0 := 0
		y1 := tileSize
		if tileR == tileR0 {
			y0 = (j - 200) % tileSize
		}
		if tileR == tileR1 {
			y1 = (j+199)%tileSize + 1
		}
		offXCanvas := 0
		for tileC := tileC0; tileC <= tileC1; tileC++ {
			data, err := os.Open(fmt.Sprintf(tileName, tileC, tileR))
			if err != nil {
				panic(err)
			}
			tile, err := png.Decode(data)
			if err != nil {
				panic(err)
			}
			x0 := 0
			x1 := tileSize
			if tileC == tileC0 {
				x0 = (i - 200) % tileSize
			}
			if tileC == tileC1 {
				x1 = (i+199)%tileSize + 1
			}
			rect := image.Rect(offXCanvas, offYCanvas, offXCanvas+x1-x0, offYCanvas+y1-y0)
			draw.Draw(canvas, rect, tile, image.Pt(x0, y0), draw.Over)
			offXCanvas += x1 - x0
		}
		offYCanvas += y1 - y0
	}
	return canvas
}

func main() {
	lat := flag.Float64("lat", 0, "Input latitude [-90, 90]")
	lon := flag.Float64("lon", 0, "Input longitude [-180, 180]")
	flag.Parse()

	im := Mosaic(*lat, *lon)
	f, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}

	png.Encode(f, im.(*image.RGBA))
}
