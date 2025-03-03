package main

import "log"

func LogInfo(msg string) {
	log.Println("INFO - ", msg)
}

func LogWarning(msg string) {
	log.Println("WARNING - ", msg)
}

func LogError(msg string) {
	log.Println("ERROR - ", msg)
}
