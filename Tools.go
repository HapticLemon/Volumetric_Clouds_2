package main

import (
	"./Ruido"
	"./Vectores"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

func loadWeatherMap() [WEATHER_WIDTH][WEATHER_HEIGHT] color.RGBA {

	var weatherMap [WEATHER_WIDTH][WEATHER_HEIGHT] color.RGBA

	imgfile, err := os.Open("./weatherMap.jpg")

	if err != nil {
		fmt.Println("img.jpg file not found!")
		os.Exit(1)
	}

	defer imgfile.Close()

	// get image height and width with image/jpeg
	// change accordinly if file is png or gif

	imgCfg, _, err := image.DecodeConfig(imgfile)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	width := imgCfg.Width
	height := imgCfg.Height

	fmt.Println("Width : ", width)
	fmt.Println("Height : ", height)

	// we need to reset the io.Reader again for image.Decode() function below to work
	// otherwise we will  - panic: runtime error: invalid memory address or nil pointer dereference
	// there is no build in rewind for io.Reader, use Seek(0,0)
	imgfile.Seek(0, 0)

	// get the image
	img, _, err := image.Decode(imgfile)

	fmt.Println(img.At(10, 10).RGBA())
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			weatherMap[x][y].R = uint8(r)
			weatherMap[x][y].G = uint8(g)
			weatherMap[x][y].B = uint8(b)
			weatherMap[x][y].A = uint8(a)

		}
	}

	return weatherMap
}

// Genera una imagen con el mapa del tiempo y la guarda en archivo para no tener que hacerlo
// cada vez, ya que Worley es MUY lento.
//
func genetareWeatherMap(){
	var  ruido2D float64 = 0
	var color color.RGBA

	img := image.NewRGBA(image.Rect(0, 0, WEATHER_WIDTH, WEATHER_HEIGHT))
	out, err := os.Create("weatherMap.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for x:= 0 ; x < WEATHER_WIDTH - 1 ; x++ {
		for y:= 0 ; y < WEATHER_HEIGHT - 1 ; y++ {

			ruido2D = Ruido.Noise2(float64(x) * LOW_COVERAGE_NOISE, float64(y) * LOW_COVERAGE_NOISE) * 255
			color.R = uint8(math.Abs(ruido2D))

			ruido2D = 255 *  Ruido.Worley3D(Vectores.Vector{X : float64(x) * HIGH_COVERAGE_NOISE, Y : float64(y) * HIGH_COVERAGE_NOISE, Z: 0})
			color.G = uint8(math.Abs(ruido2D))

			color.B = 120
			color.A = 120

			img.Set(x, y, color)
		}
	}
	var opt jpeg.Options

	opt.Quality = 80
	// ok, write out the data into the new JPEG file

	err = jpeg.Encode(out, img, &opt) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Generated image to %s \n", "weatherMap.jpg")


}

