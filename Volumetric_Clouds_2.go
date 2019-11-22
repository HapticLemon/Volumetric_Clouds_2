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
	"time"
)

// Calcula el valor del weathermap con los canales R y G
//
func calculaWMC(wc0 float64, wc1 float64) float64 {
	return math.Max(wc0, SAT(GC-0.5)*wc1*2)
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
	var snr float64
	var sng float64
	var snb float64
	var sna float64
	// Detail
	var dnfbm float64
	var dnmod float64
	var snnd float64
	var dnr float64
	var dng float64
	var dnb float64

	var t float64 = 10
	var densidad float64 = 0

	var wc0 float64
	var wc1 float64
	var wch float64
	var wcd float64

	var currColor color.RGBA

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
			xremap = int(math.Abs(R(punto.X, MINX, MAXX, 0, float64(WEATHER_X))))
			zremap = int(math.Abs(R(punto.Z, 0, MAXHORIZON, 0, float64(WEATHER_Y))))

			//var componenteG = weatherMap[xremap][zremap].R
			//return float64(weatherMap[xremap][zremap].R)
			//return float64(componenteG)

			currColor = weatherMap[xremap][zremap]

			// Según el documento, tienen que estar en [0-1]
			wc0 = float64(currColor.R) / 255
			wc1 = float64(currColor.G) / 255
			wch = float64(currColor.B) / 255
			wcd = float64(currColor.A) / 255

			wmc = calculaWMC(wc0, wc1)

			// Si no tenemos probabilidad en un punto, pasaremos a calcular el siguiente.
			if wmc == 0 {
				continue
			}
			// Ojo, las alturas van en negativo.
			ph = calculaPH(-punto.Y)

			// Inicio de Shape-altering height-function--------------------------
			srb = SAT(R(ph, 0, 0.07, 0, 1))
			srt = R(ph, wch*0.2, wch, 1, 0)
			sa = srb * srt
			// Fin de Shape-altering height-function-----------------------------

			// Inicio de Density-altering height-function +++++++++++++++++++++++
			drb = ph * SAT(R(ph, 0, 0.15, 0, 1))
			drt = SAT(R(ph, 0.9, 1.0, 1, 0))
			da = GD * drb * drt * wcd * 2
			// Fin de Density-altering height-function ++++++++++++++++++++++++++

			// Shape and detail noise********************************************
			// Shape
			xShapeCube = int(R(punto.X, MINX, MAXX, 0, float64(SHAPECUBE_X-1)))
			yShapeCube = int(R(-punto.Y, HMIN, HMAX, 0, float64(SHAPECUBE_Y-1)))
			zShapeCube = int(R(-punto.Z, 0, MAXHORIZON, 0, float64(SHAPECUBE_Z-1)))

			currColor = shapeCube[xShapeCube][yShapeCube][zShapeCube]
			snr = float64(currColor.R) / 255
			sng = float64(currColor.G) / 255
			snb = float64(currColor.B) / 255
			sna = float64(currColor.A) / 255

			snsample = R(snr, (sng*0.625+snb*0.25+sna*0.125)-1, 1, 0, 1)

			//sn := SAT (R(snsample * sa, 1 - GC * wmc , 1, 0, 1)) * da
			//sn = 0
			//return sn

			// Detail
			xDetailCube = int(R(punto.X, MINX, MAXX, 0, float64(DETAILCUBE_X)))
			yDetailCube = int(R(-punto.Y, HMIN, HMAX, 0, float64(DETAILCUBE_Y)))
			zDetailCube = int(R(-punto.Z, 0, MAXHORIZON, 0, float64(DETAILCUBE_Z)))

			currColor = detailCube[xDetailCube][yDetailCube][zDetailCube]

			dnr = float64(currColor.R) / 255
			dng = float64(currColor.G) / 255
			dnb = float64(currColor.B) / 255

			dnfbm = dnr*0.625 + dng*0.25 + dnb*0.125
			dnmod = math.Exp(GC*0.75) * Li(dnfbm, 1-dnfbm, SAT(ph*5))

			snnd = SAT(R(snsample*sa, 1-GC*wmc, 1, 0, 1))

			// Final de shape and detail noise***********************************

			// Valor final de la densidad en el punto.
			densidad += (SAT(R(snnd, dnmod, 1, 0, 1)) * da)

		}
	}
	return densidad
}

