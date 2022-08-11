package main

import (
	loggerhead "loggerhead/internal"
)

//type EventTrace struct {
//	GUID  string
//	Event []uint16
//	PID   []string
//}

func main() {
	//net := loggerhead.EventTrace{GUID: "{E13C0D23-CCBC-4E12-931B-D9CC2EEE27E4}", Event: 252} //.NET
	proc := loggerhead.EventTrace{GUID: "{22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}", Event: 1}  //Process
	conn := loggerhead.EventTrace{GUID: "{7DD42A49-5329-4832-8DFD-43D979153A88}", Event: 11} //Network
	//loggerhead.ETWSession("{0A002690-3839-4E3A-B3B6-96D8DF868D99}") //Malz
	loggerhead.ETWSession(proc)
	loggerhead.ETWSession(conn)
}
