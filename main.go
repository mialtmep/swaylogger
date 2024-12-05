package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	evdev "github.com/gvalkov/golang-evdev"
	"github.com/joshuarubin/go-sway"
)

// instructions

// 1
// get sway pid
// ps -axu | grep -v 'grep' | grep -iP "sway\$"

// 2
// get input id
// swaymsg -t get_inputs -r | jq "."

// 3
// in /dev/input/by-id find correct input id
// in /proc/$PID/fd find correct file descriptor symlinked to the input file

// 4
// arg 1 - path to the symlinked file of sway process
// arg 2 - sway input id

var layout string
var kbIndex int

var data map[int]int
var ruData map[int]int
var enData map[int]int
var ruSpData map[int]int
var enSpData map[int]int

func showResults() {
	for k, v := range data {
		if k > 1000 {
			if (k >= 1016 && k <= 1027) || (k >= 1030 && k <= 1041) || (k >= 1044 && k <= 1052) {
				ruData[k] = v
			} else {
				ruSpData[k] = v
			}
		} else {
			if (k >= 16 && k <= 25) || (k >= 30 && k <= 38) || (k >= 44 && k <= 50) {
				enData[k] = v
			} else {
				enSpData[k] = v
			}
		}
	}
	ruJSON, _ := json.MarshalIndent(ruData, "", "  ")
	enJSON, _ := json.MarshalIndent(enData, "", "  ")
	ruSpJSON, _ := json.MarshalIndent(ruSpData, "", "  ")
	enSpJSON, _ := json.MarshalIndent(enSpData, "", "  ")

	fmt.Println("\nru_data = ", string(ruJSON))
	fmt.Println("en_data = ", string(enJSON))
	fmt.Println("sp_ru_data = ", string(ruSpJSON))
	fmt.Println("sp_en_data = ", string(enSpJSON))
}

func main() {
	data = map[int]int{}
	ruData = map[int]int{}
	enData = map[int]int{}
	ruSpData = map[int]int{}
	enSpData = map[int]int{}
	ctx := context.Background()
	client, err := sway.New(ctx)
	if err != nil {
		log.Fatal(err)
	}

	dev, err := evdev.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer dev.File.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		showResults()
		os.Exit(1)
	}()

	inputs, err := client.GetInputs(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for i, v := range inputs {
		if v.Type == "keyboard" && v.Identifier == os.Args[2] {
			kbIndex = i
		}
	}

	for {
		events, err := dev.Read()
		if err != nil {
			log.Fatal(err)
		}

		for _, ev := range events {
			if ev.Type == evdev.EV_KEY {
				if ev.Value == 1 {
					inputs, err = client.GetInputs(ctx)
					if err != nil {
						log.Fatal(err)
					}
					layout = *inputs[kbIndex].XKBActiveLayoutName
					if layout == "English (US)" {
						data[int(ev.Code)] += 1
					} else if layout == "Russian" {
						data[int(ev.Code+1000)] += 1
					}
				}
			}
		}
	}
}
