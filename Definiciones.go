package main

import (
	"./Vectores"
)

const CUBENOISE_WIDTH int = 50
const CUBENOISE_HEIGHT int = 50
const CUBENOISE_DEPTH int = 50

// Multiplicador para ruído de poca frecuencia.
const LOW_COVERAGE_NOISE float64 = 0.005

const HIGH_COVERAGE_NOISE float64 = 0.05

const LOW_FREQ_NOISE = 0.005
const SIMPLEX_FREQ = 0.025
const MEDIUM_FREQ_NOISE = 0.025
const HIGH_FREQ_NOISE = 0.1
const HIGHEST_FREQ_NOISE = 0.25

const SHAPECUBE_X = 20
const SHAPECUBE_Y = 20
const SHAPECUBE_Z = 20

const DETAILCUBE_X = 32
const DETAILCUBE_Y = 32
const DETAILCUBE_Z = 32

const WEATHER_X = 512
const WEATHER_Y = 512

var WIDTH int = 640
var HEIGHT int = 480

// Ángulo para el FOV. Actúa como una especie de zoom.
var ALPHA float64 = 55.0
var ImageAspectRatio float64 = float64(WIDTH) / float64(HEIGHT)
var correccion float64 = 0.5

var EYE = Vectores.Vector{0, 0, 0}

const GC float64 = 0.5
const GD float64 = 0.2

// Alturas entre las que se encuentra la capa de nubes.
const HMIN = 20
const HMAX = 50
const HINTERVAL = HMAX - HMIN

const MAXHORIZON = 250
const MAXX = 150
const MINX = -150
