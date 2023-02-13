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
	CMDSrvStatusChanged
	CMDLogOutPut
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

func (app *App) Write(p []byte) (n int, err error) {
	msg := &CallBackMsg{
		Cmd:   CMDLogOutPut,
		Param: string(p),
	}
	app.CBJsonData(msg)
	return len(p), nil
}

type CallBackMsg struct {
	Cmd   int    `json:"cmd"`
	Param string `json:"param"`
}
