package main

import (
	. "./model"
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/websocket"
	stdlog "log"
	"net/http"
	"os"
	"strconv"
)

// standard types
//https://github.com/gorilla/websocket/blob/master/conn.go

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	formatter := log.TextFormatter{}
	formatter.FullTimestamp = true
	formatter.ForceColors = true

	log.SetFormatter(&formatter)

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//log.Out(os.Stdout)
	file, err := os.OpenFile("logs/logrus.log", os.O_RDWR|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		file, err := os.OpenFile("logs/logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			log.SetOutput(file)
		} else {
			log.Warn("Failed to log to file, using default stderr")
		}
	}

	stdlog.SetOutput(file)

	//saving server pid
	pid, err := os.OpenFile("bin/pid", os.O_CREATE|os.O_WRONLY, 0666)
	pid.WriteString(strconv.Itoa(os.Getpid()))
	pid.Close()

	// Only log the warning severity or above.
	//log.SetLevel(log.WarnLevel)
}

func main() {
	log.Info("Launching server on port 8081 with PID ", strconv.Itoa(os.Getpid()), "...")
	go GetServer().Start()

	var addr = flag.String("addr", ":8081", "http service address")
	http.HandleFunc("/", handleNewConnection)
	http.ListenAndServe(*addr, nil)
}

func handleNewConnection(w http.ResponseWriter, r *http.Request) {
	log.Info("New Connection Established:")
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn(err)
		return
	}
	var client = GetServer().GetClientFactory().NewClient(conn)
	go client.Start()
	go client.StartWriter()
}
