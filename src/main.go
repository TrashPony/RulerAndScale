package main

import (
	"github.com/TrashPony/RulerAndScale/InputData"
	"github.com/TrashPony/RulerAndScale/Log"
	"github.com/TrashPony/RulerAndScale/ParseData"
	"github.com/TrashPony/RulerAndScale/TransportData"
	"github.com/TrashPony/RulerAndScale/websocket"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	go TransportData.SelectPort()
	go Controller()

	router := mux.NewRouter()

	router.HandleFunc("/ws", websocket.HandleConnections)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../static/")))

	go websocket.Sender()

	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Panic(err)
	}
}

func Controller() {

	for {

		if len(websocket.UsersWs) > 0 {
			println("происходит дебаг :D")
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		correctWeight := -1
		widthBox, heightBox, lengthBox := -1, -1, -1
		onlyWeight := false

		scalePort := TransportData.Ports.GetPort("scale")
		if scalePort != nil {
			scaleResponse, err := scalePort.SendScaleCommand()
			if scaleResponse == nil && err.Error() != "wrong" {

				println("Весы отвалились")
				TransportData.Ports.ResetPort("scale")

			} else {

				if err != nil && err.Error() == "wrong" {
				} else {
					correctWeight = int(ParseData.ParseScaleData(scaleResponse))
					if correctWeight == 0 { // todo не уверен что это работает как надо :D
						// иногда сериал порт посылает прошлые данные и от них надо избавится или смещает биты
						scalePort.Reconnect(0)
						scalePort.ReadBytes(5)
					}
				}
			}
		}

		rulerPort := TransportData.Ports.GetPort("ruler")
		if rulerPort != nil {
			rulerResponse, err := rulerPort.SendRulerCommand([]byte{0x88}, 14)

			if err != nil && err.Error() != "wrong_data" {

				println("Линейка отвалилась")
				TransportData.Ports.ResetPort("ruler")

			} else {
				if rulerResponse != nil {
					widthBox, heightBox, lengthBox, onlyWeight = ParseData.ParseRulerData(rulerResponse, []byte{0x88})
				}
			}
		}

		if widthBox > 0 || heightBox > 0 || lengthBox > 0 || correctWeight > 0 {
			println(widthBox, heightBox, lengthBox, correctWeight, onlyWeight)
		}

		// значения не могут быть отрицаельными если это так то это ошибка
		if correctWeight < 0 {
			continue
		}

		//if !onlyWeight && correctWeight > 0 && (widthBox < 0 || heightBox < 0 || lengthBox < 0) {
		//	// если на весах что то лежит а дальномеры тупят надо колибровать
		//	rulerPort.SendRulerCommand([]byte{0x93}, 0)
		//	continue
		//}

		if scalePort != nil && rulerPort != nil {

			checkScaleData, led := ParseData.CheckData(correctWeight, widthBox, heightBox, lengthBox, onlyWeight)

			if led {
				rulerPort.SendRulerCommand([]byte{0x66}, 0) // байт готовности, включает диод
			} else {
				rulerPort.SendRulerCommand([]byte{0x55}, 0) // байт готовности, выключает диод
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

				time.Sleep(time.Second * 2)
			}
		}
	}
}
