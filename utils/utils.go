package utils

import (
	"encoding/json"
	"fmt"
	"bytes"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

func SetWhiteBackground() {
	fmt.Print("\x1B[47m")
}
func SetBlackBackground() {
	fmt.Print("\x1B[40m")
}

func PrintReset() {
	fmt.Print("\x1B[0m")
}

func PrintBlack(str string) {
	fmt.Print("\x1B[30m" + str)
}

func PrintRed(str string) {
	fmt.Print("\x1B[31m" + str)
}

func PrintGreen(str string) {
	fmt.Print("\x1B[32m" + str)
}

func PrintYellow(str string) {
	fmt.Print("\x1B[33m" + str)
}

func PrintBlue(str string) {
	fmt.Print("\x1B[34m" + str)
}

func PrintWhite(str string) {
	fmt.Print("\x1B[37m" + str)
}

func GetJson(t interface{}) string {
	b, err := json.Marshal(t)
	if err != nil {
		fmt.Print("getJson Marshell Error :", err)
	}
	return string(b)
}

func AllowedCoordinates(x int, y int) bool{
	if (x >= 0 && x <20 && y >=0 && y<20){
		return true
	} else {
		return false
	}
}

func ApiRequest(verb string, url string, data []byte) (*http.Response, map[string]interface{}, error) {
	client := &http.Client{}
	HOST := "https://acrobatt.brixbyte.com/"
	USERNAME := "golang"
	PASSWORD := "8b69c71df014c96d08ae23c11a5f63e3e8a38d75a03cf8728e90518a0c7c8be1e203aae99dce16ee929c1efd7c3deb84e708c43676f8e5b3f951ac129d9eae75"

	req, err := http.NewRequest(verb, HOST+url, bytes.NewBuffer(data))
	req.SetBasicAuth(USERNAME, PASSWORD)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	s := string(bodyText)
	log.Info(resp.Status + ": " + s)

	response := make(map[string]interface{})

	err = json.Unmarshal([]byte(s), &response)

	if err != nil {
		log.Error(err)
		return resp, response, err
	} else {
		return resp, response, nil
	}
	
}