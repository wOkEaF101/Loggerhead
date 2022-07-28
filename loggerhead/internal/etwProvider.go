package loggerhead

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/bi-zone/etw"
	"golang.org/x/sys/windows"
)

func ETWSession(GUID string) {
	guid, _ := windows.GUIDFromString(GUID)
	session, err := etw.NewSession(guid)
	if err != nil {
		log.Fatalf("Failed to create ETW session: %s", err)
	}

	cb := func(e *etw.Event) {
		if e.Header.ID != 1 {
			return
		}
		if data, err := e.EventProperties(); err == nil {
			if err != nil {
				log.Fatalf(err.Error())
			}
			output, err := json.MarshalIndent(data, "", " ")
			log.Printf("%s", output)
			if err != nil {
				log.Fatalf(err.Error())
			}
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
