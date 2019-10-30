package main

import (
	"github.com/TrashPony/RulerAndScale/config"
	log2 "github.com/TrashPony/RulerAndScale/log"
	"github.com/TrashPony/RulerAndScale/output_data"
	"github.com/TrashPony/RulerAndScale/parse_data"
	"github.com/TrashPony/RulerAndScale/transport_data"
	"github.com/TrashPony/RulerAndScale/websocket"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	go transport_data.SelectPort()
	go Controller()

	// вебсервер для страницы состояния системы провески
	router := mux.NewRouter()
	router.HandleFunc("/ws", websocket.HandleConnections)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../static/")))

	go websocket.Sender()
	err := http.ListenAndServe(":8080", router)
	if err != nil {
		log.Panic(err)
	}
}

func Controller() {
	clearBuffer := false

	for {

		scalePort := transport_data.Ports.GetPort("scale")
		rulerPort := transport_data.Ports.GetPort("ruler")

		if rulerPort != nil && !rulerPort.Init {
			clearBuffer = true
			SetSettings(rulerPort)
			continue
		}

		// если на странице состояния наъходится кто либо то дабы избежать
		// конкуретного доступа к стройствам ждем пока страницу закроют
		if len(websocket.UsersWs) > 0 {
			clearBuffer = true
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		// без этой задерржки ардуино не будет успевать отвечать
		time.Sleep(250 * time.Millisecond)

		// очищаем буфер от сообщений отладки,
		// т.к. там у сообщения больше байт
		// все не считаные байты сделают сдвиг в будущих сообщения линейки
		if clearBuffer {
			clearBuffer = false
			if scalePort != nil {
				ioutil.ReadAll(scalePort.Connection)
			}
			if rulerPort != nil {
				ioutil.ReadAll(rulerPort.Connection)
			}
		}

		correctWeight := -1
		widthBox, heightBox, lengthBox := -1, -1, -1
		onlyWeight := false

		if scalePort != nil {
			scaleResponse, err := scalePort.SendScaleCommand()
			if scaleResponse == nil && err.Error() != "wrong_data" {

				println("Весы отвалились")
				transport_data.Ports.ResetPort("scale")

			} else {
				if scaleResponse != nil {
					correctWeight = int(parse_data.ParseScaleData(scaleResponse))
				}
			}
		}

		if rulerPort != nil {
			rulerResponse, err := rulerPort.SendRulerCommand([]byte{0x88, 0x88}, 13)

			if err != nil && err.Error() != "wrong_data" {

				println("Линейка отвалилась")
				transport_data.Ports.ResetPort("ruler")

			} else {
				if rulerResponse != nil {
					widthBox, heightBox, lengthBox, onlyWeight = parse_data.ParseRulerData(rulerResponse)
				}
			}
		}

		// прост)
		if widthBox >= 0 || heightBox >= 0 || lengthBox >= 0 || correctWeight >= 0 {
			println(widthBox, heightBox, lengthBox, correctWeight, onlyWeight)
		}

		// весы не подключены, авто забитие не происходит
		if correctWeight < 0 || scalePort == nil {
			continue
		}

		// если нам нужна линейка а она тупит про пропускаем
		if rulerPort != nil && !onlyWeight && correctWeight > 0 && (widthBox < 0 || heightBox < 0 || lengthBox < 0) {
			continue
		}

		// какойто лазер возможно не откалиброван или находится за пределами измерения
		if !onlyWeight && (widthBox == 202 || heightBox == 202 || lengthBox == 202) {
			// отсылаем команду калибровки
			rulerPort.SendRulerCommand([]byte{0x93, 0x93}, 0)
			continue
		}

		// проверяем надо ли забивать данные или нет
		checkScaleData, _ := parse_data.CheckData(correctWeight, widthBox, heightBox, lengthBox, onlyWeight)

		if checkScaleData {

			// включает диод на дуине
			rulerPort.SendRulerCommand([]byte{0x66, 0x66}, 0)

			// записываем данные в место курсора)
			if onlyWeight {
				output_data.PrintResult(strconv.Itoa(correctWeight))
			} else {
				output_data.PrintResult(":" + strconv.Itoa(correctWeight) +
					":" + strconv.Itoa(widthBox) +
					":" + strconv.Itoa(heightBox) +
					":" + strconv.Itoa(lengthBox))
			}

			log2.Write(correctWeight, widthBox, heightBox, lengthBox)

			// время ожидания после успешного взвешивания
			time.Sleep(time.Second * 2)

			// выключаем диод на дуине
			rulerPort.SendRulerCommand([]byte{0x55, 0x55}, 0)
		}
	}
}

func SetSettings(rulerPort *transport_data.Port) {
	// выставляем настройки линейки
	top, width, length := config.GetConfig()

	rulerPort.SendRulerCommand([]byte{0x90, byte(top)}, 0)
	time.Sleep(300 * time.Millisecond)
	rulerPort.SendRulerCommand([]byte{0x91, byte(width)}, 0)
	time.Sleep(300 * time.Millisecond)
	rulerPort.SendRulerCommand([]byte{0x92, byte(length)}, 0)
	time.Sleep(300 * time.Millisecond)

	rulerResponse, err := rulerPort.SendRulerCommand([]byte{0x89, 0x89}, 41)

	if err != nil && err.Error() != "wrong_data" {
		println("Линейка отвалилась")
		transport_data.Ports.ResetPort("ruler")
		return
	} else {
		if rulerResponse != nil {
			_, _, _, _, widthMax, heightMax, lengthMax, _, _, _, _ := parse_data.ParseRulerIndicationData(rulerResponse)
			// обязательно проверяем установились значения или нет
			if heightMax == top && widthMax == width && lengthMax == length {
				println("линейка променила настройки.")
				rulerPort.Init = true
			}
		}
	}
}
