package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	flags "github.com/jessevdk/go-flags"

	"github.com/K-Honkawa/rtbstats"
)

type Options struct {
	Revision bool `short:"r" long:"revision" description:"Show revision information"`
	Second   int  `short:"s" long:"second" description:"0 < secound "`
}

var (
	opts     Options
	revision = ""
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
				fmt.Fprintf(os.Stderr, "{\"err\": \"%v\",\"str\":\"%v\"}\n", err, stdin.Text())
				continue
			}
			rtbStats.stack(eleInt)
		}
	}
}

func main() {
	_, err := flags.Parse(&opts)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if opts.Revision {
		fmt.Printf("Revision:%v\n", revision)
		os.Exit(0)
	}

	go staking()
	for {
		time.Sleep(time.Duration(1*opts.Second) * time.Second)
		jsonStr, err := rtbStats.spitRTBStats().ToJSON()
		if err != nil {
			fmt.Fprintf(os.Stderr, "{\"err\": \"%v\"}\n", err)
		}
		fmt.Println(jsonStr)
	}
}
