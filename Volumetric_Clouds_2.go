package main

import (
	"fmt"
	"image/color"
	"time"
)

func main() {
	var weatherMap [WEATHER_WIDTH][WEATHER_HEIGHT]color.RGBA
	//var noiseCube [NOISECUBE_X][NOISECUBE_Y][NOISECUBE_Z] byte

	start := time.Now()
	//genetareWeatherMap()
	//weatherMap = loadWeatherMap()
	generateNoiseCube()
	noiseCube := loadCubeNoise()
	duration := time.Since(start)
	noiseCube[0][0][0] = 0
	// Formatted string, such as "2h3m0.5s" or "4.503μs"
	fmt.Println(duration)
	weatherMap[0][0].R = 0

}
