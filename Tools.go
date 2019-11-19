package main

import (
	"./Ruido"
	"./Vectores"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"math"
	"os"
)

// Cargo el weatherMap desde archivo para no tener que generarlo.
//
func loadWeatherMap() [WEATHER_WIDTH][WEATHER_HEIGHT]color.RGBA {

	var weatherMap [WEATHER_WIDTH][WEATHER_HEIGHT]color.RGBA

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
func genetareWeatherMap() {
	var ruido2D float64 = 0
	var color color.RGBA

	img := image.NewRGBA(image.Rect(0, 0, WEATHER_WIDTH, WEATHER_HEIGHT))
	out, err := os.Create("weatherMap.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for x := 0; x < WEATHER_WIDTH-1; x++ {
		for y := 0; y < WEATHER_HEIGHT-1; y++ {

			ruido2D = 255 * Ruido.Worley3D2(Vectores.Vector{X: float64(x) * HIGH_COVERAGE_NOISE, Y: float64(y) * HIGH_COVERAGE_NOISE, Z: 0})
			color.R = uint8(math.Abs(ruido2D))

			//ruido2D = Ruido.Worley3D(Vectores.Vector{X : float64(x) * HIGH_COVERAGE_NOISE, Y : float64(y) * HIGH_COVERAGE_NOISE, Z: 0})
			color.G = uint8(math.Abs(ruido2D))

			color.B = uint8(math.Abs(ruido2D))
			color.A = uint8(math.Abs(ruido2D))

			/*			ruido2D = Ruido.Noise2(float64(x) * LOW_COVERAGE_NOISE, float64(y) * LOW_COVERAGE_NOISE) * 255
						color.R = uint8(math.Abs(ruido2D))

						ruido2D = 255 *  Ruido.Worley3D(Vectores.Vector{X : float64(x) * HIGH_COVERAGE_NOISE, Y : float64(y) * HIGH_COVERAGE_NOISE, Z: 0})
						//ruido2D = Ruido.Worley3D(Vectores.Vector{X : float64(x) * HIGH_COVERAGE_NOISE, Y : float64(y) * HIGH_COVERAGE_NOISE, Z: 0})
						color.G = uint8(math.Abs(ruido2D))

						color.B = 120
						color.A = 120*/

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

func generateNoiseCube() {
	var noiseCube [NOISECUBE_X][NOISECUBE_Y][NOISECUBE_Z]byte

	for x := 0; x < NOISECUBE_X-1; x++ {
		for y := 0; y < NOISECUBE_Y-1; y++ {
			for z := 0; z < NOISECUBE_Z-1; z++ {
				noiseCube[x][y][z] = byte(255 * Ruido.Worley3D2(Vectores.Vector{X: float64(x) * HIGH_COVERAGE_NOISE, Y: float64(y) * HIGH_COVERAGE_NOISE, Z: float64(z) * HIGH_COVERAGE_NOISE}))
			}
		}
	}

	buffer := make([]byte, NOISECUBE_X*NOISECUBE_Y*NOISECUBE_Z)
	for x := 0; x < NOISECUBE_X-1; x++ {
		for y := 0; y < NOISECUBE_Y-1; y++ {
			for z := 0; z < NOISECUBE_Z-1; z++ {
				buffer[x+NOISECUBE_X*(y+NOISECUBE_Z*z)] = noiseCube[x][y][z]
			}
		}
	}

	err := ioutil.WriteFile("/home/john/go/src/Volumetric_Clouds_2/noiseCube.noi", buffer, 0644)
	if err != nil {
		panic(err)
	}

}

func generateNoiseCubeRGBA() {
	var noiseCube [NOISECUBE_X][NOISECUBE_Y][NOISECUBE_Z]color.RGBA

	for x := 0; x < NOISECUBE_X-1; x++ {
		for y := 0; y < NOISECUBE_Y-1; y++ {
			for z := 0; z < NOISECUBE_Z-1; z++ {
				noiseCube[x][y][z].R = byte(255 * Ruido.Worley3D2(Vectores.Vector{X: float64(x) * HIGH_COVERAGE_NOISE, Y: float64(y) * HIGH_COVERAGE_NOISE, Z: float64(z) * HIGH_COVERAGE_NOISE}))
				noiseCube[x][y][z].G = byte(255 * Ruido.Worley3D2(Vectores.Vector{X: float64(x) * HIGH_COVERAGE_NOISE, Y: float64(y) * HIGH_COVERAGE_NOISE, Z: float64(z) * HIGH_COVERAGE_NOISE}))
				noiseCube[x][y][z].B = byte(255 * Ruido.Worley3D2(Vectores.Vector{X: float64(x) * HIGH_COVERAGE_NOISE, Y: float64(y) * HIGH_COVERAGE_NOISE, Z: float64(z) * HIGH_COVERAGE_NOISE}))
				noiseCube[x][y][z].A = byte(255 * Ruido.Worley3D2(Vectores.Vector{X: float64(x) * HIGH_COVERAGE_NOISE, Y: float64(y) * HIGH_COVERAGE_NOISE, Z: float64(z) * HIGH_COVERAGE_NOISE}))
			}
		}
	}

	buffer := make([]byte, NOISECUBE_X*NOISECUBE_Y*NOISECUBE_Z)
	for x := 0; x < NOISECUBE_X-1; x++ {
		for y := 0; y < NOISECUBE_Y-1; y++ {
			for z := 0; z < NOISECUBE_Z-1; z++ {
				buffer[x+NOISECUBE_X*(y+NOISECUBE_Z*z)] = noiseCube[x][y][z]
			}
		}
	}

	err := ioutil.WriteFile("/home/john/go/src/Volumetric_Clouds_2/noiseCube.noi", buffer, 0644)
	if err != nil {
		panic(err)
	}

}

func loadCubeNoise() [NOISECUBE_X][NOISECUBE_Y][NOISECUBE_Z]byte {
	var noiseCube [NOISECUBE_X][NOISECUBE_Y][NOISECUBE_Z]byte

	buffer, err := ioutil.ReadFile("/home/john/go/src/Volumetric_Clouds_2/noiseCube.noi")
	if err != nil {
		panic(err)
	}

	for x := 0; x < NOISECUBE_X-1; x++ {
		for y := 0; y < NOISECUBE_Y-1; y++ {
			for z := 0; z < NOISECUBE_Z-1; z++ {
				noiseCube[x][y][z] = buffer[x+NOISECUBE_X*(y+NOISECUBE_Z*z)]
			}
		}
	}

	return noiseCube
}
