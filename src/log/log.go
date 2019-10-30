package log

import (
	"os"
	"strconv"
	"time"
)

func Write(weightBox, widthBox, heightBox, lengthBox int) {
	currentTime := time.Now().Local()

	f, err := os.OpenFile("./"+currentTime.Format("2006-01-02"), os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		f, _ = os.Create("./" + currentTime.Format("2006-01-02"))
	}

	f.WriteString(currentTime.Format("15:04:05") + " - " + " weightBox: " + strconv.Itoa(weightBox) +
		" widthBox: " + strconv.Itoa(widthBox) +
		" heightBox: " + strconv.Itoa(heightBox) +
		" lengthBox: " + strconv.Itoa(lengthBox) + "\n")

	defer f.Close()
}
