package main

import (
  _ "net"
  "fmt"
  _ "bufio"
  _ "strings"
)


type Board struct {
  squareTab [20][20] Square
}

type Square struct {
  x int
  y int
  empty bool
}

type Piece struct {
  id int
  origin Square
  rotation string
  flipped bool
  player Player
}

type Player struct {
  name string
  color string
}


func main() {
  fmt.Println("#### Test ####")
  player := Player{"Bertrand", "yellow"}
  fmt.Println(player)
  mySquare := Square{1, 2, false}
  var tab [20][20]Square
  var board Board
  tab[0][0] = mySquare
  board.squareTab = tab
  fmt.Println(board.squareTab)
}
/*
func main() {

  fmt.Println("Launching server on port 8081...")

  // listen on all interfaces
  ln, _ := net.Listen("tcp", ":8081")

  // accept connection on port
  conn, _ := ln.Accept()
  // run loop forever (or until ctrl-c)
  for {
    // will listen for message to process ending in newline (\n)
    message, _ := bufio.NewReader(conn).ReadString('\n')
    // output message received
    fmt.Print("Message Received:", string(message))
    // sample process for string received
    newmessage := strings.ToUpper(message)
    // send new string back to client
    conn.Write([]byte(newmessage + "\n"))

    if (string(message) == "QUIT") {
      return
    }
  }
}*/