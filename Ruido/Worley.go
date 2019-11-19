package Ruido

import (
	"../Vectores"
	"math"
	"math/rand"
)

// TODO : Implementar Worley2D
// TODO : Es leeeeeeennnnnntooooooo

// Distribución de Poisson
// No estaría mal encontrar la forma de generarla en automático.
//
//var distP [10]int = [10]int{4, 4, 6, 5, 3, 4, 8, 8, 7, 5}
//var distP [10]int = [10]int{2, 2, 2, 2, 2, 2, 2, 2, 2, 2}
var distP [10]int = [10]int{1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

// Calculo una semilla diferente para cada uno de los cubos.
//
func calculateSeed(cube Vectores.Vector) int {
	var seed int

	seed = int(541*cube.X+79*cube.Y+31*cube.Z) % 4294967296
	return seed
}

// Devolvemos el número de puntos por cada cubo.
//
func pointNumber(seed int) int {
	var max int = 9
	var min int = 0
	var index int

	rand.Seed(int64(seed))

	index = rand.Intn(max-min) + min
	return distP[index]
}

// uniform returns a uniformly random float in [0,1).
// https://stackoverflow.com/users/5181219/ted
func uniform() float32 {
	var sig uint64

	sig = rand.Uint64() % (1 << 52)
	//return (1 + float32(sig)/(1<<52)) / math.Pow(2, geometric())
	return float32(float64(1+float32(sig)/(1<<52)) / math.Pow(2, float64(geometric())))
}

// geometric returns a number picked from a geometric
// distribution of parameter 0.5.
// https://stackoverflow.com/users/5181219/ted
func geometric() float32 {
	var b float32 = 1
	for rand.Uint64()%2 == 0 {
		b++
	}
	return b
}

// Genero un punto 3d con coordenadas en el rango 0-1
// Les añado las del cubo para poder calcular la distancia al punto
// original, ya que las coordenadas de éste no están en el rango 0-1
//
func generatePoint(seed int, cube Vectores.Vector) Vectores.Vector {

	var point Vectores.Vector
	var uniforme float32

	uniforme = uniform()
	point.X = uniforme + cube.X
	point.Y = uniforme + cube.Y
	point.Z = uniforme + cube.Z

	return point
}

// Distancia euclídea entre dos puntos.
// Quizá exista dentro de alguna librería
//
func euclidean(punto Vectores.Vector, coord Vectores.Vector) float32 {
	var distancia float32

	distancia = float32(math.Sqrt(float64(math.Pow(float64(coord.X-punto.X), 2)))) +
		float32(math.Sqrt(float64(math.Pow(float64(coord.Y-punto.Y), 2)))) +
		float32(math.Sqrt(float64(math.Pow(float64(coord.Z-punto.Z), 2))))
	return distancia

}

// Función clip. Devuelve valor siempre que esté entre min y max; de lo contrario
// retorna uno de éstos valores.
//
func Clip(valor float32, min float32, max float32) float32 {
	if valor < min {
		return min
	}
	if valor > max {
		return max
	}
	return valor
}

// Worley 3d con cubos.
// Referencias :
//   https://thebookofshaders.com/12/
//   http://www.rhythmiccanvas.com/research/papers/worley.pdf
//   https://github.com/bhickey/worley/blob/master/worley.c
//   https://www.kdnuggets.com/2017/08/comparing-distance-measurements-python-scipy.html
//
// Es la versión de varios puntos por cubo. Con solo uno funciona igual de bien.
// El problema está en la lentitud de rand.seed()
//
/*
func Worley3D(punto Vectores.Vector) float32 {
	var minimo float32 = 1000
	var seed int
	var points int

	var cx int
	var cy int
	var cz int

	var cube Vectores.Vector
	var dummy Vectores.Vector

	for cx = int(math.Floor(float64(punto.X - 1))); cx <= int(math.Floor(float64(punto.X+2))); cx++ {
		for cy = int(math.Floor(float64(punto.Y - 1))); cy <= int(math.Floor(float64(punto.Y+2))); cy++ {
			for cz = int(math.Floor(float64(punto.Z - 1))); cz <= int(math.Floor(float64(punto.Z+2))); cz++ {
				cube.X = float32(cx)
				cube.Y = float32(cy)
				cube.Z = float32(cz)

				seed = calculateSeed(cube)
				points = pointNumber(seed)

				distancias := make([]float32, points)

				for cp := 0; cp < points; cp++ {
					dummy = generatePoint(seed, cube)
					distancias[cp] = euclidean(punto, dummy)
				}

				sort.float32s(distancias)

				if distancias[0] < minimo {
					minimo = distancias[0]
				}
			}
		}
	}
	return Clip(minimo, 0, 1)

}
*/

func Worley3D2(punto Vectores.Vector) float32 {
	var minimo float32 = 1000
	var seed int

	var cx int
	var cy int
	var cz int

	var cube Vectores.Vector
	var dummy Vectores.Vector

	var distancia float32

	for cx = int(math.Floor(float64(punto.X - 1))); cx <= int(math.Floor(float64(punto.X+2))); cx++ {
		for cy = int(math.Floor(float64(punto.Y - 1))); cy <= int(math.Floor(float64(punto.Y+2))); cy++ {
			for cz = int(math.Floor(float64(punto.Z - 1))); cz <= int(math.Floor(float64(punto.Z+2))); cz++ {
				cube.X = float32(cx)
				cube.Y = float32(cy)
				cube.Z = float32(cz)

				seed = calculateSeed(cube)
				//points = pointNumber(seed)

				rand.Seed(int64(seed))
				dummy = generatePoint(seed, cube)
				distancia = euclidean(punto, dummy)

				if distancia < minimo {
					minimo = distancia
				}
			}
		}
	}
	return Clip(minimo, 0, 1)

}
