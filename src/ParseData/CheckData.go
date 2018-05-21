package ParseData

var OldWeightValue = 0
var faultWeight = 40

func CheckData(weightBox, widthBox, heightBox, lengthBox int) (checkScaleData bool, checkRulerData bool, led bool) {
	// значение нуля попадает сюда когда весы не готовы, если весы не готовы значит ждем калибровки с последующим автозабитием
	if weightBox <= 0 {
		OldWeightValue = 0
	}

	// значение весов регулирует авто забитие, оно происходит только при изменение весе, если вес не изменялся то автозабитие не происходит

	if weightBox > 0 && widthBox > 0 && heightBox > 0 && lengthBox > 0 {
		if (OldWeightValue - faultWeight) <= weightBox && weightBox <= (OldWeightValue + faultWeight) {
			return false,false, true
		} else {
			OldWeightValue = weightBox
			return true, true, true
		}
	} else {
		if weightBox > 0 && heightBox > 0 {
			if (OldWeightValue - faultWeight) <= weightBox && weightBox <= (OldWeightValue + faultWeight) {
				return false, false, true
			} else {
				OldWeightValue = weightBox
				return true, false, true
			}
		} else {
			return false, false, false
		}
	}
}
