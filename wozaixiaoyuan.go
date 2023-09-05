package main

import (
	"github.com/robfig/cron"
	"log"
	"wobuzaixiaoyuan/pkg/logServer"
)

func main() {
	log.SetFlags(log.Ltime | log.Ldate)
	logServer.StartLogServer()
	instance, err := NewInstance()
	if err != nil {
		log.Fatal(err)
	}
	c := cron.New()
	updateSpec := "0 30 * * * *"
	checkinSpec := "0 */5 * * * *"
	resetSpec := "0 0 0 * * *"
	err = c.AddFunc(updateSpec, instance.UpdateData)
	if err != nil {
		log.Fatal(err)
	}
	err = c.AddFunc(checkinSpec, instance.CheckInTask)
	if err != nil {
		log.Fatal(err)
	}
	err = c.AddFunc(resetSpec, instance.ResetEventMap)
	if err != nil {
		log.Fatal(err)
	}
	c.Start()
	select {}
}
