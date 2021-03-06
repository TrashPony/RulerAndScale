package main

import (
	"./TransportData"
	"./ParseData"
	"./InputData"
	"./Log"
	"strconv"
	"time"
)

var scalePort, rulerPort *TransportData.Port

func main() {
	Controller()
}

func Controller() {
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
				widthBox, heightBox, lengthBox, onlyWeight := ParseData.ParseRulerData(rulerResponse)

				correctWeight := int(weightBox)

				checkScaleData, led := ParseData.CheckData(correctWeight, widthBox, heightBox, lengthBox, onlyWeight)

				if led {
					rulerPort.Connection.Write([]byte{0x66}) // байт готовности, включает диод
				} else {
					rulerPort.Connection.Write([]byte{0x55}) // байт готовности, выключает диод
				}

				if checkScaleData {

					if onlyWeight {

						InputData.ToClipBoard(strconv.Itoa(correctWeight))

					} else {
						InputData.ToClipBoard(":" + strconv.Itoa(correctWeight) +
							":" + strconv.Itoa(widthBox) +
							":" + strconv.Itoa(heightBox) +
							":" + strconv.Itoa(lengthBox))
					}

					InputData.ToClipBoard("_ESC_Save")

					Log.Write(correctWeight, widthBox, heightBox, lengthBox)

					time.Sleep(time.Second * 3)
				}
			}
		}
	}
}
