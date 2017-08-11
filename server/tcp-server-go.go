package main
//
//import (
//	"net"
//	//"bufio"
//	"fmt"
//	_ "log"
//	"github.com/tkanos/gonfig"
//	"os"
//	"path/filepath"
//	"time"
//	"math/rand"
//	"strconv"
//	"strings"
//)
//
//type configuration struct {
//	ServerTcpAddress string `json:"ServerTcpAddress"`
//}
//
//var cnfg configuration
//var _conn net.Conn
//
//func main() {
//
//	//getConfig()
//
//	if len(os.Args) > 1 {
//		if (len(os.Args)>=1) { cnfg.ListenAddress = os.Args[1] }
//		if (len(os.Args)>=2) {
//			cnfg.RandomInterval = parseRandomRange(os.Args[2])
//		}
//	} else {
//		panic("ServerAdress is not defined")
//	}
//
//	fmt.Printf("Launching server ... %v",cnfg.ListenAddress)
//
//	ln,err := net.Listen("tcp", cnfg.ListenAddress)
//	if err!=nil {
//		panic("Interface ")
//	}
//
//	//Random seed initializion
//	rand.Seed(time.Now().Unix())
//
//	for {
//		conn, _ := ln.Accept()
//		_conn = conn;
//
//		//
//		tmr := time.NewTicker(2 * time.Second)
//
//		bcontinue := true
//
//		for bcontinue {
//			select {
//			case <-tmr.C:
//
//				data := strconv.Itoa(getNextRandomValue()) + ".\r"
//
//				err := sendData(conn, data)
//				if err != nil {
//					fmt.Printf("SERVER: error send data '%s'\n",err)
//					bcontinue=false
//				}
//
//				break
//			}
//		}
//	}
//
//	fmt.Printf("SERVER == Finalizado main ==")
//}
//
//func getNextRandomValue() int{
//	return  rand.Intn(cnfg.RandomInterval.Max-cnfg.RandomInterval.Min)+cnfg.RandomInterval.Min
//}
//
//func sendData(conn net.Conn, data string) error {
//	_,err := conn.Write([]byte(data+"\n"))
//	fmt.Println("TCPServer: message send ..." + data)
//
//	return err
//}
//
//func getConfig() {
//
//	cnfg = configuration{}
//
//	fconfig := "config.dev.json"
//
//	//Defaults
//	//cnfg.ListenAddress="127.0.0.1:8888"
//
//	f,_ := filepath.Abs("./server/" + fconfig)
//
//	//Obtenemos la configuracion
//	err := gonfig.GetConf(f,&cnfg)
//	if err!=nil {
//		panic("ERROR obtener configuracion: " + err.Error())
//	}
//
//}
//
//func parseRandomRange(strrange string) (rg *randomInterval) {
//
//	//Defaults
//	rg = &randomInterval{Max: 1000, Min: 0}
//
//	if strrange == "" || !strings.Contains(strrange, ":") {
//		return rg
//	}
//
//	defer func() {
//		//Simplemente TryCatch
//	}()
//
//	s := strings.Split(strrange, ":")
//	if len(s) <= 1 {
//		return rg
//	}
//
//	rg.Max, _ = strconv.Atoi(s[0])
//	rg.Min, _ = strconv.Atoi(s[1])
//
//	return rg
//}
//
////func startTcpclient() {
////
////	ln,err := net.Listen("tcp","127.0.0.1:35349")
////	if err!=nil {
////		log.Fatal("FATAL")
////	}
////
////}
//
//
