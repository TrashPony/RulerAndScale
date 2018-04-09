package ParseData

var OldWeightValue int = 0

func CheckData(weightBox, widthBox, heightBox, lengthBox int) (check bool) {

	if weightBox <= 0 {
		OldWeightValue = 0
	}

	// значение весов регулирует авто забитие, оно происходит только при изменение весе, если вес не изменялся то автозабитие не происходит

	if weightBox != 0 && widthBox != 0 && heightBox != 0 && lengthBox != 0 {
		if (OldWeightValue - 4) <= weightBox && weightBox <= (OldWeightValue + 4) {
			return false
		} else {
			OldWeightValue = weightBox
			return true
		}
	} else {
		return false
	}
}
