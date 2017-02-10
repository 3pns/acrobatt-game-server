package utils

import (
  "fmt"
)

func SetWhiteBackground (){
  fmt.Print("\x1B[47m")
}
func SetBlackBackground (){
  fmt.Print("\x1B[40m")
}

func PrintReset(){
  fmt.Print("\x1B[0m")
}

func PrintBlack(str string){
  fmt.Print("\x1B[30m" + str)
}

func PrintRed(str string){
  fmt.Print("\x1B[31m" + str)
}

func PrintGreen(str string){
  fmt.Print("\x1B[32m" + str)
}

func PrintYellow(str string){
  fmt.Print("\x1B[33m" + str)
}

func PrintBlue(str string){
  fmt.Print("\x1B[34m" + str)
}

func PrintWhite(str string){
  fmt.Print("\x1B[37m" + str)
}