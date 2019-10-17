package main

import (
	"github.com/TrashPony/RulerAndScale/InputData"
	"github.com/TrashPony/RulerAndScale/Log"
	"github.com/TrashPony/RulerAndScale/ParseData"
	"github.com/TrashPony/RulerAndScale/TransportData"
	"strconv"
	"time"
)

var scalePort, rulerPort *TransportData.Port

func main() {
	go TransportData.SelectPort()
	Controller()
}

func Controller() {

	for {

		correctWeight := -1
		widthBox, heightBox, lengthBox := -1, -1, -1
		onlyWeight := false

		scalePort := TransportData.Ports.GetPort("scale")
		if scalePort != nil {
			scaleResponse := TransportData.SendScaleCommand(scalePort)
			if scaleResponse == nil {

				println("Весы отвалились")
				TransportData.Ports.ResetPort("scale")

			} else {
				correctWeight = int(ParseData.ParseScaleData(scaleResponse))
			}
		}

		rulerPort := TransportData.Ports.GetPort("ruler")
		if rulerPort != nil {
			rulerResponse := TransportData.SendRulerCommand(rulerPort)
			if rulerResponse == nil {

				println("Линейка отвалилась")
				TransportData.Ports.ResetPort("ruler")

			} else {
				widthBox, heightBox, lengthBox, onlyWeight = ParseData.ParseRulerData(rulerResponse)
			}
		}

		// значения не могут быть отрицаельными если это так то это ошибка
		if correctWeight < 0 || widthBox < 0 || heightBox < 0 || lengthBox < 0 {

			if widthBox > 0 || heightBox > 0 || lengthBox > 0 {
				println(widthBox, heightBox, lengthBox)
			}

			continue
		}

		if scalePort != nil && rulerPort != nil {

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
