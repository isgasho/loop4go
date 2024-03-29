package main

import (
	"fmt"
	"github.com/smartwalle/loop4go"
	"time"
)

func main() {
	var queue = loop4go.NewEventQueue()

	var loop loop4go.Loop
	var count = 0
	loop = loop4go.NewLoop(time.Second, queue, func(l loop4go.Loop) {
		count++
		fmt.Println(time.Now(), count)

		if count > 2 {
			queue.Stop()
		}
	})
	loop.Start()

	queue.Start()
	queue.Wait()
}
