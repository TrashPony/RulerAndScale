package TransportData

import (
	"github.com/jacobsa/go-serial/serial"
	"io"
)

type Port struct {
	Name       string
	Config     *serial.OpenOptions
	Connection io.ReadWriteCloser
}

func (p *Port) SendBytes(command []byte, countRead int) (data []byte) {

	//ioutil.ReadAll(p.Connection)

	_, err := p.Connection.Write(command)
	if err != nil {
		println("ошибка записи" + err.Error())
		return data
	}

	//time.Sleep(150 * time.Millisecond)
	data = make([]byte, countRead)

	n, err := p.Connection.Read(data)
	if err != nil {
		println("ошибка чтения: " + err.Error())
	}

	if n == countRead {
		return data
	} else {
		println("ошибка чтения: неправильное количество байт")
	}

	return
}
