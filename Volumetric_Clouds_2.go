package main

import "image/color"

func main(){
	var weatherMap [WEATHER_WIDTH][WEATHER_HEIGHT] color.RGBA

	genetareWeatherMap()
	weatherMap = loadWeatherMap()
	weatherMap[0][0].R = 0

}
