package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/K-Honkawa/rtbstats"
)

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		elementStrs := strings.Split(stdin.Text(), ";")
		for _, elementStr := range elementStrs {
			eleInt, err := strconv.Atoi(elementStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "{\"err\": \"%v\"}\n", err)
				continue
			}
			jsonBytes, err := rtbstats.NewRTBStat(eleInt).ToJSON()
			if err != nil {
				fmt.Fprintf(os.Stderr, "{\"err\": \"%v\"}\n", err)
				continue
			}
			fmt.Println(string(jsonBytes))
		}
	}
}
