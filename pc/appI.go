package main

/*
#include "callback.h"
*/
import "C"

func main() {
}

type App struct {
	callback C.UserInterfaceAPI
	setErr   C.SetLastErr
}

//export InitLib
func InitLib() {
}

//export OpenWallet
func OpenWallet(auth string) bool {
	return false
}
