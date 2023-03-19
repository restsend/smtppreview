package main

import "log"

const (
	ColorBlack  = "\u001b[30m"
	ColorRed    = "\u001b[31m"
	ColorGreen  = "\u001b[32m"
	ColorYellow = "\u001b[33m"
	ColorBlue   = "\u001b[34m"
	ColorReset  = "\u001b[0m"
)
const logFlags = log.Lshortfile | log.LstdFlags

func colorize(color, message string) string {
	return color + message + ColorReset
}
