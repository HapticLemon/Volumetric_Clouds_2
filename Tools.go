package main

import (
	"./Ruido"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"
)

// La clásica función de clamp pero con el nombre que aparece en el documento.
//
func SAT(valor float64) float64 {
	if valor < 0 {
		return 0.0
	}
	if valor > 1 {
		return 1.0
	}
	return valor
}

// Carga del shapeCube
//
func loadShapeCube() [SHAPECUBE_X][SHAPECUBE_Y][SHAPECUBE_Z]color.RGBA {
	var shapeCube [SHAPECUBE_X][SHAPECUBE_Y][SHAPECUBE_Z]color.RGBA

	bufferR, errR := ioutil.ReadFile("/home/john/go/src/noiseGenerator/shapeCubeR.dat")
	if errR != nil {
		panic(errR)
	}

	for x := 0; x < SHAPECUBE_X-1; x++ {
		for y := 0; y < SHAPECUBE_Y-1; y++ {
			for z := 0; z < SHAPECUBE_Z-1; z++ {
				shapeCube[x][y][z].R = bufferR[x+SHAPECUBE_X*(y+SHAPECUBE_Z*z)]
			}
		}
	}

	bufferG, errG := ioutil.ReadFile("/home/john/go/src/noiseGenerator/shapeCubeG.dat")
	if errG != nil {
		panic(errG)
	}

	for x := 0; x < SHAPECUBE_X-1; x++ {
		for y := 0; y < SHAPECUBE_Y-1; y++ {
			for z := 0; z < SHAPECUBE_Z-1; z++ {
				shapeCube[x][y][z].G = bufferG[x+SHAPECUBE_X*(y+SHAPECUBE_Z*z)]
			}
		}
	}

	bufferB, errB := ioutil.ReadFile("/home/john/go/src/noiseGenerator/shapeCubeB.dat")
	if errB != nil {
		panic(errB)
	}

	for x := 0; x < SHAPECUBE_X-1; x++ {
		for y := 0; y < SHAPECUBE_Y-1; y++ {
			for z := 0; z < SHAPECUBE_Z-1; z++ {
				shapeCube[x][y][z].B = bufferB[x+SHAPECUBE_X*(y+SHAPECUBE_Z*z)]
			}
		}
	}

	bufferA, errA := ioutil.ReadFile("/home/john/go/src/noiseGenerator/shapeCubeA.dat")
	if errA != nil {
		panic(errA)
	}

	for x := 0; x < SHAPECUBE_X-1; x++ {
		for y := 0; y < SHAPECUBE_Y-1; y++ {
			for z := 0; z < SHAPECUBE_Z-1; z++ {
				shapeCube[x][y][z].A = bufferA[x+SHAPECUBE_X*(y+SHAPECUBE_Z*z)]
			}
		}
	}

	return shapeCube
}

// Carga del detailCube
//
func loadDetailCube() [DETAILCUBE_X][DETAILCUBE_Y][DETAILCUBE_Z]color.RGBA {
	var detailCube [DETAILCUBE_X][DETAILCUBE_Y][DETAILCUBE_Z]color.RGBA

	bufferR, errR := ioutil.ReadFile("/home/john/go/src/noiseGenerator/detailCubeR.dat")
	if errR != nil {
		panic(errR)
	}

	for x := 0; x < DETAILCUBE_X-1; x++ {
		for y := 0; y < DETAILCUBE_Y-1; y++ {
			for z := 0; z < DETAILCUBE_Z-1; z++ {
				detailCube[x][y][z].R = bufferR[x+DETAILCUBE_X*(y+DETAILCUBE_Z*z)]
			}
		}
	}

	bufferG, errG := ioutil.ReadFile("/home/john/go/src/noiseGenerator/detailCubeG.dat")
	if errG != nil {
		panic(errG)
	}

	for x := 0; x < DETAILCUBE_X-1; x++ {
		for y := 0; y < DETAILCUBE_Y-1; y++ {
			for z := 0; z < DETAILCUBE_Z-1; z++ {
				detailCube[x][y][z].G = bufferG[x+DETAILCUBE_X*(y+DETAILCUBE_Z*z)]
			}
		}
	}

	bufferB, errB := ioutil.ReadFile("/home/john/go/src/noiseGenerator/detailCubeB.dat")
	if errB != nil {
		panic(errB)
	}

	for x := 0; x < DETAILCUBE_X-1; x++ {
		for y := 0; y < DETAILCUBE_Y-1; y++ {
			for z := 0; z < DETAILCUBE_Z-1; z++ {
				detailCube[x][y][z].B = bufferB[x+DETAILCUBE_X*(y+DETAILCUBE_Z*z)]
			}
		}
	}

	bufferA, errA := ioutil.ReadFile("/home/john/go/src/noiseGenerator/detailCubeA.dat")
	if errA != nil {
		panic(errA)
	}

	for x := 0; x < DETAILCUBE_X-1; x++ {
		for y := 0; y < DETAILCUBE_Y-1; y++ {
			for z := 0; z < DETAILCUBE_Z-1; z++ {
				detailCube[x][y][z].A = bufferA[x+DETAILCUBE_X*(y+DETAILCUBE_Z*z)]
			}
		}
	}

	return detailCube
}

