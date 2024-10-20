package painter

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type Pixel struct {
	R int
	G int
	B int
	A int
} //TIP See GoLand help at <a href="https://www.jetbrains.com/help/go/">jetbrains.com/help/go/</a>.
type MyPixel struct {
	Y     int
	X     int
	Color string
}

func New(path string) [][]MyPixel {
	//TIP Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined or highlighted text
	// to see how GoLand suggests fixing it.
	return openIMG(path)
}

func openIMG(path string) [][]MyPixel {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()
	_, _, err = image.Decode(file)
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println(imageData)

	//fmt.Println(imageType)
	file.Seek(0, 0)
	loadedImage, err := png.Decode(file)
	if err != nil {
		fmt.Println(err)
	}
	all_pixels := getPixels(loadedImage, 16)
	return all_pixels
}

func toHex(number int) string {
	Newnumber := number
	result := ""
	convertToHex := map[int]string{
		0:  "0",
		1:  "1",
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "A",
		11: "B",
		12: "C",
		13: "D",
		14: "E",
		15: "F",
	}
	for Newnumber != 0 {
		tmp := Newnumber % 16

		result = convertToHex[tmp] + result
		Newnumber /= 16
	}
	if len(result) < 1 {
		result = "0" + result
	}
	if len(result) < 2 {
		result = "0" + result
	}
	return result
}
func getPixels(img image.Image, scale int) [][]MyPixel {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	var pixels [][]MyPixel
	for y := 0; y < height; y += scale {
		midY := y + scale
		var row []MyPixel
		for x := 0; x < width; x += scale {
			midX := x + scale
			pixel := rgbaToPixel(img.At(midX, midY).RGBA())
			hexCode := "#" + toHex(pixel.R) + toHex(pixel.G) + toHex(pixel.B)
			row = append(row, MyPixel{X: x / scale, Y: y / 16, Color: hexCode})
		}
		pixels = append(pixels, row)
	}
	return pixels
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) Pixel {
	return Pixel{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

// Pixel struct example
