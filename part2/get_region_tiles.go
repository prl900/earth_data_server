package main

import (
	"flag"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang/snappy"
)

const (
	xSize    = 21600
	pixDeg   = xSize / 360
	tileSize = 400
	tileName = "world.topo.bathy.200412.3x400x400.%02d.%02d.%s"
)

var colChans []string = []string{"red", "green", "blue"}

func MosaicPNG(lat, lon float64, colChan int) image.Image {
	i := int(.5+(lon+180)) * pixDeg
	j := int(.5+(90-lat)) * pixDeg
	tileC0 := (i - 200) / tileSize
	tileC1 := (i + 199) / tileSize
	tileR0 := (j - 200) / tileSize
	tileR1 := (j + 199) / tileSize
	canvas := image.NewGray(image.Rect(0, 0, tileSize, tileSize))
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
			data, err := os.Open(fmt.Sprintf(tileName+".png", tileC, tileR, colChans[colChan]))
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
			draw.Draw(canvas, rect, tile.(*image.Gray), image.Pt(x0, y0), draw.Over)
			offXCanvas += x1 - x0
		}
		offYCanvas += y1 - y0
	}
	return canvas
}

func MosaicRaw(lat, lon float64, colChan int) image.Image {
	i := int(.5+(lon+180)) * pixDeg
	j := int(.5+(90-lat)) * pixDeg
	tileC0 := (i - 200) / tileSize
	tileC1 := (i + 199) / tileSize
	tileR0 := (j - 200) / tileSize
	tileR1 := (j + 199) / tileSize
	canvas := image.NewGray(image.Rect(0, 0, tileSize, tileSize))
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
			data, err := ioutil.ReadFile(fmt.Sprintf(tileName+".raw", tileC, tileR, colChans[colChan]))
			if err != nil {
				panic(err)
			}
			tile := &image.Gray{Pix: data, Stride: 400, Rect: image.Rect(0, 0, 400, 400)}
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

func SnappyReader(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName + ".snp")
	cdata, err := snappy.Decode(nil, data)
	fmt.Printf("Reading Snappy File from disk: %v\n", time.Since(start))

	return cdata, err
}
func MosaicSnappy(lat, lon float64, colChan int) image.Image {
	i := int(.5+(lon+180)) * pixDeg
	j := int(.5+(90-lat)) * pixDeg
	tileC0 := (i - 200) / tileSize
	tileC1 := (i + 199) / tileSize
	tileR0 := (j - 200) / tileSize
	tileR1 := (j + 199) / tileSize
	canvas := image.NewGray(image.Rect(0, 0, tileSize, tileSize))
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
			data, err := ioutil.ReadFile(fmt.Sprintf(tileName+".snpy", tileC, tileR, colChans[colChan]))
			cdata, err := snappy.Decode(nil, data)
			if err != nil {
				panic(err)
			}
			tile := &image.Gray{Pix: cdata, Stride: 400, Rect: image.Rect(0, 0, 400, 400)}
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
	chann := flag.Int("chan", 0, "Colour channel R=0, G=1, B=2")
	flag.Parse()

	start := time.Now()
	im := MosaicPNG(*lat, *lon, *chann)
	f, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generating PNG tile: %v\n", time.Since(start))

	png.Encode(f, im.(*image.Gray))

	start = time.Now()
	im = MosaicRaw(*lat, *lon, *chann)
	f, err = os.Create("out2.png")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generating Raw tile: %v\n", time.Since(start))

	png.Encode(f, im.(*image.Gray))

	start = time.Now()
	im = MosaicSnappy(*lat, *lon, *chann)
	f, err = os.Create("out3.png")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Generating Snappy tile: %v\n", time.Since(start))

	png.Encode(f, im.(*image.Gray))
}
