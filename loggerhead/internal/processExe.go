package loggerhead

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/bi-zone/etw"
	"github.com/mitchellh/mapstructure"
	"golang.org/x/sys/windows"
)

type ProcessExeStruct struct {
	CreateTime             string
	ImageName              string
	ProcessID              string
	ProcessTokenIsElevated string
	SessionID              string
}

func ProcETW() {
	guid, _ := windows.GUIDFromString("{22FB2CD6-0E7B-422B-A0C7-2FAD1FD0E716}")
	session, err := etw.NewSession(guid)
	if err != nil {
		log.Fatalf("Failed to create ETW session: %s", err)
	}

	cb := func(e *etw.Event) {
		if e.Header.ID != 1 {
			return
		}
		if data, err := e.EventProperties(); err == nil {
			Process := &ProcessExeStruct{}
			mapstructure.Decode(data, &Process)
			procJSON, err := json.MarshalIndent(Process, "", " ")
			if err != nil {
				log.Fatalf(err.Error())
			}
			log.Printf("%s", procJSON)
		}
	}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := session.Process(cb); err != nil {
			log.Printf("[ERR] Got error processing events: %s", err)
		}
		wg.Done()
	}()

	// Trap cancellation.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	<-sigCh

	if err := session.Close(); err != nil {
		log.Printf("[ERR] Got error closing the session: %s", err)
	}
	wg.Wait()
}
