package main

import (
	"./TransportData"
	"github.com/tarm/serial"
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
				TransportData.ParseScaleData(scaleResponse)
				TransportData.ParseRulerData(rulerResponse)
			}
		}
	}
}
