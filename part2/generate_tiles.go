package main

import (
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang/snappy"
)

const (
	xSize    = 21600
	ySize    = 10800
	tileSize = 400
	tileName = "world.topo.bathy.200412.3x400x400.%02d.%02d.%s"
)

var chanCodes []string = []string{"red", "green", "blue"}

func GetChannels(img image.Image) []*image.Gray {
	rgba := img.(*image.RGBA)
	rect := rgba.Bounds()
	ch1 := make([]byte, len(rgba.Pix)/4)
	ch2 := make([]byte, len(rgba.Pix)/4)
	ch3 := make([]byte, len(rgba.Pix)/4)
	for i := 0; i < len(ch1); i++ {
		ch1[i] = rgba.Pix[i*4]
		ch2[i] = rgba.Pix[i*4+1]
		ch3[i] = rgba.Pix[i*4+2]
		// Alpha channel is discarded (opaque image)
	}
	red := image.Gray{Pix: ch1, Stride: rect.Dx(), Rect: rect}
	green := image.Gray{Pix: ch2, Stride: rect.Dx(), Rect: rect}
	blue := image.Gray{Pix: ch3, Stride: rect.Dx(), Rect: rect}
	return []*image.Gray{&red, &green, &blue}
}

func GenerateTiles(img image.Image, colour int) {
	for i := 0; i < xSize/tileSize; i++ {
		for j := 0; j < ySize/tileSize; j++ {
			rect := image.Rect(i*tileSize, j*tileSize,
				(i+1)*tileSize, (j+1)*tileSize)
			tile := img.(*image.Gray).SubImage(rect)
			pngOut, _ := os.Create(fmt.Sprintf(tileName+".png", i, j, chanCodes[colour]))
			png.Encode(pngOut, tile)

			b := tile.Bounds()
			width := b.Max.X - b.Min.X
			height := b.Max.Y - b.Min.Y
			pix := make([]uint8, width*height)
			fmt.Println(width, height)

			for y := b.Min.Y; y < b.Max.Y; y++ {
				for x := b.Min.X; x < b.Max.X; x++ {
					pix[(y-b.Min.Y)*width+(x-b.Min.X)] = tile.(*image.Gray).GrayAt(x, y).Y
				}
			}

			SnappyWriter(fmt.Sprintf(tileName+".snpy", i, j, chanCodes[colour]), pix)
			RawWriter(fmt.Sprintf(tileName+".raw", i, j, chanCodes[colour]), pix)
		}
	}
}

func RawWriter(fName string, data []byte) error {
	start := time.Now()

	err := ioutil.WriteFile(fName, data, 0644)
	fmt.Printf("Writting Raw File to disk: %v\n", time.Since(start))

	return err
}

func RawReader(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName)
	fmt.Printf("Reading Raw File from disk: %v\n", time.Since(start))

	return data, err
}

func SnappyWriter(fName string, data []byte) error {
	start := time.Now()

	err := ioutil.WriteFile(fName, snappy.Encode(nil, data), 0644)
	fmt.Printf("Writting Snappy File to disk: %v\n", time.Since(start))

	return err
}

func main() {
	data, err := os.Open("world.topo.bathy.200412.3x21600x10800.png")
	if err != nil {
		panic(err)
	}
	img, _ := png.Decode(data)
	channs := GetChannels(img)
	for i, chann := range channs {
		GenerateTiles(chann, i)
	}
}
