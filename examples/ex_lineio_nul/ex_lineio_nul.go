package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/l4go/cmdio"
	"github.com/l4go/lineio"
	"github.com/l4go/task"
)

func main() {
	log.Println("START")
	defer log.Println("END")

	m := task.NewMission()

	signal_ch := make(chan os.Signal, 1)
	signal.Notify(signal_ch, syscall.SIGINT, syscall.SIGTERM)

	std_rw, err := cmdio.StdDup()
	if err != nil {
		defer log.Println("Error:", err)
		return
	}
	go func(cm *task.Mission) {
		defer std_rw.Close()
		echo_worker(cm, std_rw)
	}(m.New())

	select {
	case <-m.Recv():
	case <-signal_ch:
		m.Cancel()
	}
}

func echo_worker(m *task.Mission, rw io.ReadWriter) {
	defer m.Done()
	log.Println("start: echo worker")
	defer log.Println("end: echo worker")

	line_r := lineio.NewReaderByDelim(rw, 0x00)
	for {
		var ln []byte
		var ok bool
		select {
		case ln, ok = <-line_r.Recv():
		case <-m.RecvCancel():
			return
		}
		if !ok {
			break
		}

		fmt.Fprintln(rw, ">", string(ln))
	}
	if err := line_r.Err(); err != nil {
		log.Println("Error:", err)
		return
	}
}
