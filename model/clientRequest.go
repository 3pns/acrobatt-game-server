package model

import (
	"encoding/json"
	_"flag"
	"fmt"
	_"github.com/gorilla/websocket"
	_ "io"
	_"log"
	_"net/http"
	_ "strings"
)

type ClientRequest struct {
	RequestType string `json:"requestType"`
	DataType    string `json:"dataType"`
	Data        []byte `json:"data"`
}

func (request *ClientRequest) MarshalData(t interface{}) {
	fmt.Print("Marshalling")
	value, ok := t.(Board)
	if ok {
		b, err := json.Marshal(value)
		if err != nil {
			fmt.Println(err)
		}
		request.DataType = "Board"
		request.Data = b
	}

}

func (request *ClientRequest) Marshal() []byte{
	marshaleldrequest, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}
	return marshaleldrequest
}


func (request *ClientRequest) Unmarshal() {

}
