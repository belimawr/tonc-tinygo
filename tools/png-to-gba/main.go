package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"

	"github.com/belimawr/tonc-tinygo/types"
)

/*
Create an image on Gimp
To make it paletted: Image > Mode > Indexed
Select "Generate optimum palette"
Set maximum number of colours to 256
Export as PNG
*/

func main() {
	var (
		in   string
		mode string
		out  string
	)

	flag.Usage = func() {
		msg := `Convert PNG to GameBoy Advance images

Create a paletted PNG using Gimp:
 - Create a new image 240 x 160 pixels (to fill the whole screen) or smaller
 - Go to Image > Mode > Indexed
   - Select "Generate optimum palette"
   - Set "Maximum number of colours" to 256
   - Click "Convert"
 - When using the colour selection dialog, you will be able to choose colours
   and their index. Make sure the background is index 0

`
		fmt.Println(msg)

		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&in, "in", "img.png", "Input image")
	flag.StringVar(&out, "out", "out.dat", "Output file")
	flag.StringVar(&mode, "mode", "paletted", "Output format: paletted|bitmap|goBitmap")
	flag.Parse()

	f, err := os.ReadFile(in)
	if err != nil {
		log.Fatalf("cannot read image: %s", err)
	}

	imgReader := bytes.NewReader(f)

	img, err := png.Decode(imgReader)
	if err != nil {
		log.Fatalf("cannot decode png: %s", err)
	}

	imgP, ok := img.(*image.Paletted)
	if !ok {
		log.Fatal("not a paletted png")
	}

	switch mode {
	case "paletted":
		if err := writeImgToFile(out, img, paletted(imgP)); err != nil {
			log.Fatalf("cannot write paletted %s:", err)
		}
	case "bitmap":

		if err := writeImgToFile(out, img, bitmap(img)); err != nil {
			log.Fatalf("cannot write bitmap %s:", err)
		}

	case "goBitmap":
		if err := writeImgAsGoBitmap("go.dat", "main", "myImg", img); err != nil {
			log.Fatalf("cannot write paletted %s:", err)
		}
	default:
		log.Fatalf("mode %q not supported", mode)
	}

}

type ImgWriter func(w io.Writer, c color.Color) error

func writeImg(w io.Writer, img image.Image, fn ImgWriter) error {
	rect := img.Bounds()

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {

			rawC := img.At(x, y)
			if err := fn(w, rawC); err != nil {
				return fmt.Errorf("cannot write colour at (%d, %d): %w", x, y, err)
			}
		}
	}

	return nil
}

func writeImgToFile(filename string, img image.Image, fn ImgWriter) error {
	buff := bytes.Buffer{}

	writeImg(&buff, img, fn)

	if err := os.WriteFile(filename, buff.Bytes(), 0666); err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	return nil
}

func paletted(img *image.Paletted) func(w io.Writer, c color.Color) error {
	return func(w io.Writer, c color.Color) error {

		idx := uint8(img.Palette.Index(c))
		if err := binary.Write(w, binary.LittleEndian, idx); err != nil {
			return fmt.Errorf("cannot convert to bytes: %w", err)
		}

		return nil
	}
}

func bitmap(img image.Image) ImgWriter {
	return func(w io.Writer, c color.Color) error {
		gbaC := rgbaToColour(c)
		binary.Write(w, binary.LittleEndian, gbaC)
		return nil
	}
}

func writeImgAsGoBitmap(filename, pkg, varName string, img image.Image) error {
	buff := bytes.Buffer{}

	fmt.Fprintf(&buff, "package %s\n\n", pkg)
	fmt.Fprintf(&buff, "var %s = []byte{\n", varName)

	b := bytes.NewBuffer(make([]byte, 2))
	w := func(w io.Writer, c color.Color) error {
		b.Reset()

		gbaColour := rgbaToColour(c)
		if err := binary.Write(b, binary.LittleEndian, gbaColour); err != nil {
			return fmt.Errorf("cannot convert %X to binary: %s", c, err)
		}

		bb := b.Bytes()
		if _, err := fmt.Fprintf(&buff, "0x%02x, 0x%02x\n", bb[0], bb[1]); err != nil {
			return fmt.Errorf("cannot write to buffer: %w", err)
		}

		return nil
	}

	writeImg(&buff, img, w)

	fmt.Fprintln(&buff, "}")
	if err := os.WriteFile(filename, buff.Bytes(), 0666); err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	return nil
}

func rgbaToColour(rawC color.Color) types.Colour {
	r, g, b, _ := rawC.RGBA()
	r = min(r, 31)
	g = min(g, 31)
	b = min(b, 31)

	return types.NewColour(uint16(r), uint16(g), uint16(b))
}
