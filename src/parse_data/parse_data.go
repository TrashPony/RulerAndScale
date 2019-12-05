package parse_data

func ParseScaleData(data []byte) (weightBox float64) {
	/*
		   80 00
		   EC 00

		80             - готовность 128 - готов, 0 - не готов
		00 (1я строка) - 00-1г, 01-0.1г, 04-0.01г, 05-0.1кг
		EC             - вес в определившейся дискретности
		00 (2я строка) - 0 это (+) 80 это (-) отрицательный , положительный вес.
	*/

	if data[0] == 128 && (data[2] != 0 || data[3] != 0) {

		// data.ReadyAndDiscreteness[0] - готовность
		// data.ReadyAndDiscreteness[1] - дискретность

		if data[1] == 0 {
			if data[3] == 0 { // вес уместился в 1н байт

				weightBox = float64(data[2])

				if weightBox <= 256 {
					return
				} else {
					return 0
				}
			}

			if data[3] != 0 { // не уместился

				weightBox = (256 * float64(data[3])) + float64(data[2])

				if weightBox < 15000 {
					return
				} else {
					return 0
				}
			}
		}

		if data[1] == 4 {
			if data[3] == 0 { // вес уместился в 1н байт

				weightBox = float64(data[2]) * 10

				if weightBox <= 2560 {
					return
				} else {
					return 0
				}
			}

			if data[3] != 0 { // не уместился

				weightBox = ((256 * float64(data[3])) + float64(data[2])) * 10

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

func ParseRulerData(data []byte) (widthBox, heightBox, lengthBox int, weight bool) {

	/*

			команда 0x95, ответ 0x7F - успешное подключение к устройству

			команда 0x88 - запрос габаритов

			протокол обмена состоит из 3х строк такого типа "0x2D 0x0B 0x64 0x7B"

		    они состоят из "0x2D - начало строки, 0xB - датчик, 0x64 - растояние, 0x7B - конец строки" все в 16ричной системе счисления

			0x0B - ширина
			0x16 - высота
			0x21 - длина

	*/

	widthBox = rulerParse([]byte{data[0], data[1], data[2], data[3]}, 0x0B)    //ширина
	heightBox = rulerParse([]byte{data[4], data[5], data[6], data[7]}, 0x16)   //высота
	lengthBox = rulerParse([]byte{data[8], data[9], data[10], data[11]}, 0x21) //длина

	if data[12] == 0 {
		weight = false
	} else {
		weight = true
	}

	if widthBox > 200 {
		widthBox = -1
	}

	if heightBox > 200 {
		heightBox = -1
	}

	if lengthBox > 200 {
		lengthBox = -1
	}

	return
}

func ParseRulerIndicationData(data []byte) (left, right, top, back, wMax, tMax, lMax, widthBox, heightBox, lengthBox int, weight bool) {
	// TODO это была не лучшая моя идея :D

	left = rulerParse([]byte{data[0], data[1], data[2], data[3]}, 0x0B)
	right = rulerParse([]byte{data[4], data[5], data[6], data[7]}, 0xBB)
	top = rulerParse([]byte{data[8], data[9], data[10], data[11]}, 0x16)
	back = rulerParse([]byte{data[12], data[13], data[14], data[15]}, 0x21)

	wMax = rulerParse([]byte{data[16], data[17], data[18], data[19]}, 0x0B) //ширина max
	tMax = rulerParse([]byte{data[20], data[21], data[22], data[23]}, 0x16) //высота max
	lMax = rulerParse([]byte{data[24], data[25], data[26], data[27]}, 0x21) //длина max

	widthBox = rulerParse([]byte{data[28], data[29], data[30], data[31]}, 0x0B)  //ширина
	heightBox = rulerParse([]byte{data[32], data[33], data[34], data[35]}, 0x16) //высота
	lengthBox = rulerParse([]byte{data[36], data[37], data[38], data[39]}, 0x21) //длина

	if data[40] == 0 {
		weight = false
	} else {
		weight = true
	}

	if widthBox > 202 {
		widthBox = -1
	}

	if heightBox > 202 {
		heightBox = -1
	}

	if lengthBox > 202 {
		lengthBox = -1
	}

	if left > 202 {
		left = -1
	}

	if right > 202 {
		right = -1
	}

	if top > 202 {
		top = -1
	}

	if back > 202 {
		back = -1
	}

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