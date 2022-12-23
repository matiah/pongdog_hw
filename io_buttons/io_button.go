package io_button

import (
	"fmt"
	"time"

	"github.com/warthog618/gpiod"
	"github.com/warthog618/gpiod/device/rpi"
)

func handler(evt gpiod.LineEvent) {
	go notify(evt.Offset, evt.Timestamp)
}

func main() {
	cc := gpiod.Chips()
	fmt.Println(cc)
	l, _ := gpiod.RequestLine("gpiochip0", rpi.J8p37, gpiod.AsInput, gpiod.WithPullUp, gpiod.AsActiveLow, gpiod.WithEventHandler(handler), gpiod.WithRisingEdge, gpiod.WithDebounce(time.Millisecond*20))
	fmt.Println(l)
	for {
		select {}
	}
}
func notify(line int, timestamp time.Duration) {
	fmt.Printf("%s - Rising edge detected, button %d\n", timestamp, line)
}
