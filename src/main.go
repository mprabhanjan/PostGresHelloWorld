package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "sync"
    "syscall"
)

func main() {
    DbInitialize()

    wg := &sync.WaitGroup{}
    errChnl := make(chan error, 1)
    signalChan := make(chan os.Signal, 0)
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

    server, err := StartWebServer(8080, wg, errChnl)
    if err != nil {
        log.Fatalf("Unable to start webserver. Error - %v", err)
    }

    for {
        select {
        case osSig := <- signalChan:
            log.Printf("Rcvd signal %v. Calling server Shutdown()", osSig)
            server.Shutdown(context.Background())

        case err := <- errChnl:
            log.Printf("Got notified on channel about listener shutdown - error %v", err)
            wg.Wait()
            return
        }
    }
}