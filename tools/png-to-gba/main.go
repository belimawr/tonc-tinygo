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
	"strings"

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

	switch mode {
	case "paletted":
		imgP, ok := img.(*image.Paletted)
		if !ok {
			log.Fatal("not a paletted png")
		}

		if len(imgP.Palette) > 256 {
			log.Fatalf("the maximum palette size is 256, current palette is %d", len(imgP.Palette))
		}

		idx := strings.Index(out, ".")
		pName := out[:]
		if idx > 0 {
			pName = out[:idx]
		}
		pName += ".palette"

		if err := writePaletteToFile(pName, imgP); err != nil {
			log.Fatalf("cannot write palette to disk %s:", err)
		}

		if err := writeImgToFile(out, img, paletted(imgP.Palette)); err != nil {
			log.Fatalf("cannot write paletted %s:", err)
		}

	case "bitmap":
		if err := writeImgToFile(out, img, bitmap); err != nil {
			log.Fatalf("cannot write bitmap %s:", err)
		}

	case "goBitmap":
		if err := writeImgAsGoBitmap(out, "main", "myImg", img); err != nil {
			log.Fatalf("cannot write paletted %s:", err)
		}

	default:
		log.Fatalf("mode %q not supported", mode)
	}
}

func writePaletteToFile(filename string, img *image.Paletted) error {
	buff := bytes.NewBuffer(make([]byte, 0, len(img.Palette)*2))
	for _, c := range img.Palette {
		gbaC := rgbaToColour(c)
		if err := binary.Write(buff, binary.LittleEndian, gbaC); err != nil {
			return fmt.Errorf("cannot convert to bytes: %w", err)
		}
	}

	return os.WriteFile(filename, buff.Bytes(), 0666)
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

func paletted(palette color.Palette) func(w io.Writer, c color.Color) error {
	return func(w io.Writer, c color.Color) error {
		idx := uint8(palette.Index(c))
		if err := binary.Write(w, binary.LittleEndian, idx); err != nil {
			return fmt.Errorf("cannot convert to bytes: %w", err)
		}

		return nil
	}
}

func bitmap(w io.Writer, c color.Color) error {
	gbaC := rgbaToColour(c)
	binary.Write(w, binary.LittleEndian, gbaC)
	return nil
}

func writeImgAsGoBitmap(filename, pkg, varName string, img image.Image) error {
	idx := strings.Index(filename, ".")
	if idx > 0 {
		filename = filename[:idx]
	}
	filename += ".go"

	buff := bytes.Buffer{}

	fmt.Fprintf(&buff, "package %s\n\n", pkg)
	fmt.Fprintf(&buff, "var %s = []uint16{\n", varName)

	i := 0
	w := func(w io.Writer, c color.Color) error {
		gbaColour := rgbaToColour(c)
		fmt.Fprintf(&buff, "0x%04X, ", gbaColour)
		i++
		if i%8 == 0 {
			fmt.Fprint(&buff, "\n")
		}
		return nil
	}

	writeImg(&buff, img, w)

	fmt.Fprintln(&buff, "\n}")
	if err := os.WriteFile(filename, buff.Bytes(), 0666); err != nil {
		return fmt.Errorf("cannot write file: %w", err)
	}

	return nil
}

func rgbaToColour(rawC color.Color) types.Colour {
	// RGBA() returns 0-65535
	r16, g16, b16, _ := rawC.RGBA()

	// Shift right by 11 to map [0-65535] down to [0-31]
	r := uint16(r16 >> 11)
	g := uint16(g16 >> 11)
	b := uint16(b16 >> 11)

	return types.NewColour(r, g, b)
}
