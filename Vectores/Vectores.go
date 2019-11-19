package Vectores

import "math"

// Vector - struct holding X Y Z values of a 3D vector
type Vector struct {
	X, Y, Z float32
}

func (a Vector) Add(b Vector) Vector {
	a.X += b.X
	a.Y += b.Y
	a.Z += b.Z
	return a
}

func (a Vector) Sub(b Vector) Vector {
	return Vector{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

func (a Vector) MultiplyByScalar(s float32) Vector {
	return Vector{
		X: a.X * s,
		Y: a.Y * s,
		Z: a.Z * s,
	}
}

func (a Vector) Dot(b Vector) float32 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vector) Length() float32 {
	return float32(math.Sqrt(float64(a.Dot(a))))
}

func (a Vector) Normalize() Vector {
	return a.MultiplyByScalar(1. / a.Length())
}

// Valor absoluto de las coordenadas de un vector.
//
func (a Vector) Abs() Vector {
	return Vector{
		X: float32(math.Abs(float64(a.X))),
		Y: float32(math.Abs(float64(a.Y))),
		Z: float32(math.Abs(float64(a.Z))),
	}
}
