package output_data

import (
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	"log"
	"time"
)

func PrintResult(data string) {

	println(data)

	// забитие данных происходит через буфер обмена ОС, и эмуляции нажатия клавиш CTRL+V
	toClipBoard(data)
	toClipBoard("_ESC_Save")
}

func toClipBoard(data string) {
	err := clipboard.WriteAll(data)

	if err != nil {
		log.Fatal(err)
	}

	pressCtrlV()
	time.Sleep(time.Millisecond * 300)
}

func pressCtrlV() {
	// возможно есть более простой путь..
	robotgo.KeyTap("v", "control")
	robotgo.KeyTap("enter")
}
