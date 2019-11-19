package main

import (
	"./Vectores"
	"image/color"
	"math"
)

// Calcula el valor del weathermap con los canales R y G
//
func calculaWMC(wc0 float32, wc1 float32) float32 {
	return float32(math.Max(float64(wc0), float64(SAT(GC-0.5)*wc1*2)))
}

func main() {
	var weatherMap [WEATHER_WIDTH][WEATHER_HEIGHT]color.RGBA
	var NDC_x float32
	var NDC_y float32
	var PixelScreen_x float32
	var PixelScreen_y float32
	var PixelCamera_x float32
	var PixelCamera_y float32

	var ro Vectores.Vector
	var rd Vectores.Vector
	var nuevo Vectores.Vector
	var punto Vectores.Vector

	var t float32 = 10

	//var noiseCube [NOISECUBE_X][NOISECUBE_Y][NOISECUBE_Z] color.RGBA

	//start := time.Now()
	//genetareWeatherMap()
	//weatherMap = loadWeatherMap()
	//generateNoiseCubeRGBA()
	//noiseCube = loadCubeNoiseRGBA()
	//duration := time.Since(start)
	//noiseCube[0][0][0] = color.RGBA{0,0,0,0}
	// Formatted string, such as "2h3m0.5s" or "4.503μs"
	//fmt.Println(duration)
	//weatherMap[0][0].R = 0

	weatherMap = loadWeatherMap()

	// Calculo el Field of View. El ángulo es de 45 grados.
	//
	var FOV float32 = float32(math.Tan(float64(ALPHA / 2.0 * math.Pi / 180.0)))

	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			// Hacemos las conversiones de espacios
			//
			NDC_x = (float32(x) + correccion) / float32(WIDTH)
			NDC_y = (float32(y) + correccion) / float32(HEIGHT)

			PixelScreen_x = 2*NDC_x - 1
			PixelScreen_y = 2*NDC_y - 1

			PixelCamera_x = PixelScreen_x * ImageAspectRatio * FOV
			PixelCamera_y = PixelScreen_y * FOV

			// Origen y dirección

			ro = EYE
			nuevo.X = PixelCamera_x
			nuevo.Y = PixelCamera_y
			nuevo.Z = -1

			rd = nuevo.Sub(ro).Normalize()

			// Generamos puntos en la trayectoria del rayo hasta que alguno
			// tiene altura como para estar en la capa de nubes.
			//
			for {
				punto = ro.Add(rd.MultiplyByScalar(t))
				if punto.Y <= -HMIN || punto.Z <= -MAXHORIZON {
					break
				}
				t += 10
			}

		}
	}
	weatherMap[0][0].R = 0

}
