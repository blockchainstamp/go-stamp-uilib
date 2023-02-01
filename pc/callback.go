package main

/*
#include "callback.h"

char* bridge_func(UserInterfaceAPI f , const char* v){
	return f(v);
}


void bridge_Error(SetLastErr f , const char* v){
	f(v);
}
*/
import "C"
import "encoding/json"

const (
	_ = iota
	LoadDns
	LoadInnerIP
)

func (app *App) CBJsonData(msg *CallBackMsg) string {
	bts, _ := json.Marshal(msg)
	result := C.bridge_func(app.callback, C.CString(string(bts)))
	if result == nil {
		return ""
	}
	return C.GoString(result)
}

func (app *App) SetError(str string) {
	C.bridge_Error(app.setErr, C.CString(str))
}

type CallBackMsg struct {
	Cmd   int    `json:"cmd"`
	Param string `json:"param"`
}
