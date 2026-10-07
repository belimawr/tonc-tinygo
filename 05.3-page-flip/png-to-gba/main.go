package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"
	_ "image/png"
	"io"
	"log"
	"os"
	"strconv"

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
	f, err := os.ReadFile("img.png")
	if err != nil {
		log.Fatalf("cannot read image: %s", err)
	}

	imgReader := bytes.NewReader(f)

	img, err := png.Decode(imgReader)
	if err != nil {
		log.Fatalf("cannot decode png: %s", err)
	}

	fmt.Println("[0,0]:", img.At(0, 0))
	fmt.Println("[0,1]:", img.At(1, 0))
	fmt.Println("[0,2]:", img.At(2, 0))
	fmt.Println("[0,3]:", img.At(3, 0))
	fmt.Println("[0,4]:", img.At(4, 0))
	fmt.Printf("Img: %T\n", img)
	p, ok := img.(*image.Paletted)
	if !ok {
		log.Fatal("not a paletted png")
	}

	fmt.Println("Palette: ", p.Palette)
	fmt.Println("Rect: ", p.Rect)
	fmt.Println("Stride: ", p.Stride)
	fmt.Println("Pix[:10]: ", p.Pix[:10])
	fmt.Println("foo", p.Palette.Index(img.At(0, 0)))
	fmt.Println("Palette[0,0]:", p.Palette.Index(img.At(0, 0)))
	fmt.Println("Palette[0,1]:", p.Palette.Index(img.At(1, 0)))
	fmt.Println("Palette[0,2]:", p.Palette.Index(img.At(2, 0)))
	fmt.Println("Palette[0,3]:", p.Palette.Index(img.At(3, 0)))
	fmt.Println("Palette[0,4]:", p.Palette.Index(img.At(4, 0)))

	r, g, b, _ := p.Palette[0].RGBA()
	c := types.NewColour(uint16(r), uint16(g), uint16(b))
	fmt.Printf("formated: 0x%s\n", strconv.FormatUint(uint64(c), 16))

	bb := make([]byte, 0, 10)
	buff := bytes.NewBuffer(bb)

	for i := range 5 {
		r, g, b, _ := img.At(i, 0).RGBA()
		r = min(r, 31)
		g = min(g, 31)
		b = min(b, 31)
		c := types.NewColour(uint16(r), uint16(g), uint16(b))
		fmt.Printf("0b%016b\n", c)
		binary.Write(buff, binary.LittleEndian, c)
	}

	if err := os.WriteFile("bin.dat", buff.Bytes(), 0666); err != nil {
		log.Fatalf("cannot write dat file: %s", err)
	}

	writeImg(os.Stdout, buff.Bytes(), "main", "myImg")
}

func writeImg(w io.Writer, data []byte, pkg, name string) {
	fmt.Fprintf(w, "package %s\n\n", pkg)
	fmt.Fprintf(w, "var %s = {\n", name)
	for _, d := range data {
		fmt.Fprintf(w, "0x%02x,\n", d)
	}
	fmt.Fprintln(w, "}")
}