/*
// Genera una imagen con el mapa del tiempo y la guarda en archivo para no tener que hacerlo
// cada vez, ya que Worley es MUY lento.
//
func generateWeatherMap() {
	var colorR color.RGBA
	var color color.RGBA


	img := image.NewRGBA(image.Rect(0, 0, WEATHER_X, WEATHER_Y))
	out, err := os.Create("./weatherMap.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	imgR := image.NewRGBA(image.Rect(0, 0, WEATHER_X, WEATHER_Y))
	outR, errR := os.Create("./weatherMapR.jpg")
	if errR != nil {
		fmt.Println(errR)
		os.Exit(1)
	}
	for x := 0; x < WEATHER_X-1; x++ {
		for y := 0; y < WEATHER_Y-1; y++ {
			if (x >= 260 && x <=280) && (y >= 210 && y <=230) {
				fmt.Println("gaatete")
			}
			color.R = byte(255 * SAT(Ruido.Noise2(float64(x) * LOW_FREQ_NOISE,float64(y) * LOW_FREQ_NOISE)))
			color.G = byte(255 *  Ruido.Worley3D(Vectores.Vector{X : float64(x) * MEDIUM_FREQ_NOISE, Y : float64(y) * MEDIUM_FREQ_NOISE, Z: 0}))

			color.B = 120
			color.A = 12

			colorR.R = color.R
			colorR.G = color.R
			colorR.B = color.R
			colorR.A = 255

			img.Set(x, y, color)
			imgR.Set(x, y, colorR)
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
	errR = jpeg.Encode(outR, imgR, &opt) // put quality to 80%
	if errR != nil {
		fmt.Println(errR)
		os.Exit(1)
	}
	fmt.Printf("Generated image to %s \n", "./weatherMap.jpg")

}
*/
// Cargo el weatherMap desde archivo para no tener que generarlo.
//
/*
func loadWeatherMap() [WEATHER_X][WEATHER_Y]color.RGBA {

	var weatherMap [WEATHER_X][WEATHER_Y]color.RGBA

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
	for x := 0; x < WEATHER_X - 1; x++ {
		for y := 0; y < WEATHER_Y - 1; y++ {
			if (x >= 260 && x <=280) && (y >= 210 && y <=230) {
				fmt.Println("gaatete")
			}
			r, g, b, a := img.At(x, y).RGBA()
			r, g, b, a = r>>8, g>>8, b>>8, a>>8
			weatherMap[x][y].R = uint8(r)
			weatherMap[x][y].G = uint8(g)
			weatherMap[x][y].B = uint8(b)
			weatherMap[x][y].A = uint8(a)

		}
	}

	return weatherMap
}
*/

