package InputData

import (
	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	"log"
	"time"
)

func PrintResult(data string) {

	println(data)

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
	robotgo.KeyTap("v", "control")
	robotgo.KeyTap("enter")
}
