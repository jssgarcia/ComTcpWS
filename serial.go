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
	ReadInterval int //Intervalo de tiempo (ms) que revisa si ha llegado nuevos datos

	openConfig *rs.Config
	portConn *rs.Port
	chanSyncroOP chan int
}

var modName = "SERIAL-READER"

func initSerial(ctx context.Context,info *SerialReaderInfo) {

	info.openOptionsTx = "SerialOptions: [" + Utils.PrettyPrint(info.OpenOptions) + "]"

	lg.Lgdef.Info("SERIAL-READER:: INIT " + info.openOptionsTx)

	info.chanSyncroOP = make(chan int,0)

	ticker := time.NewTicker(5 * time.Second)
	//ticker.C <- time.Now()

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
		Name:        info.OpenOptions.Portname,
		Baud:        int(info.OpenOptions.Baudrate),
		Size:        byte(info.OpenOptions.Databits),
		StopBits:    rs.StopBits(info.OpenOptions.Stopbits),
		Parity:      rs.Parity(info.OpenOptions.ParityMode),
		ReadTimeout: time.Duration(info.OpenOptions.ReadTimeout) * time.Millisecond,
	}

	info.openConfig = cnfg;
	readData := ""
	readCount := 0
	closePortAndOpen := false

	if info.ReadInterval == 0 {
		info.ReadInterval = 1000
	}

	tmrCheckRead := time.NewTimer(time.Duration(info.ReadInterval))
	tmrCheckRead.Stop()

	//Check if Data received for Port
	go func() {
		for {
			select {
			case <-tmrCheckRead.C:
				if readData == "" {
					continue
				}

				lg.Lgdef.Printf("[%s] DATA READED: Bytes:%n Dato:%s", modName, readCount, readData)
				tmrCheckRead.Stop()

				sendDataToServer(info, readData)
				closePortAndOpen=true

				//Asignamos
				readCount = 0
				readData = ""

				break
			}
		}
	}()

	// Open the port.
	conn, err := rs.OpenPort(cnfg)
	if err != nil {
		lg.Lgdef.Errorf("%s: ERROR Open Serial Port. %v", modName, err)
		return err
	}

	lg.Lgdef.Printf("[%s] PUERTO '%s' ABIERTO CON EXITO.", modName, cnfg.Name)

	info.portConn = conn

	defer func() {
		if r := recover(); r != nil {
			lg.Lgdef.Errorf("[%s] UNHANDLER ERROR [%s] %s", modName, info.openOptionsTx, r)
		}
		if conn != nil {
			conn.Close()
		}
	}()

	buf := make([]byte, 128)

	lg.Lgdef.Infof("[%s.readSerialPort] Waiting for read data ....", modName)

	for {
		if (closePortAndOpen) {
			newconn,err := closeAndOpen(info);
			if (err != nil) {
				return nil
			}
			closePortAndOpen=false
			lg.Lgdef.Infof("[%s.readSerialPort] Waiting for read data ....", modName)

			conn = newconn
		}

		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		if (readCount == 0 && n > 0) {
			//tmrCheckRead.Reset(time.Duration(info.ReadCheckInterval))
			tmrCheckRead.Reset(5 * time.Second)
			lg.Lgdef.Printf("[%s] Starting data reading... >>>", modName)
		}

		readCount = readCount + n
		readData = readData + string(buf[:n])
		//lg.Lgdef.Printf("     [%s] Bytes readed %n",modName,readCount)
		//lg.Lgdef.Printf("  >> [%s] ckuck-data read: 0x %s %s",modName,buf[:n],string(buf[:n]),buf[:n])
	}

	time.Sleep(1 * time.Hour)

	return nil
}

func closeAndOpen(info *SerialReaderInfo) (*rs.Port,error) {

	err:=closeConn(info)
	if (err!=nil) {
		return nil,err;
	}

	info.portConn = nil;

	conn, err := rs.OpenPort(info.openConfig)
	if err != nil {
		lg.Lgdef.Errorf("%s.closeAndOpen: ERROR Open Serial Port. %v",modName, err)
		return nil, err
	}

	lg.Lgdef.Printf("[%s] PUERTO '%s' ABIERTO CON EXITO.",modName,info.openConfig.Name)

	info.portConn = conn;

	return conn,nil
}

func sendDataToServer(info *SerialReaderInfo,data string) {

	info.ChannelReceiveData <- data

}

func closeConn(info *SerialReaderInfo) (error) {
	if info.portConn!=nil {
		err := info.portConn.Close()
		if (err!=nil) {
			lg.Lgdef.Errorf("[%s.closeConn] ERROR CERRAR PUERTO '%s'",modName,err)
			return err
		}

		lg.Lgdef.Printf("[%s.closeConn] PUERTO CERRADO con exito'",modName)
	}

	return nil
}

func disposeSerial(info *SerialReaderInfo) {
	lg.Lgdef.Debugf("[%s] : == Dispose INIT === ",modName)
	closeConn(info)
	lg.Lgdef.Debugf("[%s]): === Dispose FINISH === ",modName)
}