package TransportData

func ParseScaleData(data *ScaleResponse) (weightBox float64) {
	/*
	   80 00
	   EC 00

	80             - готовность 128 - готов, 0 - не готов
	00 (1я строка) - 00-1г, 01-0.1г, 04-0.01г, 05-0.1кг
	EC             - вес в определившейся дискретности
	00 (2я строка) - 0 это (+) 80 это (-) отрицательный , положительный вес.
	*/

	if data.ReadyAndDiscreteness[0] == 128 && (data.Weight[0] !=0 || data.Weight[1] !=0) {

		// data.ReadyAndDiscreteness[0] - готовность
		// data.ReadyAndDiscreteness[1] - дискретность

		if data.ReadyAndDiscreteness[1] == 0 {
			if  data.Weight[1] == 0 { // вес уместился в 1н байт
				weightBox = float64(data.Weight[0]) * 1
				return
			}

			if data.Weight[1] != 0 { // не уместился
				weightBox = ((256 * float64(data.Weight[1])) + float64(data.Weight[0])) * 1
				return
			}
		}

		if data.ReadyAndDiscreteness[1] == 4 {
			if  data.Weight[1] == 0 { // вес уместился в 1н байт
				weightBox = float64(data.Weight[0]) * 0.01
				return
			}

			if data.Weight[1] != 0 { // не уместился
				weightBox = ((256 * float64(data.Weight[1])) + float64(data.Weight[0])) * 0.01
				return
			}
		}
	}

	return
}

func ParseRulerData(data *RulerResponse) (widthBox, heightBox, lengthBox int) {

	/*
	команда 0x95, ответ 0x7F - успешное подключение к устройству

	команда 0x99 - запрос высоты
	команда 0x88 - запрос ширины
	команда 0x77 - запрос длины

	протокол измерения "0x2D 0x7F/0x7E 0x0B 0x64 0x7B"
	команда 99, ответ "0x2D - начало строки, 0x7F/0x7E - флаг готовности,  0xB - датчик, 0x64 - растояние, 0x7B - конец строки" все в 16ричной системе счисления

	0x0B - ширина
	0x16 - высота
	0x21 - длина

	0x7F - готов
	0x7E - неготов
	*/

	widthBox = rulerParse(data.Width, 0x0B)   //ширина
	heightBox = rulerParse(data.Height, 0x16) //высота
	lengthBox = rulerParse(data.Length, 0x21) //длина

	return
}

func rulerParse(data []byte, id byte) (result int) {
	if data[0] == 45 {
		if data[1] == 127 && data[2] == id && data[4] == 123 { // ширина
			result = int(data[3])
			return
		} else {
			result = 0
			return
		}
	}
	return
}
