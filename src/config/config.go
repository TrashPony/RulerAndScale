package config

import (
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func GetConfig() (int, int, int) {

	_, err := ioutil.ReadDir("../src/config")
	if err != nil {
		os.MkdirAll("../src/config", os.ModePerm)
	}

	_, err = os.OpenFile("../src/config/config", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		os.Create("../src/config/config")
	}

	configFile, err := ioutil.ReadFile("../src/config/config")
	if err != nil {
		log.Fatal(err)
	}

	configLines := strings.Split(string(configFile), "\n")

	configMap := make(map[string]string)

	for i := 0; i < len(configLines); i++ {
		if configLines[i] != "" {
			configLine := strings.Split(configLines[i], ":")
			configMap[configLine[0]] = configLine[1]
		}
	}

	top, err := strconv.Atoi(configMap["top"])
	if err != nil {
		top = 0
	}

	width, err := strconv.Atoi(configMap["width"])
	if err != nil {
		width = 0
	}

	length, err := strconv.Atoi(configMap["length"])
	if err != nil {
		length = 0
	}

	return top, width, length
}

func WriteConfig(top, width, length int) {
	f, err := os.OpenFile("../src/config/config", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		f, _ = os.Create("../src/config/config")
	}

	ioutil.WriteFile("../src/config/config", nil, 0600)

	f.WriteString("top:" + strconv.Itoa(top) + "\n")
	f.WriteString("width:" + strconv.Itoa(width) + "\n")
	f.WriteString("length:" + strconv.Itoa(length) + "\n")

	f.Sync()
}
