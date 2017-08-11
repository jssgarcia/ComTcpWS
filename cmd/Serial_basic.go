package main

import (
	rs "github.com/tarm/serial"
	"time"
	"log"
)

func main() {
	port0 := "COM1"
	port1 := "COM2"

	c0 := &rs.Config{Name: port0, Baud: 9600,Size:8,StopBits:1,Parity:0}
	c1 := &rs.Config{Name: port1, Baud: 9600,Size:8,StopBits:1,Parity:0}

	s1, err := rs.OpenPort(c0)
	if err != nil {
		log.Fatal(err)
	}

	s2, err := rs.OpenPort(c1)
	if err != nil {
		log.Fatal(err)
	}

	ch := make(chan int, 1)
	go func() {
		buf := make([]byte, 128)
		var readCount int
		for {
			n, err := s2.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			readCount++
			log.Printf("Read %v %v bytes: % 02x %s", readCount, n, buf[:n], buf[:n])
			select {
			case <-ch:
				ch <- readCount
				close(ch)
			default:
			}
		}
	}()

	if _, err = s1.Write([]byte("hello")); err != nil {
		log.Fatal(err)
	}
	if _, err = s1.Write([]byte(" ")); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second)
	if _, err = s1.Write([]byte("world")); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second / 10)

	ch <- 0
	s1.Write([]byte(" ")) // We could be blocked in the read without this
	c := <-ch
	exp := 5
	if c >= exp {
		log.Fatalf("Expected less than %v read, got %v", exp, c)
	}
}



