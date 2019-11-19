package main

import (
	"./Vectores"
	"image/color"
	"math"
)

// Calcula el valor del weathermap con los canales R y G
//
func calculaWMC(wc0 uint8, wc1 uint8) float32 {
	return float32(math.Max(float64(wc0), float64(SAT(GC-0.5)*float32(wc1)*2)))
}

// Calculo el porcentaje de la altura a la que estamos en la zona de nubes [0-1]
//
func calculaPH(h float32) float32 {
	return (h - HMIN) / HINTERVAL
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

	var xremap uint8 = 0
	var zremap uint8 = 0

	var wmc float32 = 0.0
	var ph float32 = 0.0

	// Shape Round top y bottom
	var srb float32
	var srt float32
	var sa float32

	// Density Round top y bottom
	var drb float32
	var drt float32
	var da float32

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
				// Si nos salimos de los límites pasamos al siguiente punto.
				// También si el rayo no apunta hacia el cielo, ya que nunca tocará la capa de nubes.
				//
				if punto.Z <= -MAXHORIZON || punto.X > MAXX || punto.X < MINX || punto.Y >= 0 {
					break
				}
				t += 10
				if punto.Y <= -HMIN {
					// Coordenadas remapeadas para acceder al weeathermap.
					// Creo que no tendría que haber problema con los valores - de X/Z ya que siempre se mapean a +
					//
					xremap = uint8(R(punto.X, MINX, MAXX, 0, float32(WEATHER_WIDTH)))
					zremap = uint8(R(punto.Z, 0, MAXHORIZON, 0, float32(WEATHER_HEIGHT)))

					wmc = calculaWMC(weatherMap[xremap][zremap].R, weatherMap[xremap][zremap].G)

					// Ojo, las alturas van en negativo.
					ph = calculaPH(-punto.Y)

					// Inicio de Shape-altering height-function--------------------------
					srb = SAT(R(ph, 0, 0.07, 0, 1))
					srt = R(ph, float32(weatherMap[xremap][zremap].B)*0.2, float32(weatherMap[xremap][zremap].B), 1, 0)
					sa = srb * srt
					// Fin de Shape-altering height-function-----------------------------

					// Inicio de Density-altering height-function +++++++++++++++++++++++
					drb = ph * SAT(R(ph, 0, 0.15, 0, 1))
					drt = SAT(R(ph, 0.9, 1.0, 1, 0))
					da = GD * drb * drt * float32(weatherMap[xremap][zremap].A) * 2
					// Fin de Density-altering height-function ++++++++++++++++++++++++++

					// Carroña para que el compilador no proteste.

					wmc *= 2
					ph *= 2
					sa *= 2
				}
			}

		}
	}
	weatherMap[0][0].R = 0

}
