package main

import (
	"context"
	"ty/csi/ws/ComTcpWS/Global"
	lg "ty/csi/ws/ComTcpWS/lgg"
	"ty/csi/ws/ComTcpWS/Utils"
	"github.com/jacobsa/go-serial/serial"
	"time"
	"io"
	"encoding/hex"
)

type SerialReaderInfo struct {
	OpenOptions Global.SerialOption
	openOptionsTx string

	portConn io.ReadWriteCloser
}

var modName = "SERIAL-READER"

func initSerial(ctx context.Context,info *SerialReaderInfo) {

	info.openOptionsTx = "SerialOptions: [" + Utils.PrettyPrint(info.openOptionsTx) + "]"

	lg.Lgdef.Info("SERIAL-READER:: INIT " + info.openOptionsTx)

	ticker := time.NewTicker(5 * time.Second)

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

	options := serial.OpenOptions{
		PortName:   info.OpenOptions.Portname,
		BaudRate:   uint(info.OpenOptions.Baudrate),
		DataBits:   uint(info.OpenOptions.Databits),
		StopBits:   uint(info.OpenOptions.Stopbits),
		ParityMode: intToParityMode(info.OpenOptions.ParityMode),
		InterCharacterTimeout: 0,
		MinimumReadSize:32,
	}

	// Open the port.
	port, err := serial.Open(options)
	if err != nil {
		lg.Lgdef.Errorf("%s: ERROR Open Serial Port. %v",modName, err)
		return err
	}

	if err != nil {
		lg.Lgdef.Errorf("[%s] Error opening serial port: %s",modName, err)
		return err
	} else {
		defer port.Close()
	}

	info.portConn = port

	defer func() {
		if r := recover(); r != nil {
			lg.Lgdef.Errorf("UNHANDLER ERROR [%s %s] %s",modName, info.openOptionsTx,r)
		}
		if port!=nil {
			port.Close()
		}
	}()

	for {
		buf := make([]byte, 32)
		n, err := port.Read(buf)
		if err != nil {
			if err != io.EOF {
				lg.Lgdef.Errorf("[%s] Error opening serial port: %s",modName, err)
			}
		} else {
			buf = buf[:n]
			readtext :=  hex.EncodeToString(buf)
			lg.Lgdef.Infof("%s Read data from port. '%s'",modName,readtext)
		}
	}
}

func intToParityMode(value int) serial.ParityMode {
	switch value {
	case 1: return serial.PARITY_ODD
	case 2: return serial.PARITY_EVEN
	default:
		return serial.PARITY_NONE
	}
}

func closeConn(info *SerialReaderInfo) {
	if info.portConn!=nil {
		info.portConn.Close()
	}
}

func disposeSerial(info *SerialReaderInfo) {
	lg.Lgdef.Debugf("%s(%s): == Dispose INIT === ",modName)
	closeConn(info)
	lg.Lgdef.Debugf("%s(%s): === Dispose FINISH === ",modName)
}