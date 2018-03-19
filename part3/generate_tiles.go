package main

import (
	"fmt"
	"image"
	"image/png"
	"os"

	"cloud.google.com/go/storage"
	"github.com/golang/snappy"
	"golang.org/x/net/context"
)

const (
	xSize    = 21600
	ySize    = 10800
	tileSize = 400
	tileName = "world.topo.bathy.200412.3x400x400.%02d.%02d.%s"
	bktName  = "bluemarble"
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

func WriteObject(bktName, objName string, contents []byte) error {
	ctx := context.Background()

	// Sets your Google Cloud Platform project ID.
	//projectID := "YOUR_PROJECT_ID"
	//projectID := "nci-gce"

	// Creates a client.
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("Failed to create client: %v", err)
	}

	// Creates a Bucket instance.
	bucket := client.Bucket(bktName)
	w := bucket.Object(objName).NewWriter(ctx)
	w.ContentType = "application/octet-stream"

	if _, err := w.Write([]byte(contents)); err != nil {
		return fmt.Errorf("Failed to write object to bucket: %v", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("Failed to close writer to bucket: %v", err)
	}
	// Close the client when finished.
	if err := client.Close(); err != nil {
		return fmt.Errorf("Failed to close client: %v", err)
	}

	return nil
}

func GenerateTiles(img image.Image, colour int) {
	for i := 0; i < xSize/tileSize; i++ {
		for j := 0; j < ySize/tileSize; j++ {
			rect := image.Rect(i*tileSize, j*tileSize,
				(i+1)*tileSize, (j+1)*tileSize)
			tile := img.(*image.Gray).SubImage(rect)

			b := tile.Bounds()
			width := b.Max.X - b.Min.X
			height := b.Max.Y - b.Min.Y
			pix := make([]uint8, width*height)

			for y := b.Min.Y; y < b.Max.Y; y++ {
				for x := b.Min.X; x < b.Max.X; x++ {
					pix[(y-b.Min.Y)*width+(x-b.Min.X)] = tile.(*image.Gray).GrayAt(x, y).Y
				}
			}

			oName := fmt.Sprintf(tileName, i, j, chanCodes[colour])
			err := WriteObject(bktName, oName, snappy.Encode(nil, pix))
			if err != nil {
				panic(err)
			}
		}
	}
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
