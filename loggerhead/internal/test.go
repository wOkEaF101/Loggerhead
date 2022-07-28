package loggerhead

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/bi-zone/etw"
	"golang.org/x/sys/windows"
)

func TestETW() {
	guid, _ := windows.GUIDFromString("{7DD42A49-5329-4832-8DFD-43D979153A88}")
	session, err := etw.NewSession(guid)
	if err != nil {
		log.Fatalf("Failed to create ETW session: %s", err)
	}

	cb := func(e *etw.Event) {
		if data, err := e.EventProperties(); err == nil {
			if err != nil {
				log.Fatalf(err.Error())
			}
			log.Printf("%s", data)
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
