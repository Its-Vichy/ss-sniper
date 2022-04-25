package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/inancgumus/screen"
)

func log(content string) {
	fmt.Printf("[%s] %s.\n", time.Now().Format("02-Jan-2006 15:04:05"), content)

	file, err := os.OpenFile("Logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}

	_, err = file.WriteString(fmt.Sprintf("[%s] %s.\n", time.Now().Format("02-Jan-2006 15:04:05"), content))
	if err != nil {
		return
	}
}

func print_logo() {
	screen.Clear()
	screen.MoveTopLeft()

	fmt.Printf(`
	 ___ ___     ___      _               
	/ __/ __|___/ __|_ _ (_)_ __  ___ _ _ 
	\__ \__ \___\__ \ ' \| | '_ \/ -_) '_|
	|___/___/   |___/_||_|_| .__/\___|_|  
	                       |_| V ~~> %s
						   
`, version)
}

func block_console() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
