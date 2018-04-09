package main

import (
	"./TransportData"
	"./ParseData"
	"strconv"
)

var scalePort, rulerPort *TransportData.Port

func main() {
	Controller()
}

func Controller()  {

	for {

		if scalePort == nil || rulerPort == nil {

			scalePort, rulerPort = TransportData.SelectPort()

		} else {

			scaleResponse := TransportData.SendScaleCommand(scalePort)
			if scaleResponse == nil {
				println("Весы отвалились")
				scalePort = nil
			}

			rulerResponse := TransportData.SendRulerCommand(rulerPort)
			if rulerResponse == nil {
				println("Линейка отвалилась")
				rulerPort = nil
			}

			if scalePort != nil && rulerPort != nil {

				weightBox := ParseData.ParseScaleData(scaleResponse)
				widthBox, heightBox, lengthBox := ParseData.ParseRulerData(rulerResponse)

				check := ParseData.CheckData(int(weightBox), widthBox, heightBox, lengthBox)

				if check {
					println("Вес коробки: " + strconv.Itoa(int(weightBox)))
					println("Ширина коробки: " + strconv.Itoa(widthBox))
					println("Высота коробки: " + strconv.Itoa(heightBox))
					println("Длинна коробки: " + strconv.Itoa(lengthBox))
					println("-------------------")
				} else {
					//println("не откалиброван")
				}
			}
		}
	}
}
