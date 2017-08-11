package main

import (
	"context"
	"ty/csi/ws/ComTcpWS/Global"
	lg "ty/csi/ws/ComTcpWS/lgg"
	"ty/csi/ws/ComTcpWS/Utils"
	rs "github.com/tarm/serial"
	"time"
	"log"
)

type SerialReaderInfo struct {
	OpenOptions Global.SerialOption
	openOptionsTx string
	ChannelReceiveData chan string		//Canal para enviar datos al servidor

	portConn *rs.Port
}

var modName = "SERIAL-READER"

func initSerial(ctx context.Context,info *SerialReaderInfo) {

	info.openOptionsTx = "SerialOptions: [" + Utils.PrettyPrint(info.openOptionsTx) + "]"

	lg.Lgdef.Info("SERIAL-READER:: INIT " + info.openOptionsTx)

	ticker := time.NewTicker(5 * time.Second)

	defer func() {
		lg.Lgdef.Printf("[%s] initSerial exit",modName)
	}()

	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			lg.Lgdef.Warnf("%s(%s): Cancelación recibida: Salimos", modName,info.openOptionsTx)
			disposeSerial(info)
			return

		case <- ticker.C:
			err := readSerialPort(ctx,info)
			if err!=nil {
				lg.Lgdef.Error(err)
			}

			break
		}
	}
}

func readSerialPort(ctx context.Context,info *SerialReaderInfo) error {

	cnfg := &rs.Config{
		Name: info.OpenOptions.Portname,
		Baud: int(info.OpenOptions.Baudrate),
		Size:byte(info.OpenOptions.Databits),
		StopBits: rs.StopBits(info.OpenOptions.Stopbits),
		Parity:rs.Parity(info.OpenOptions.ParityMode)}

	readData := ""
	readCount :=0

	tmrCheckRead := time.NewTimer(5 * time.Second)
	tmrCheckRead.Stop()

	//Check if Data received for Port
	go func() {
		for {
			select {
			case <- tmrCheckRead.C:
				if readData=="" {
					continue
				}

				log.Printf("[%s] DATA READED: Bytes:%n Dato:%s",modName, readCount,readData)
				tmrCheckRead.Stop()

				sendDataToServer(info,readData)

				//Asignamos
				readCount=0
				readData=""

				break
			}
		}
	}()

	// Open the port.
	conn, err := rs.OpenPort(cnfg)
	if err != nil {
		lg.Lgdef.Errorf("%s: ERROR Open Serial Port. %v",modName, err)
		return err
	}

	info.portConn = conn

	defer func() {
		if r := recover(); r != nil {
			lg.Lgdef.Errorf("[%s] UNHANDLER ERROR [%s] %s",modName, info.openOptionsTx,r)
		}
		if conn!=nil {
			conn.Close()
		}
	}()

	buf := make([]byte, 128)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		if readCount==0 {
			tmrCheckRead.Reset(5 * time.Second)
			lg.Lgdef.Printf("[%s] Starting data reading... >>>",modName)
		}

		readCount++
		readData = readData + string(buf[:n])
	}
}

func sendDataToServer(info *SerialReaderInfo,data string) {

	info.ChannelReceiveData <- data

}

func closeConn(info *SerialReaderInfo) {
	if info.portConn!=nil {
		info.portConn.Close()
	}
}

func disposeSerial(info *SerialReaderInfo) {
	lg.Lgdef.Debugf("[%s] : == Dispose INIT === ",modName)
	closeConn(info)
	lg.Lgdef.Debugf("[%s]): === Dispose FINISH === ",modName)
}