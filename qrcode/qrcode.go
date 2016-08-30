package main

import (
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"image"
	"image/png"
	"log"
	"os"
)

func main() {
	str := `{"Status":true,"IDDealer":135}`

	code, err := qr.Encode(str, qr.L, qr.Unicode)

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Encoded data: ", code.Content())

	if base64 != code.Content() {
		log.Fatal("data did not match")
	}

	code, err = barcode.Scale(code, 90, 90)

	if err != nil {
		log.Fatal(err)
	}

	writePNG("static/test1.png", code)
}

func writePNG(filename string, img image.Image) {
	file, err := os.Create(filename)

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	err = png.Encode(file, img)

	if err != nil {
		log.Fatal(err)
	}
}
