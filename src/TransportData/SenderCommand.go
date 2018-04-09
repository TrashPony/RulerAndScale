package TransportData

type ScaleResponse struct {
	ReadyAndDiscreteness []byte
	Weight               []byte
}

type RulerResponse struct {
	Width  []byte
	Height []byte
	Length []byte
}

func SendScaleCommand(port *Port) (*ScaleResponse) {

	var response ScaleResponse
	countRead := 2

	// не/готовность 0/128 и дискретность 0х00-1г,0х01-0.1г,0х04-0.01кг,0х05-0.1кг
	response.ReadyAndDiscreteness = port.SendBytes([]byte{0x48}, countRead)

	//вес в виде 2х байтов n х n
	response.Weight = port.SendBytes([]byte{0x45}, countRead)

	if response.ReadyAndDiscreteness != nil && response.Weight != nil {
		return &response
	} else {
		return nil
	}
}

func SendRulerCommand(port *Port) (*RulerResponse) {

	var response RulerResponse
	countRead := 5

	// запрос ширины коробки
	response.Width = port.SendBytes([]byte{0x88}, countRead)

	// запрос высоты коробки
	response.Height = port.SendBytes([]byte{0x99}, countRead)

	// запрос длинны коробки
	response.Length = port.SendBytes([]byte{0x77}, countRead)

	if response.Width != nil && response.Height != nil && response.Length != nil {
		return &response
	} else {
		return nil
	}
}