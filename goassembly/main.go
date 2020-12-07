package main

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"syscall/js"
)

func main() {
	c := make(chan bool)
	//1. Adding an <h1> element in the HTML document
	document := js.Global().Get("document")
	p := document.Call("createElement", "h1")
	p.Set("innerHTML", "Hello from Golang!")
	document.Get("body").Call("appendChild", p)

	//2. Exposing go functions/values in javascript variables.
	js.Global().Set("goVar", "I am a variable set from Go")
	js.Global().Set("sayHello", js.FuncOf(sayHello))

	js.Global().Set("changeImage", js.FuncOf(changeImage))

	//3. This channel will prevent the go program to exit
	<-c
}

func sayHello(this js.Value, inputs []js.Value) interface{} {
	firstArg := inputs[0].String()
	return "Hi " + firstArg + " from Go!"
}

func changeImage(this js.Value, inputs []js.Value) interface{} {
	firstArg := inputs[0]

	inBuf := make([]uint8, firstArg.Get("byteLength").Int())
	js.CopyBytesToGo(inBuf, firstArg)

	img, _, err := image.Decode(bytes.NewReader(inBuf))

	if err != nil {
		return err.Error()
	}
	// var point image.Point
	// point = img.Bounds().Max()

	cimg := image.NewRGBA(img.Bounds())

	for x := 0; x < img.Bounds().Max.X; x++ {
		for y := 0; y < img.Bounds().Max.Y; y++ {
			// r, g, b, a := img.At(x, y).RGBA()

			// prevColor := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			// if img.Bounds().Max.X > x && x > 1 &&
			// 	img.Bounds().Max.Y > y && y > 1 {
			// 	prevColor = color.RGBAModel.Convert(img.At(x-1, y-1)).(color.RGBA)
			// }
			originalColor, _ := color.RGBAModel.Convert(img.At(x, y)).(color.RGBA)

			// grey := uint8(float64(originalColor.R)*0.21 + float64(originalColor.G)*0.72 + float64(originalColor.B)*0.07)

			cimg.Set(x, y, color.RGBA{
				R: originalColor.R / 2,
				G: originalColor.G / 5,
				B: originalColor.B / 5,
				A: originalColor.A,
			})
		}
	}

	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, cimg, nil)
	newArray := buf.Bytes()

	// this.Global().Set("newImageArrayByte", newArray)

	// r, g, b, a := img.At(20, 20).RGBA()
	dst := js.Global().Get("Uint8Array").New(len(newArray))
	js.CopyBytesToJS(dst, newArray)

	return dst

}
