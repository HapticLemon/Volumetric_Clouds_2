package main

import (
	"./Vectores"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
	"time"
)

// Calcula el valor del weathermap con los canales R y G
//
func calculaWMC(wc0 uint8, wc1 uint8) float64 {
	return math.Max(float64(wc0), SAT(GC-0.5)*float64(wc1)*2)
}

// Calculo el porcentaje de la altura a la que estamos en la zona de nubes [0-1]
//
func calculaPH(h float64) float64 {
	return (h - HMIN) / HINTERVAL
}

func calculaDensidad(ro Vectores.Vector, rd Vectores.Vector, weatherMap [WEATHER_X][WEATHER_Y]color.RGBA, shapeCube [SHAPECUBE_X][SHAPECUBE_Y][SHAPECUBE_Z]color.RGBA, detailCube [DETAILCUBE_X][DETAILCUBE_Y][DETAILCUBE_Z]color.RGBA) float64 {
	var punto Vectores.Vector
	var xremap int = 0
	var zremap int = 0

	// Nos servirán de coordenadas para el cubeNoise
	var xShapeCube int = 0
	var yShapeCube int = 0
	var zShapeCube int = 0

	// Nos servirán de coordenadas para el cubeNoise
	var xDetailCube int = 0
	var yDetailCube int = 0
	var zDetailCube int = 0

	var wmc float64 = 0.0
	var ph float64 = 0.0

	// Shape Round top y bottom
	var srb float64
	var srt float64
	var sa float64

	// Density Round top y bottom
	var drb float64
	var drt float64
	var da float64

	// Shape and detail noise.
	// Shape
	var snsample float64
	//var sn float32
	var snr byte
	var sng byte
	var snb byte
	var sna byte
	// Detail
	var dnfbm float64
	var dnmod float64
	var snnd float64
	var dnr byte
	var dng byte
	var dnb byte

	var t float64 = 10
	var densidad float64 = 0

	// Generamos puntos en la trayectoria del rayo hasta que alguno
	// tiene altura como para estar en la capa de nubes.
	//
	for {
		punto = ro.Add(rd.MultiplyByScalar(t))
		// Si nos salimos de los límites pasamos al siguiente punto.
		// También si el rayo no apunta hacia el cielo, ya que nunca tocará la capa de nubes.
		//
		if punto.Z <= -MAXHORIZON || punto.X > MAXX || punto.X < MINX || punto.Y >= 0 || densidad >= 1 {
			break
		}
		t += 1
		// Limitamos a la banda de alturas en la que hay nubes.
		if punto.Y <= -HMIN && punto.Y >= -HMAX {
			// Coordenadas remapeadas para acceder al weeathermap.
			// Creo que no tendría que haber problema con los valores - de X/Z ya que siempre se mapean a +
			//
			//xremap = int(R(punto.X, MINX, MAXX, 0, float64(WEATHER_X)))
			//zremap = int(R(punto.Z, 0, MAXHORIZON, 0, float64(WEATHER_Y)))

			xremap = int(math.Abs(R(punto.X, MINX, MAXX, 0, float64(WEATHER_X))))
			zremap = int(math.Abs(R(punto.Z, 0, MAXHORIZON, 0, float64(WEATHER_Y))))

			if (xremap >= 260 && xremap <= 280) && (zremap >= 210 && zremap <= 230) {
				fmt.Println("gaatete")
			}
			var componenteR = weatherMap[xremap][zremap].R
			//return float64(weatherMap[xremap][zremap].R)
			return float64(componenteR)

			wmc = calculaWMC(weatherMap[xremap][zremap].R, weatherMap[xremap][zremap].R)

			// Ojo, las alturas van en negativo.
			ph = calculaPH(-punto.Y)

			// Inicio de Shape-altering height-function--------------------------
			srb = SAT(R(ph, 0, 0.07, 0, 1))
			srt = R(ph, float64(weatherMap[xremap][zremap].B)*0.2, float64(weatherMap[xremap][zremap].B), 1, 0)
			sa = srb * srt
			// Fin de Shape-altering height-function-----------------------------

			// Inicio de Density-altering height-function +++++++++++++++++++++++
			drb = ph * SAT(R(ph, 0, 0.15, 0, 1))
			drt = SAT(R(ph, 0.9, 1.0, 1, 0))
			da = GD * drb * drt * float64(weatherMap[xremap][zremap].A) * 2
			// Fin de Density-altering height-function ++++++++++++++++++++++++++

			// Shape and detail noise********************************************
			// Shape
			xShapeCube = int(R(punto.X, MINX, MAXX, 0, float64(SHAPECUBE_X-1)))
			yShapeCube = int(R(-punto.Y, HMIN, HMAX, 0, float64(SHAPECUBE_Y-1)))
			zShapeCube = int(R(-punto.Z, 0, MAXHORIZON, 0, float64(SHAPECUBE_Z-1)))

			snr = shapeCube[xShapeCube][yShapeCube][zShapeCube].R
			sng = shapeCube[xShapeCube][yShapeCube][zShapeCube].G
			snb = shapeCube[xShapeCube][yShapeCube][zShapeCube].B
			sna = shapeCube[xShapeCube][yShapeCube][zShapeCube].A

			snsample = R(float64(snr), float64((float64(sng)*0.625+float64(snb)*0.25+float64(sna)*0.125)-1), 1, 0, 1)

			//sn = SAT (R(snsample * sa, 1 - GC * wmc , 1, 0, 1)) * da

			// Detail
			xDetailCube = int(R(punto.X, MINX, MAXX, 0, float64(DETAILCUBE_X)))
			yDetailCube = int(R(-punto.Y, HMIN, HMAX, 0, float64(DETAILCUBE_Y)))
			zDetailCube = int(R(-punto.Z, 0, MAXHORIZON, 0, float64(DETAILCUBE_Z)))

			dnr = detailCube[xDetailCube][yDetailCube][zDetailCube].R
			dng = detailCube[xDetailCube][yDetailCube][zDetailCube].G
			dnb = detailCube[xDetailCube][yDetailCube][zDetailCube].B

			dnfbm = float64(dnr)*0.625 + float64(dng)*0.25 + float64(dnb)*0.125
			dnmod = float64(math.Exp(float64(GC*0.75))) * Li(dnfbm, 1-dnfbm, SAT(ph*5))

			snnd = SAT(R(snsample*sa, 1-GC*wmc, 1, 0, 1))

			// Final de shape and detail noise***********************************

			// Valor final de la densidad en el punto.
			densidad += (SAT(R(snnd, dnmod, 1, 0, 1)) * da) * 0.01

		}
	}
	return densidad
}

