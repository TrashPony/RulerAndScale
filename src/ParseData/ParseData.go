package ParseData

import "github.com/TrashPony/RulerAndScale/TransportData"

func ParseScaleData(data *TransportData.ScaleResponse) (weightBox float64) {
	/*
		   80 00
		   EC 00

		80             - готовность 128 - готов, 0 - не готов
		00 (1я строка) - 00-1г, 01-0.1г, 04-0.01г, 05-0.1кг
		EC             - вес в определившейся дискретности
		00 (2я строка) - 0 это (+) 80 это (-) отрицательный , положительный вес.
	*/

	if data.ReadyAndDiscreteness[0] == 128 && (data.Weight[0] != 0 || data.Weight[1] != 0) {

		// data.ReadyAndDiscreteness[0] - готовность
		// data.ReadyAndDiscreteness[1] - дискретность

		if data.ReadyAndDiscreteness[1] == 0 {
			if data.Weight[1] == 0 { // вес уместился в 1н байт

				weightBox = float64(data.Weight[0])

				if weightBox <= 256 {
					return
				} else {
					return 0
				}
			}

			if data.Weight[1] != 0 { // не уместился

				weightBox = (256 * float64(data.Weight[1])) + float64(data.Weight[0])

				if weightBox < 15000 {
					return
				} else {
					return 0
				}
			}
		}

		if data.ReadyAndDiscreteness[1] == 4 {
			if data.Weight[1] == 0 { // вес уместился в 1н байт

				weightBox = float64(data.Weight[0]) * 10

				if weightBox <= 2560 {
					return
				} else {
					return 0
				}
			}

			if data.Weight[1] != 0 { // не уместился

				weightBox = ((256 * float64(data.Weight[1])) + float64(data.Weight[0])) * 10

				if weightBox < 60000 {
					return
				} else {
					return 0
				}
			}
		}
	}

	return
}

func ParseRulerData(data []byte, command []byte) (widthBox, heightBox, lengthBox int) {

	/*

			команда 0x95, ответ 0x7F - успешное подключение к устройству

			команда 0x88 - запрос габаритов

			протокол обмена состоит из 3х строк такого типа "0x2D 0x0B 0x64 0x7B"

		    они состоят из "0x2D - начало строки, 0xB - датчик, 0x64 - растояние, 0x7B - конец строки" все в 16ричной системе счисления

			0x0B - ширина
			0x16 - высота
			0x21 - длина

	*/

	if command != nil && (command[0] == 0x88 || command[0] == 0x89) {
		widthBox = rulerParse([]byte{data[1], data[2], data[3], data[4]}, 0x0B)     //ширина
		heightBox = rulerParse([]byte{data[5], data[6], data[7], data[8]}, 0x16)    //высота
		lengthBox = rulerParse([]byte{data[9], data[10], data[11], data[12]}, 0x21) //длина
	}

	return
}

func ParseRulerIndicationData(data []byte, command []byte) (left, right, top, back int) {

	left = rulerParse([]byte{data[1], data[2], data[3], data[4]}, 0x0B)
	right = rulerParse([]byte{data[5], data[6], data[7], data[8]}, 0xBB)
	top = rulerParse([]byte{data[9], data[10], data[11], data[12]}, 0x16)
	back = rulerParse([]byte{data[13], data[14], data[15], data[16]}, 0x21)

	return
}

func rulerParse(data []byte, id byte) (result int) {
	if data[0] == 45 {
		if data[1] == id && data[3] == 123 {
			result = int(data[2])
			return
		} else {
			result = -1
			return
		}
	}
	return
}
