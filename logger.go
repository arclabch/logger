// Copyright (C) 2018 ARClab, Lionel Riem - https://arclab.ch/
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

// This package simplifies logging to Syslog or the screen, depending on the
// "verbose" setting.

package logger

import (
	"fmt"
	"errors"
	"os"
	"time"
	"log/syslog"

	"github.com/mattn/go-isatty"
)

const (

	// Logging levels
	L_EMERGENCY  = 0
	L_ALERT      = 1
	L_CRITICAL   = 2
	L_ERROR      = 3
	L_WARNING    = 4
	L_NOTICE     = 5
	L_INFO       = 6
	L_DEBUG      = 7

	// Message header
	M_EMERGENCY  = "EMERGENCY"
	M_ALERT      = "ALERT    "
	M_CRITICAL   = "CRITICAL "
	M_ERROR      = "ERROR    "
	M_WARNING    = "WARNING  "
	M_NOTICE     = "NOTICE   "
	M_INFO       = "INFO     "
	M_DEBUG      = "DEBUG    "
)

var (
	s            *syslog.Writer
	debug        = false
	verbose      = false
	color        = true

	// Colors
	C_BLUE       = string([]byte{27, 91, 57, 55, 59, 52, 52, 109})
	C_CYAN       = string([]byte{27, 91, 57, 55, 59, 52, 54, 109})
	C_GREEN      = string([]byte{27, 91, 57, 55, 59, 52, 50, 109})
	C_WHITE      = string([]byte{27, 91, 57, 48, 59, 52, 55, 109})
	C_YELLOW     = string([]byte{27, 91, 57, 55, 59, 52, 51, 109})
	C_MAGENTA    = string([]byte{27, 91, 57, 55, 59, 52, 53, 109})
	C_RED        = string([]byte{27, 91, 57, 55, 59, 52, 49, 109})
	C_RESET      = string([]byte{27, 91, 48, 109})

)

// Starts the logging system.
// Takes a tag parameter to specify the name of the program.
// Returns an error if unable to start logging.
func Open(tag string) error {
	var err error

	if tag != "" {
		s, err =  syslog.New(syslog.LOG_WARNING|syslog.LOG_DAEMON, tag)
		if err != nil {
			return err
		}
	} else {
		return errors.New("logger: tag cannot be empty")
	}

	if (os.Getenv("TERM") == "dumb" || (!isatty.IsTerminal(os.Stdout.Fd()) && !isatty.IsCygwinTerminal(os.Stdout.Fd()))) {
		color = false
	}

	return nil
}

// Stops the logging system.
// Should be called at the end of the program.
// Returns an error if unable to stop logging.
func Close() error {
	err := s.Close()
	return err
}

// Prints a message to the screen.
// Will check if color can be used or not.
func PrintToScreen(level int, message string) {
	var (
		mColor  string
		mReset  string
		mHeader string
	)

	switch level {
	case L_EMERGENCY:
		mColor  = C_RED
		mHeader = M_EMERGENCY
	case L_ALERT:
		mColor  = C_RED
		mHeader = M_ALERT
	case L_CRITICAL:
		mColor  = C_YELLOW
		mHeader = M_CRITICAL
	case L_ERROR:
		mColor  = C_YELLOW
		mHeader = M_ERROR
	case L_WARNING:
		mColor  = C_MAGENTA
		mHeader = M_WARNING
	case L_NOTICE:
		mColor  = C_CYAN
		mHeader = M_NOTICE
	case L_INFO:
		mColor  = C_WHITE
		mHeader = M_INFO
	case L_DEBUG:
		mColor  = C_GREEN
		mHeader = M_DEBUG
	}
	mReset = C_RESET

	if(!color) {
		mColor = ""
		mReset = ""
	}

	fmt.Printf("%s - %s%s%s %s\n", time.Now().Format(time.RFC3339), mColor, mHeader, mReset, message)
}

// Disable colors in messages printed to screen.
func DisableColor() {
	color = false
}

// Sets the logging to debug mode using the supplied boolean.
// When set to true, logs will be printed on screen instead of being
// sent to Syslog.
// Off (false) by default.
func SetDebug(b bool) {
	debug = b
}

// Sets the logging to verbose mode using the supplied boolean.
// When set to true, Info and Debug level messages will be logged.
// Otherwise, they are simply ignored.
// Off (false) by default.
func SetVerbose(b bool) {
	verbose = b
}

// Logs an Emergency-evel event.
// Emergency messages will always be sent to Syslog and printed on screen.
// Returns an error if unable to log it.
func Emergency(message string) error {
	PrintToScreen(L_EMERGENCY, message)
	err := s.Emerg(message)
	return err
}

// Logs an Alert-level event.
// Returns an error if unable to log it.
func Alert(message string) error {
	if debug == true {
		PrintToScreen(L_ALERT, message)
		return nil
	} else {
		err := s.Alert(message)
		return err
	}
}

// Logs a Critical-level event.
// Returns an error if unable to log it.
func Critical(message string) error {
	if debug == true {
		PrintToScreen(L_CRITICAL, message)
		return nil
	} else {
		err := s.Crit(message)
		return err
	}
}

// Logs an Error-level event.
// Returns an error if unable to log it.
func Error(message string) error {
	if debug == true {
		PrintToScreen(L_ERROR, message)
		return nil
	} else {
		err := s.Err(message)
		return err
	}
}

// Logs a Warning-level event.
// Returns an error if unable to log it.
func Warning(message string) error {
	if debug == true {
		PrintToScreen(L_WARNING, message)
		return nil
	} else {
		err := s.Warning(message)
		return err
	}
}

// Logs a Notice-level event.
// Returns an error if unable to log it.
func Notice(message string) error {
	if debug == true {
		PrintToScreen(L_NOTICE, message)
		return nil
	} else {
		err := s.Notice(message)
		return err
	}
}

// Logs an Info-level event.
// Will not be logged unless Verbose is set to true.
// Returns an error if unable to log it.
func Info(message string) error {
	if verbose == true {
		if debug == true {
			PrintToScreen(L_INFO, message)
			return nil
		} else {
			err := s.Info(message)
			return err
		}
	}
	return nil
}

// Logs a Debug-level event.
// Will not be logged unless Verbose is set to true.
// Returns an error if unable to log it.
func Debug(message string) error {
	if verbose == true {
		if debug == true {
			PrintToScreen(L_DEBUG, message)
			return nil
		} else {
			err := s.Debug(message)
			return err
		}
	}
	return nil
}
