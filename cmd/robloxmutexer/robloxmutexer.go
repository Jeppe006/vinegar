//go:build windows
// +build windows

package main

import (
	"errors"
	"log"

	"golang.org/x/sys/windows"
)

func main() {
	log.SetPrefix("robloxmutexer: ")
	log.SetFlags(log.Lmsgprefix | log.LstdFlags)

	name, err := windows.UTF16PtrFromString("ROBLOX_singletonMutex")
	if err != nil {
		log.Fatal(err)
	}

	handle, err := windows.CreateMutex(nil, false, name)
	if err != nil {
		if errors.Is(err, windows.ERROR_ALREADY_EXISTS) {
			log.Fatal("Roblox's Mutex is already locked!")
		} else {
			log.Fatal(err)
		}
	}

	_, err = windows.WaitForSingleObject(handle, 0)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Locked Roblox singleton Mutex")

	_, _ = windows.WaitForSingleObject(windows.CurrentProcess(), windows.INFINITE)
}
