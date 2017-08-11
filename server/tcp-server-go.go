package server

import (
	"net"
	"context"
	"fmt"
	_ "log"
	lg "ty/csi/ws/ComTcpWS/lgg"
	"strings"
	"errors"
)

type ServerTcpInfo struct {
	ServerTcpAddress string
	ChannelReceiveData chan string		//Canal para enviar datos al servidor

	serverOptionsTx string
	conn net.Conn
}

var modName = "SERVER-TCP"
var _conn net.Conn

func InitServer(ctx context.Context,info *ServerTcpInfo) {

	if info.ServerTcpAddress =="" {
		panic("ERROR Initiation: Server adddres is not defined. " + info.ServerTcpAddress)
	}

	lg.Lgdef.Infof("[%s]Launching server ... %v",modName, info.ServerTcpAddress)

	ln,err := net.Listen("tcp", info.ServerTcpAddress)
	if err!=nil {
		panic("Interface ")
	}

	info.conn = nil


	go func() {

		for {

			select {
			case data := <- info.ChannelReceiveData:
				//Dato recivido desde le Puerto

				if (info.conn != nil) {
					//Envio datos al cliente
					lg.Lgdef.Printf("[%s]: DATARAW RECEIVED FOR SEND ... %s \n",modName, data)
					sendData(info.conn,data)
				}

				break
			}
		}

	}()

	for {
		conn, _ := ln.Accept()
		info.conn = conn;

		//Cliente conectado...
		lg.Lgdef.Infof("[%s] Cliente conectado... ",modName)

		////
		//tmr := time.NewTicker(2 * time.Second)
		//
		//bcontinue := true
		//
		//for bcontinue {
		//	select {
		//	case <-tmr.C:
		//
		//		data := strconv.Itoa(getNextRandomValue()) + ".\r"
		//
		//		err := sendData(conn, data)
		//		if err != nil {
		//			fmt.Printf("SERVER: error send data '%s'\n",err)
		//			bcontinue=false
		//		}
		//
		//		break
		//	}
		//}
	}

	fmt.Printf("SERVER == Finalizado main ==")
}

func sendData(conn net.Conn, data string) error {

	dataparsed,err := parseData(data)
	if (err!=nil) {
		lg.Lgdef.Error(err)
		return err
	}

	_,errWrite := conn.Write([]byte(dataparsed+"\r"))
	lg.Lgdef.Infof("[%s]: DATA SEND ... %s \r\n",modName, dataparsed)

	return errWrite
}

func parseData(data string) (string,error){

	//FORMATO DE DATOS ESPERADO
	//    06:10   11.6 23.1ºC\n
	/*
	a) Trim b) Reemplazar '   ' x ' ' c) Dividir por espacios d) Obtener pos(1) */
	formatspec := "    06:10   11.6 23.1ºC"
	newdata := strings.TrimSpace(data)
	newdata = strings.Replace(data,"   "," ",-1)
	newdata = removeNoValidCharacters(strings.TrimSpace(newdata))
	parts := strings.Split(newdata," ")

	if (len(parts)<1){
		lg.Lgdef.Errorf("[%s] PARSE DATA. '%s'.  NOT FORMAT SPEC. #%s#",modName,data,formatspec)
		return  "",errors.New("No format spec. Format expexted #" + formatspec + "#")
	}

	partsNoEmpty := 0
	for i:=0;i<=len(parts);i++  {
		if (strings.TrimSpace(parts[i])!="") {
			partsNoEmpty++
			if (partsNoEmpty==2) {
				return strings.TrimSpace(parts[i]), nil
			}
		}
	}

	return "",errors.New("ERROR procesing format. No format spec. Format expexted #" + formatspec + "#")

}

func removeNoValidCharacters(msg string) string {
	b := make([]byte, len(msg))
	var bl int
	for i := 0; i < len(msg); i++ {
		c := msg[i]
		if c >= 32 && c < 127 {
			b[bl] = c
			bl++
		}
	}
	return string(b[:bl])
}


