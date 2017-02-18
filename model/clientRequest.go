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

type Request struct {
	Type string `json:"type"`
	DataType    string `json:"dataType"`
	Data        []byte `json:"data"`
}

func (request *Request) MarshalData(t interface{}) {
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

func (request *Request) Marshal() []byte{
	marshaleldrequest, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}
	return marshaleldrequest
}


func (request *Request) Unmarshal() {
	fmt.Print("Unmarshalling")
}
