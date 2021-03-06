package rtbstats

import (
	"image/color"
	"testing"

	"gonum.org/v1/plot"
)

// {"dspid":212,"status":31,"price":546,"res_time":11}
const testInt = 1145141919

func TestRTBStats_Stack(t *testing.T) {
	var rss = NewRTBStats()

	const loopNum = 3
	for i := 0; i < loopNum; i++ {
		rss.Stack(testInt)
	}

	rsJSON, err := NewRTBStat(testInt).ToJSON()
	if err != nil {
		t.Error("error:", err)
	}

	rssJSON, err := rss.ToJSON()
	if err != nil {
		t.Error("error:", err)
	}

	t.Log("loopNum:", loopNum, "\ntestStats:", rsJSON, "Result", rssJSON)

	if len(rss.DSPStats) != 1 {
		t.Error("difference DSPStats length")
	}
	if _, exist := rss.DSPStats[212]; !exist {
		t.Fatal("found 212 DSP number")
	}
	if rss.DSPStats[212].SUMCount() != (loopNum) {
		t.Error("difference SUMCount")
	}
	if rss.DSPStats[212].SUMPrice != (546 * loopNum) {
		t.Error("difference SUMCount")
	}
	if len(rss.DSPStats[212].StatusStats) != 1 {
		t.Error("difference StateuStats length")
	}
	if _, exist := rss.DSPStats[212].StatusStats[31]; !exist {
		t.Fatal("found 31 status")
	}
	if rss.DSPStats[212].StatusStats[31] != 3 {
		t.Error("difference StateuStats[31] num")
	}
}

func TestRtbStatsVector_Png(t *testing.T) {
	plotConf, err := plot.New()
	if err != nil {
		t.Fatal(err)
	}
	plotConf.X.Min = 0
	plotConf.X.Max = 10
	plotConf.Y.Min = 0
	plotConf.Y.Max = 10

	rsv := NewRTBStatsPlotter(plotConf)
	rsv.SetLine(func(RTBStats) float64 { return 5 }, color.RGBA{R: 255, G: 0, B: 0, A: 255}, "test")
	if err := rsv.Png("./test.png"); err != nil {
		t.Fatal(err)
	}
}
