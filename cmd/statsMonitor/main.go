package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/K-Honkawa/rtbstats"
)

type muRTBStats struct {
	rtbStats *rtbstats.RTBStats
	mu       sync.Mutex
}

func newRTBStats() muRTBStats {
	return muRTBStats{rtbStats: rtbstats.NewRTBStats()}
}

func (murs *muRTBStats) stack(i int) {
	murs.mu.Lock()
	murs.rtbStats.Stack(i)
	murs.mu.Unlock()
}

func (murs *muRTBStats) spitRTBStats() *rtbstats.RTBStats {
	murs.mu.Lock()
	oldStats := murs.rtbStats
	murs.rtbStats = rtbstats.NewRTBStats()
	murs.mu.Unlock()
	return oldStats
}

var rtbStats = newRTBStats()

func staking() {
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		elementStrs := strings.Split(stdin.Text(), ";")
		for _, elementStr := range elementStrs {
			eleInt, err := strconv.Atoi(elementStr)
			if err != nil {
				fmt.Fprintf(os.Stderr, "{\"err\": \"%v\"}\n", err)
				continue
			}
			rtbStats.stack(eleInt)
		}
	}
}

func main() {
	go staking()
	for {
		time.Sleep(10 * time.Second)
		jsonStr, err := rtbStats.spitRTBStats().ToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "{\"err\": \"%v\"}\n", err)
		}
		fmt.Println(jsonStr)
	}
}
