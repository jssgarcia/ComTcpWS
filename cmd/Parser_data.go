package main

import (
	"strings"
	"errors"
	"fmt"
)

func main() {

	strparsed,_ := parseData("    06:10   11.6 23.1ºC\n")
	println( strparsed)
}

func parseData(data string) (string,error){

	//FORMATO DE DATOS ESPERADO
	//    06:10   11.6 23.1ºC\n
	/*
	a) Trim b) Reemplazar '   ' x ' ' c) Dividir por espacios d) Obtener pos(1) */
	formatspec := "    06:10   11.6 23.1ºC"
	newdata := strings.TrimSpace(data)
	newdata = strings.Replace(data,"   "," ",-1)
	newdata = strings.TrimSpace(newdata)
	parts := strings.Split(newdata," ")

	if (len(parts)<1){
		fmt.Errorf("[%s] PARSE DATA. '%s'.  NOT FORMAT SPEC. #%s#","__",formatspec)
		return  "",errors.New("No format spec. Format expexted #" + formatspec + "#")
	}

	return strings.TrimSpace(parts[1]),nil

}