// Genera una imagen con el mapa del tiempo y la guarda en archivo para no tener que hacerlo
// cada vez, ya que Worley es MUY lento.
//
func generateWeatherMap() {
	var weatherMap [WEATHER_X][WEATHER_Y]color.RGBA
	var color color.RGBA

	img := image.NewRGBA(image.Rect(0, 0, WEATHER_X, WEATHER_Y))
	out, err := os.Create("/home/john/go/src/noiseGenerator/weatherMap.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p := Ruido.NewPerlin(alpha, beta, n, seed)

	for x := 0; x < WEATHER_X-1; x++ {
		for y := 0; y < WEATHER_Y-1; y++ {
			weatherMap[x][y].R = byte(255 * SAT(p.Noise2D(float64(x)*TEST_FREQU_NOISE, float64(y)*TEST_FREQU_NOISE)))
			weatherMap[x][y].G = byte(255 * SAT(p.Noise2D(float64(x)*HIGH_FREQ_NOISE, float64(y)*HIGH_FREQ_NOISE)))
			weatherMap[x][y].B = 100 //byte(255 * SAT(p.Noise2D(float64(x) * LOW_FREQ_NOISE,float64(y) * LOW_FREQ_NOISE)))
			weatherMap[x][y].A = 11

			color.R = weatherMap[x][y].R
			color.G = weatherMap[x][y].G
			color.B = weatherMap[x][y].B
			color.A = weatherMap[x][y].A

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

	fmt.Printf("Generated image to %s \n", "/home/john/go/src/noiseGenerator/weatherMap.jpg")

	// Lo grabo en 4 archivos en lugar de en uno. Es un poco guarro pero funciona.
	//
	bufferR := make([]byte, WEATHER_X*WEATHER_Y)
	bufferG := make([]byte, WEATHER_X*WEATHER_Y)
	bufferB := make([]byte, WEATHER_X*WEATHER_Y)
	bufferA := make([]byte, WEATHER_X*WEATHER_Y)

	for x := 0; x < WEATHER_X-1; x++ {
		for y := 0; y < WEATHER_Y-1; y++ {
			bufferR[WEATHER_X*x+y] = weatherMap[x][y].R
			bufferG[WEATHER_X*x+y] = weatherMap[x][y].G
			bufferB[WEATHER_X*x+y] = weatherMap[x][y].B
			bufferA[WEATHER_X*x+y] = weatherMap[x][y].A
		}
	}

	errR := ioutil.WriteFile("/home/john/go/src/noiseGenerator/weatherMapR.dat", bufferR, 0644)
	if errR != nil {
		panic(errR)
	}
	errG := ioutil.WriteFile("/home/john/go/src/noiseGenerator/weatherMapG.dat", bufferG, 0644)
	if errG != nil {
		panic(errG)
	}
	errB := ioutil.WriteFile("/home/john/go/src/noiseGenerator/weatherMapB.dat", bufferB, 0644)
	if errB != nil {
		panic(errB)
	}
	errA := ioutil.WriteFile("/home/john/go/src/noiseGenerator/weatherMapA.dat", bufferA, 0644)
	if errA != nil {
		panic(errA)
	}
}

// Carga del detailCube
//
func loadWeatherMap() [WEATHER_X][WEATHER_Y]color.RGBA {
	var weatherMap [WEATHER_X][WEATHER_Y]color.RGBA

	bufferR, errR := ioutil.ReadFile("/home/john/go/src/noiseGenerator/weatherMapR.dat")
	if errR != nil {
		panic(errR)
	}

	for x := 0; x < WEATHER_X-1; x++ {
		for y := 0; y < WEATHER_Y-1; y++ {
			weatherMap[x][y].R = bufferR[WEATHER_X*x+y]
		}
	}

	bufferG, errG := ioutil.ReadFile("/home/john/go/src/noiseGenerator/weatherMapG.dat")
	if errG != nil {
		panic(errG)
	}

	for x := 0; x < WEATHER_X-1; x++ {
		for y := 0; y < WEATHER_Y-1; y++ {
			weatherMap[x][y].G = bufferG[WEATHER_X*x+y]
		}
	}

	bufferB, errB := ioutil.ReadFile("/home/john/go/src/noiseGenerator/weatherMapB.dat")
	if errB != nil {
		panic(errB)
	}

	for x := 0; x < WEATHER_X-1; x++ {
		for y := 0; y < WEATHER_Y-1; y++ {
			weatherMap[x][y].B = bufferB[WEATHER_X*x+y]
		}
	}

	bufferA, errA := ioutil.ReadFile("/home/john/go/src/noiseGenerator/weatherMapA.dat")
	if errA != nil {
		panic(errA)
	}

	for x := 0; x < WEATHER_X-1; x++ {
		for y := 0; y < WEATHER_Y-1; y++ {
			weatherMap[x][y].A = bufferA[WEATHER_X*x+y]
		}
	}

	return weatherMap
}

// Función de remap para cambiar el espectro de valores.
//
func R(v float64, l0 float64, h0 float64, ln float64, hn float64) float64 {
	return ln + (((v - l0) * (hn - ln)) / (h0 - l0))
}

// Interpolación lineal entre dos valores.
//
func Li(v0 float64, v1 float64, ival float64) float64 {
	return (1-ival)*v0 + ival*v1
}

// Interpolación entre dos colores.
//
func mixColor(x color.RGBA, y color.RGBA, a float64) color.RGBA {
	var resultado color.RGBA

	resultado.R = uint8(float64(x.R)*(1-a) + float64(y.R)*a)
	resultado.G = uint8(float64(x.G)*(1-a) + float64(y.G)*a)
	resultado.B = uint8(float64(x.B)*(1-a) + float64(y.B)*a)

	return resultado
}