func main() {
	var weatherMap [WEATHER_X][WEATHER_Y]color.RGBA
	var NDC_x float64
	var NDC_y float64
	var PixelScreen_x float64
	var PixelScreen_y float64
	var PixelCamera_x float64
	var PixelCamera_y float64

	var ro Vectores.Vector
	var rd Vectores.Vector
	var nuevo Vectores.Vector

	var densidad float64

	var shapeCube [SHAPECUBE_X][SHAPECUBE_Y][SHAPECUBE_Z]color.RGBA
	var detailCube [DETAILCUBE_X][DETAILCUBE_Y][DETAILCUBE_Z]color.RGBA
	var color color.RGBA

	//generateWeatherMap()
	weatherMap = loadWeatherMap()
	//generateNoiseCubeRGBA()
	//noiseCube = loadCubeNoiseRGBA()

	//noiseCube[0][0][0] = color.RGBA{0,0,0,0}
	// Formatted string, such as "2h3m0.5s" or "4.503μs"
	//fmt.Println(duration)
	//weatherMap[0][0].R = 0

	start := time.Now()

	shapeCube = loadShapeCube()
	detailCube = loadDetailCube()
	//checkNoiseCube(noiseCube)
	duration := time.Since(start)
	fmt.Println(duration)

	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
	out, err := os.Create("clouds.jpg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Calculo el Field of View. El ángulo es de 45 grados.
	//
	var FOV float64 = math.Tan(float64(ALPHA / 2.0 * math.Pi / 180.0))

	for x := 0; x < WIDTH; x++ {
		for y := 0; y < HEIGHT; y++ {
			// Hacemos las conversiones de espacios
			//
			NDC_x = (float64(x) + correccion) / float64(WIDTH)
			NDC_y = (float64(y) + correccion) / float64(HEIGHT)

			PixelScreen_x = 2*NDC_x - 1
			PixelScreen_y = 2*NDC_y - 1

			PixelCamera_x = PixelScreen_x * ImageAspectRatio * FOV
			PixelCamera_y = PixelScreen_y * FOV

			// Origen y dirección

			ro = EYE
			nuevo.X = PixelCamera_x
			nuevo.Y = PixelCamera_y
			nuevo.Z = -1

			densidad = 0
			rd = nuevo.Sub(ro).Normalize()

			densidad = calculaDensidad(ro, rd, weatherMap, shapeCube, detailCube)
			color.R = uint8(densidad * 255)
			color.G = uint8(densidad * 255)
			color.B = uint8(densidad * 255)
			color.A = 255

			img.Set(x, y, color)
		}
	}
	weatherMap[0][0].R = 0
	var opt jpeg.Options

	opt.Quality = 80
	// ok, write out the data into the new JPEG file

	err = jpeg.Encode(out, img, &opt) // put quality to 80%
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Generated image to %s \n", "clouds.jpg")

}
