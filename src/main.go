package main

import (
	"./TransportData"
	"./ParseData"
	"./InputData"
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
				widthBox, heightBox, lengthBox := ParseData.ParseRulerData(rulerResponse)

				checkData, led := ParseData.CheckData(int(weightBox), widthBox, heightBox, lengthBox)
				//6*1523530450259

				if led {
					rulerPort.Connection.Write([]byte{0x66}) // байт готовности, включает диод
				} else {
					rulerPort.Connection.Write([]byte{0x55}) // байт готовности, выключает диод
				}

				if checkData {

					InputData.ToClipBoard(":" + strconv.Itoa(int(weightBox)) +
										":" + strconv.Itoa(widthBox) +
										":" + strconv.Itoa(heightBox) +
										":" + strconv.Itoa(lengthBox))

					InputData.ToClipBoard("_ESC_Save")
					
					time.Sleep(time.Second * 3)
				}
			}
		}
	}
}