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
	rtbStats := rtbstats.NewRTBStats()
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		elementStrs := strings.Split(stdin.Text(), ";")
		for _, elementStr := range elementStrs {
			eleInt, err := strconv.Atoi(elementStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "{\"err\": \"%v\"}\n", err)
				continue
			}
			rtbStats.Stack(eleInt)
		}
	}
	jsonStr, _ := rtbStats.ToJSON()
	fmt.Println(jsonStr)
}
