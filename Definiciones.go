package main

import (
	"./Vectores"
)

const WEATHER_WIDTH int = 512
const WEATHER_HEIGHT int = 512

// Multiplicador para ruído de poca frecuencia.
const LOW_COVERAGE_NOISE float32 = 0.005

const HIGH_COVERAGE_NOISE float32 = 0.05

const NOISECUBE_X int = 50
const NOISECUBE_Y int = 50
const NOISECUBE_Z int = 50

const WORLEY_MEDIUM float32 = 0.005
const WORLEY_HIGH float32 = 0.05
const WORLEY_HIGHEST float32 = 0.5

var WIDTH int = 640
var HEIGHT int = 480

// Ángulo para el FOV. Actúa como una especie de zoom.
var ALPHA float32 = 55.0
var ImageAspectRatio float32 = float32(WIDTH) / float32(HEIGHT)
var correccion float32 = 0.5

var EYE = Vectores.Vector{0, 0, 0}

const GC float32 = 0.5
const GD float32 = 0.2

// Alturas entre las que se encuentra la capa de nubes.
const HMIN = 20
const HMAX = 50
const HINTERVAL = HMAX - HMIN

const MAXHORIZON = 500
const MAXX = 250
const MINX = -250
