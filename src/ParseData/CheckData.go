package ParseData

var OldWeightValue = 0

func CheckData(weightBox, widthBox, heightBox, lengthBox int) (checkData bool, led bool) {
	// значение нуля попадает сюда когда весы не готовы, если весы не готовы значит ждем калибровки с последующим автозабитием
	if weightBox <= 0 {
		OldWeightValue = 0
	}

	// значение весов регулирует авто забитие, оно происходит только при изменение весе, если вес не изменялся то автозабитие не происходит

	if weightBox > 0 && widthBox > 0 && heightBox > 0 && lengthBox > 0 {
		if (OldWeightValue - 4) <= weightBox && weightBox <= (OldWeightValue + 4) {
			return false, true
		} else {
			OldWeightValue = weightBox
			return true, true
		}
	} else {
		return false, false
	}
}
