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
				widthBox, heightBox, lengthBox := ParseData.ParseRulerData(rulerResponse)

				correctWeight := int(weightBox)

				checkScaleData, _, led := ParseData.CheckData(correctWeight, widthBox, heightBox, lengthBox)
/*
:2450:30:6:22
_ESC_Save
:930:31:4:23
_ESC_Save
:930:30:3:24
_ESC_Save
:930:30:3:23
_ESC_Save
:920:29:2:20
_ESC_Save
:940:31:2:20
_ESC_Save
:210:0:1:20
_ESC_Save
:220:7:1:21
_ESC_Save
:210:7:1:21
_ESC_Save
:210:8:1:21
_ESC_Save
:220:0:1:18
_ESC_Save
:210:0:1:18
_ESC_Save
:210:0:3:19
_ESC_Save
:940:28:3:22
_ESC_Save
:940:30:2:23
_ESC_Save
:560:28:3:26
_ESC_Save
:110:0:3:16
_ESC_Save
:550:31:4:25
_ESC_Save
:220:22:30:31
_ESC_Save
:230:21:30:31
_ESC_Save


*/
				if led {
					rulerPort.Connection.Write([]byte{0x66}) // байт готовности, включает диод
				} else {
					rulerPort.Connection.Write([]byte{0x55}) // байт готовности, выключает диод
				}

				if checkScaleData {

					InputData.ToClipBoard(":" + strconv.Itoa(correctWeight) +
										":" + strconv.Itoa(widthBox) +
										":" + strconv.Itoa(heightBox) +
										":" + strconv.Itoa(lengthBox))

					InputData.ToClipBoard("_ESC_Save")

					Log.Write(correctWeight, widthBox, heightBox, lengthBox)

					time.Sleep(time.Second * 3)
				} else {
					/*if checkScaleData {

						InputData.ToClipBoard(":" + strconv.Itoa(correctWeight))
						InputData.ToClipBoard("_ESC_Save")

						Log.Write(correctWeight, widthBox, heightBox, lengthBox)

						time.Sleep(time.Second * 3)

					}*/
				}
			}
		}
	}
}