func calculaDensidadTest(ro Vectores.Vector, rd Vectores.Vector, weatherMap [WEATHER_X][WEATHER_Y]color.RGBA) float64 {
	var punto Vectores.Vector
	var xremap int = 0
	var zremap int = 0

	// Nos servirán de coordenadas para el cubeNoise
	/*
		var xShapeCube int = 0
		var yShapeCube int = 0
		var zShapeCube int = 0
	*/
	// Nos servirán de coordenadas para el cubeNoise
	/*
		var xDetailCube int = 0
		var yDetailCube int = 0
		var zDetailCube int = 0
	*/
	//var wmc float64 = 0.0
	//var ph float64 = 0.0

	// Shape Round top y bottom
	/*
		var srb float64
		var srt float64
		var sa float64
	*/
	// Density Round top y bottom
	/*
		var drb float64
		var drt float64
		var da float64
	*/
	// Shape and detail noise.
	// Shape
	//	var snsample float64
	//var sn float32
	/*
		var snr float64
		var sng float64
		var snb float64
		var sna float64
	*/

	// Detail
	/*
		var dnfbm float64
		var dnmod float64
		var snnd float64
		var dnr float64
		var dng float64
		var dnb float64
	*/
	var t float64 = 10
	var densidad float64 = 0

	var wc0 float64
	//var wc1 float64
	var wch float64
	var wcd float64

	var hCloud float64
	var currColor color.RGBA
	var noiseValue float64

	// Generamos puntos en la trayectoria del rayo hasta que alguno
	// tiene altura como para estar en la capa de nubes.
	//
	p := Ruido.NewPerlin(alpha, beta, n, seed)
	for {
		punto = ro.Add(rd.MultiplyByScalar(t))
		// Si nos salimos de los límites pasamos al siguiente punto.
		// También si el rayo no apunta hacia el cielo, ya que nunca tocará la capa de nubes.
		//
		if punto.Z <= -MAXHORIZON || punto.X > MAXX || punto.X < MINX || punto.Y >= 0 || densidad >= 1 {
			break
		}
		t += 0.25
		// Limitamos a la banda de alturas en la que hay nubes.
		if punto.Y <= -HMIN && punto.Y >= -HMAX {
			// Coordenadas remapeadas para acceder al weeathermap.
			// Creo que no tendría que haber problema con los valores - de X/Z ya que siempre se mapean a +
			//
			xremap = int(math.Abs(R(punto.X, MINX, MAXX, 0, float64(WEATHER_X))))
			zremap = int(math.Abs(R(punto.Z, 0, MAXHORIZON, 0, float64(WEATHER_Y))))

			//var componenteG = weatherMap[xremap][zremap].R
			//return float64(weatherMap[xremap][zremap].R)
			//return float64(componenteG)

			currColor = weatherMap[xremap][zremap]

			// Según el documento, tienen que estar en [0-1]
			wc0 = float64(currColor.R) / 255
			//wc1 = float64(currColor.G) / 255
			wch = float64(currColor.B) / 255
			wcd = float64(currColor.A) / 255

			//return wc0
			// Si el punto no cae en el canal R que indica las nubes, salimos.
			if wc0 == 0 {
				continue
			}

			//wmc = calculaWMC(wc0, wc1)

			// Si no tenemos probabilidad en un punto, pasaremos a calcular el siguiente.
			/*
				if wmc == 0 {
					continue
				}
			*/
			// Compruebo si estoy "dentro" de la nube.
			// La altura del punto está entre HMIN y ALTURANUBE
			hCloud = wch*HINTERVAL + HMIN
			if -punto.Y > hCloud {
				break
			}

			// Shape and detail noise********************************************
			// Shape
			/*
				xShapeCube = int(R(punto.X, MINX, MAXX, 0, float64(SHAPECUBE_X-1)))
				yShapeCube = int(R(-punto.Y, HMIN, HMAX, 0, float64(SHAPECUBE_Y-1)))
				zShapeCube = int(R(-punto.Z, 0, MAXHORIZON, 0, float64(SHAPECUBE_Z-1)))

				currColor = shapeCube[xShapeCube][yShapeCube][zShapeCube]
				snr = float64(currColor.R) / 255
				sng = float64(currColor.G) / 255
				snb = float64(currColor.B) / 255
				sna = float64(currColor.A) / 255

				//snsample = R(snr, (sng*0.625+snb*0.25+sna*0.125)-1, 1, 0, 1)
				//snsample += 12
				snr *= 2
				sng *= 2
				snb *= 2
				sna *= 2
			*/

			// Si estamos dentro, pasamos a calcular la densidad del punto.
			noiseValue = math.Abs(p.Noise3D(float64(xremap)*HIGH_FREQ_NOISE, float64(zremap)*HIGH_FREQ_NOISE, hCloud*HIGH_FREQ_NOISE))
			//noiseValue *= 2
			wcd = 0.3 * wcd * noiseValue
			densidad += wcd

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

	//var shapeCube [SHAPECUBE_X][SHAPECUBE_Y][SHAPECUBE_Z]color.RGBA
	//var detailCube [DETAILCUBE_X][DETAILCUBE_Y][DETAILCUBE_Z]color.RGBA
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

	//shapeCube = loadShapeCube()
	//detailCube = loadDetailCube()
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

			densidad = SAT(calculaDensidadTest(ro, rd, weatherMap))

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
