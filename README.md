logger
======

`logger` is a simple Go package wrapping around Go's `log\syslog` package.

Usage example:

```go
package main

import (
        "fmt"
        "os"
        "github.com/arclabch/logger"
)

func main() {
        err := logger.Open("Test")
        if err !=nil {
                fmt.Printf("Fatal error logger: %s\n", err.Error())
                os.Exit(-1)
        }

        logger.SetDebug(true)
        logger.SetVerbose(true)

        m := "Test Message - Disregard."
        logger.Emergency(m)
        logger.Alert(m)
        logger.Critical(m)
        logger.Error(m)
        logger.Warning(m)
        logger.Notice(m)
        logger.Info(m)
        logger.Debug(m)

        logger.Close()
}
```