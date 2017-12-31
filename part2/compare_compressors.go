package main

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/lzw"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"time"

	"github.com/golang/snappy"
	"github.com/pierrec/lz4"
)

const (
	fileName = "world.topo.bathy.200412.3x21600x10800.png"
)

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

func PNGWriter(fName string, gray *image.Gray) error {
	start := time.Now()
	f, err := os.Create(fName + ".png")
	if err != nil {
		return err
	}

	err = png.Encode(f, gray)
	fmt.Printf("Writting PNG File to disk: %v %v\n", time.Since(start), err)
	return err
}

func PNGReader(fName string) error {
	start := time.Now()
	data, err := os.Open(fName + ".png")
	if err != nil {
		return err
	}
	_, err = png.Decode(data)
	fmt.Printf("Reading PNG File from disk: %v\n", time.Since(start))
	return err
}

func RawWriter(fName string, data []byte) error {
	start := time.Now()

	err := ioutil.WriteFile(fName+".raw", data, 0644)
	fmt.Printf("Writting Raw File to disk: %v\n", time.Since(start))

	return err
}

func RawReader(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName + ".raw")
	fmt.Printf("Reading Raw File from disk: %v\n", time.Since(start))

	return data, err
}

func FlateWriter(fName string, data []byte, level int) error {
	start := time.Now()

	outFile, err := os.Create(fName + ".flt")
	if err != nil {
		return err
	}

	flateWriter, err := flate.NewWriter(outFile, level)
	if err != nil {
		return err
	}
	_, err = flateWriter.Write(data)
	if err != nil {
		return err
	}
	err = flateWriter.Close()

	fmt.Printf("Writting Flate File to disk: %v\n", time.Since(start))
	return err
}

func FlateReader(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName + ".flt")
	if err != nil {
		return []byte{}, err
	}

	buf := bytes.NewBuffer(data)

	flateReader := flate.NewReader(buf)

	var resB bytes.Buffer
	_, err = resB.ReadFrom(flateReader)
	if err != nil {
		return []byte{}, err
	}

	fmt.Printf("Reading Flate File from disk: %v\n", time.Since(start))

	return resB.Bytes(), nil
}

func LZWWriter(fName string, data []byte) error {
	start := time.Now()

	outFile, err := os.Create(fName + ".lzw")
	if err != nil {
		return err
	}

	lzwWriter := lzw.NewWriter(outFile, lzw.LSB, 8)

	if _, err = lzwWriter.Write(data); err != nil {
		return err
	}

	if err = lzwWriter.Close(); err != nil {
		return err
	}

	fmt.Printf("Writting LZW File to disk: %v\n", time.Since(start))

	return nil
}

func LZWReader(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName + ".lzw")
	if err != nil {
		return []byte{}, err
	}

	buf := bytes.NewBuffer(data)

	lzwReader := lzw.NewReader(buf, lzw.LSB, 8)

	var resB bytes.Buffer
	_, err = resB.ReadFrom(lzwReader)
	if err != nil {
		return []byte{}, err
	}

	fmt.Printf("Reading LZW File from disk: %v\n", time.Since(start))

	return resB.Bytes(), nil
}

func GZipWriter(fName string, data []byte) error {
	start := time.Now()

	outFile, err := os.Create(fName + ".gzip")
	if err != nil {
		return err
	}

	gzipWriter := gzip.NewWriter(outFile)

	if _, err = gzipWriter.Write(data); err != nil {
		return err
	}

	if err = gzipWriter.Close(); err != nil {
		return err
	}

	fmt.Printf("Writting Zip File to disk: %v\n", time.Since(start))

	return nil
}

func GZipReader(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName + ".gzip")
	if err != nil {
		return []byte{}, err
	}

	buf := bytes.NewBuffer(data)

	gzipReader, err := gzip.NewReader(buf)
	if err != nil {
		return []byte{}, err
	}

	var resB bytes.Buffer
	_, err = resB.ReadFrom(gzipReader)
	if err != nil {
		return []byte{}, err
	}

	fmt.Printf("Reading Zip File from disk: %v\n", time.Since(start))

	return resB.Bytes(), nil
}

func LZ4Writer(fName string, data []byte) error {
	start := time.Now()

	comp := make([]byte, len(data))

	l, err := lz4.CompressBlock(data, comp, 0)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(fName+".lz4", comp[:l], 0644)
	fmt.Printf("Writting LZ4 File to disk: %v\n", time.Since(start))

	return err
}

func LZ4Reader(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName + ".lz4")

	decomp := make([]byte, len(data)*3)
	l, err := lz4.UncompressBlock(data, decomp, 0)
	if err != nil {
		return []byte{}, err
	}
	fmt.Printf("Reading LZ4 File from disk: %v\n", time.Since(start))

	return decomp[:l], err
}

func SnappyWriter(fName string, data []byte) error {
	start := time.Now()

	err := ioutil.WriteFile(fName+".snp", snappy.Encode(nil, data), 0644)
	fmt.Printf("Writting Snappy File to disk: %v\n", time.Since(start))

	return err
}

func SnappyReader(fName string) ([]byte, error) {
	start := time.Now()

	data, err := ioutil.ReadFile(fName + ".snp")
	cdata, err := snappy.Decode(nil, data)
	fmt.Printf("Reading Snappy File from disk: %v\n", time.Since(start))

	return cdata, err
}

func main() {
	data, _ := os.Open(fileName)
	img, _ := png.Decode(data)
	bands := GetChannels(img)

	fNames := []string{"red", "green", "blue"}

	for i := 0; i < 3; i++ {
		fmt.Println(fNames[i])
		PNGWriter(fNames[i], bands[i])
		RawWriter(fNames[i], bands[i].Pix)
		FlateWriter(fNames[i], bands[i].Pix, 1)
		LZWWriter(fNames[i], bands[i].Pix)
		GZipWriter(fNames[i], bands[i].Pix)
		LZ4Writer(fNames[i], bands[i].Pix)
		SnappyWriter(fNames[i], bands[i].Pix)

		PNGReader(fNames[i])
		RawReader(fNames[i])
		FlateReader(fNames[i])
		LZWReader(fNames[i])
		GZipReader(fNames[i])
		LZ4Reader(fNames[i])
		SnappyReader(fNames[i])
	}
}
