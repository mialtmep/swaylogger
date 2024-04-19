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
var kb_index int

var data map[int]int
var ru_data map[int]int
var en_data map[int]int
var sp_ru_data map[int]int
var sp_en_data map[int]int

func showResults() {
	for k, v := range data {
		if k > 1000 {
			if (k >= 1016 && k <= 1027) || (k >= 1030 && k <= 1041) || (k >= 1044 && k <= 1052) {
				ru_data[k] = v
			} else {
				sp_ru_data[k] = v
			}
		} else {
			if (k >= 16 && k <= 25) || (k >= 30 && k <= 38) || (k >= 44 && k <= 50) {
				en_data[k] = v
			} else {
				sp_en_data[k] = v
			}
		}
	}
	ru_json, _ := json.MarshalIndent(ru_data, "", "  ")
	en_json, _ := json.MarshalIndent(en_data, "", "  ")
	ru_sp_json, _ := json.MarshalIndent(sp_ru_data, "", "  ")
	en_sp_json, _ := json.MarshalIndent(sp_en_data, "", "  ")

	fmt.Println("\nru_data = ", string(ru_json))
	fmt.Println("en_data = ", string(en_json))
	fmt.Println("sp_ru_data = ", string(ru_sp_json))
	fmt.Println("sp_en_data = ", string(en_sp_json))
}

func main() {
	data = map[int]int{}
	ru_data = map[int]int{}
	en_data = map[int]int{}
	sp_ru_data = map[int]int{}
	sp_en_data = map[int]int{}
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

	c := make(chan os.Signal)
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
			kb_index = i
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
					layout = *inputs[kb_index].XKBActiveLayoutName
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
