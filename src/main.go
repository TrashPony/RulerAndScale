package main

import (
	"./TransportData"
	"github.com/tarm/serial"
	"strconv"
)

func main() {
	for {
		var scalePort, rulerPort *serial.Port

		if scalePort == nil || rulerPort == nil {
			scalePort, rulerPort = TransportData.SelectPort()
		} else {
			scaleResponse := TransportData.SendScaleCommand(scalePort)
			if scaleResponse == nil {
				scalePort = nil
			}

			rulerResponse := TransportData.SendRulerCommand(rulerPort)
			if rulerResponse == nil {
				rulerPort = nil
			}

			if scalePort != nil && rulerPort != nil {
				weightBox := TransportData.ParseScaleData(scaleResponse)
				widthBox, heightBox, lengthBox := TransportData.ParseRulerData(rulerResponse)

				println("Вес коробки: " + strconv.FormatFloat(weightBox, 'E', -1, 64))
				println("Ширина коробки: " + strconv.Itoa(widthBox))
				println("Высота коробки: " + strconv.Itoa(heightBox))
				println("Длинна коробки: " + strconv.Itoa(lengthBox))
			}
		}
	}
}
