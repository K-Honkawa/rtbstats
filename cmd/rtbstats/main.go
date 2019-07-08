package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type element struct {
	DSPID  int `json:"dspid"`
	Status  int `json:"status"`
	Price int `json:"price"`
	ResTime int `json:"res_time"`
}

const (
	offsetStatus = 0
	offsetDSPID = 5
	offsetResTime = 13
	offsetPrice = 21
	maskStatus = 5
	maskDSPID = 8
	maskResTime = 8
	maskPrice = 24
)

func newElement(i int) element {
	return element {
		DSPID: (i >> offsetDSPID)%(1 << maskDSPID),
		Status: (i >> offsetStatus)%(1 << maskStatus),
		Price: (i >> offsetPrice)%(1 << maskPrice),
		ResTime: (i >> offsetResTime)%(1 << maskResTime),
	}
}

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		elementStrs := strings.Split(stdin.Text(), ";")
		for _,elementStr := range elementStrs {
			eleInt ,err := strconv.Atoi(elementStr) 
			if err != nil {
				fmt.Fprintf(os.Stderr, "{\"err\": \"%v\"}\n", err)
				continue
			}
			jsonBytes, err := json.Marshal(newElement(eleInt))
			if err != nil {
				fmt.Fprintf(os.Stderr, "{\"err\": \"%v\"}\n", err)
				continue
			}
			fmt.Println(string(jsonBytes))
		}
	}
}
