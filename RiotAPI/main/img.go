package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

type Changeable interface {
	Set(x, y int, c color.Color)
}


func imgtest (){
	reader, err := os.Open("map11.png")
	if err != nil {
		log.Fatal(err)
	}

	m, _, err := image.Decode(reader)
	defer reader.Close()

	if cimg, ok := m.(Changeable); ok {
		// cimg is of type Changeable, you can call its Set() method (draw on it)
		cimg.Set(0,0, color.RGBA{255, 0, 0, 255})
		cimg.Set(0,1, color.RGBA{255, 0, 0, 255} )
		cimg.Set(1,1, color.RGBA{255, 0, 0, 255} )
		cimg.Set(1,0, color.RGBA{255, 0, 0, 255} )

		output, _ := os.Create("outfile.png")
		tmp := cimg.(image.Image)
		fmt.Printf("%d\n", tmp.At(0,0 ))
		png.Encode(output, tmp)
		output.Close()
	}


}