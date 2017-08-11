package main

import (
	rs "github.com/tarm/serial"
	_ "time"
	"log"
	"flag"
	"fmt"
	"os"
	"time"
)

//func main() {
//	port0 := "COM1"
//	port1 := "COM2"
//
//	c0 := &rs.Config{Name: port0, Baud: 9600,Size:8,StopBits:1,Parity:0}
//	c1 := &rs.Config{Name: port1, Baud: 9600,Size:8,StopBits:1,Parity:0}
//
//	s1, err := rs.OpenPort(c0)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	s2, err := rs.OpenPort(c1)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	ch := make(chan int, 1)
//	go func() {
//		buf := make([]byte, 128)
//		var readCount int
//		for {
//			n, err := s2.Read(buf)
//			if err != nil {
//				log.Fatal(err)
//			}
//			readCount++
//			log.Printf("Read %v %v bytes: % 02x %s", readCount, n, buf[:n], buf[:n])
//			select {
//			case <-ch:
//				ch <- readCount
//				close(ch)
//			default:
//			}
//		}
//	}()
//
//	if _, err = s1.Write([]byte("hello")); err != nil {
//		log.Fatal(err)
//	}
//
//	if _, err = s1.Write([]byte(" ")); err != nil {
//		log.Fatal(err)
//	}
//	time.Sleep(time.Second)
//	if _, err = s1.Write([]byte("world")); err != nil {
//		log.Fatal(err)
//	}
//	time.Sleep(time.Second / 10)
//
//	ch <- 0
//	s1.Write([]byte(" ")) // We could be blocked in the read without this
//	c := <-ch
//	exp := 5
//	if c >= exp {
//		log.Fatalf("Expected less than %v read, got %v", exp, c)
//	}
//}


func main() {

	log.Printf("SERIAL-READER: INIT... ")

	port := flag.String("port", "COM1", "serial port name (/dev/ttyUSB0,COMX, etc)")
	baud := flag.Int("baud", 9600, "Baud rate")
	stopbits := flag.Int("stopbits", 1, "Stop bits")
	databits := flag.Uint("databits", 8, "Data bits")
	parity :=flag.Uint("parity", 0, "Parity mode")

	flag.Parse()

	if *port == "" {
		fmt.Println("Must specify port")
		usage()
	}

	c0 := &rs.Config{Name: *port, Baud: int(*baud),Size:byte(*databits),StopBits: rs.StopBits(*stopbits),Parity:rs.Parity(*parity)}

	log.Printf("SERIAL-READER: [%s]. CONFIG",c0)

	read(c0)

	log.Printf("SERIAL-READER: [%s]. FINISH",c0.Name)


}

func read(cnfg *rs.Config) {

	conn, err := rs.OpenPort(cnfg)
	if err != nil {
		log.Fatal(err)
	}

	readData := ""
	readCount :=0

	timer := time.NewTimer(5 * time.Second)
	timer.Stop()
	//Sincro channel
	//chsync := make(chan int, 0)

	go func() {
		for {
			select {
			case <- timer.C:
				if (readData=="" ){
					continue
				}

				log.Printf("DATA READED: Bytes:%n Dato:%s",readCount,readData)
				timer.Stop()

				//Asignamos
				readCount=0
				readData=""

				break
			}
		}
	}()

	buf := make([]byte, 128)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		if readCount==0 {
			timer.Reset(5 * time.Second)
			log.Printf("DATA READ. Timer Stated >>>")
		}

		readCount++
		readData = readData + string(buf[:n])
		log.Printf("Read %v %v bytes: % 02x %s", readCount, n, buf[:n], buf[:n])
	}
}




func usage() {
	fmt.Println("go-serial-test usage:")
	flag.PrintDefaults()
	os.Exit(-1)
}


