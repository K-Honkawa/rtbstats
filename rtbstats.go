package rtbstats

import (
	"bufio"
	"encoding/json"
	"image/color"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

const (
	offsetStatus  = 0
	offsetDSPID   = 5
	offsetResTime = 13
	offsetPrice   = 21
	maskStatus    = 5
	maskDSPID     = 8
	maskResTime   = 8
	maskPrice     = 24
)

// RTBStat has RTB info of one Auction on one DSP
type RTBStat struct {
	DSPID   int `json:"dspid"`
	Status  int `json:"status"`
	Price   int `json:"price"`
	ResTime int `json:"res_time"`
}

// ToJSON return json
func (rs RTBStat) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(rs)
	return string(jsonBytes), err
}

// NewRTBStat return rtbstat
func NewRTBStat(i int) RTBStat {
	return RTBStat{
		DSPID:   (i >> offsetDSPID) % (1 << maskDSPID),
		Status:  (i >> offsetStatus) % (1 << maskStatus),
		Price:   (i >> offsetPrice) % (1 << maskPrice),
		ResTime: (i >> offsetResTime) % (1 << maskResTime),
	}
}

type stats struct {
	//DSPID int
	StatusStats map[int]int `json:"stats"`
	SUMPrice    int         `json:"price"`
}

func newStats() *stats {
	return &stats{StatusStats: map[int]int{}, SUMPrice: 0}
}

func (s *stats) stack(rs RTBStat) {
	s.StatusStats[rs.Status]++
	s.SUMPrice += rs.Price
}

func (s stats) SUMCount() int {
	return sumInt(s.StatusStats)
}

func sumInt(is map[int]int) int {
	ret := 0
	for _, n := range is {
		ret += n
	}
	return ret
}

// RTBStats is Stats map[DSPID]stats
type RTBStats struct {
	DSPStats map[int]*stats `json:"result"`
}

// NewRTBStats is
func NewRTBStats() *RTBStats {
	return &RTBStats{DSPStats: map[int]*stats{}}
}

func NewRTBStatsFromJson(jsonBytes []byte) (*RTBStats, error) {
	ret := &RTBStats{}
	err := json.Unmarshal(jsonBytes, ret)
	return ret, err
}

// ToJSON return json
func (rss RTBStats) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(rss)
	return string(jsonBytes), err
}

func (rss RTBStats) GetStatus(dspID, statID int) int {
    stats, exist := rss.DSPStats[dspID];
    if !exist {
        return 0
    }
    ret, exist  := stats.StatusStats[statID]
    if !exist {
        return 0
    }
	return ret
}

func (rss RTBStats) GetPrice(dspID int) int {
    stats, exist := rss.DSPStats[dspID];
    if !exist {
        return 0
    }
    return stats.SUMPrice
}

func (rss RTBStats) GetSumRequest(dspID int) int {
    stats, exist := rss.DSPStats[dspID];
    if !exist {
        return 0
    }
    return stats.SUMCount()
}

// Stack is
func (rss *RTBStats) Stack(statInt int) {
	stat := NewRTBStat(statInt)
	if _, exist := rss.DSPStats[stat.DSPID]; !exist {
		rss.DSPStats[stat.DSPID] = newStats()
	}
	rss.DSPStats[stat.DSPID].stack(stat)
}

type rtbStatsPlotter struct {
	rtbStatsVec []RTBStats
	plotConf    *plot.Plot
	lines       []plotLine
}

func NewRTBStatsPlotter(p *plot.Plot) *rtbStatsPlotter {
	ret := &rtbStatsPlotter{plotConf: p}
	ret.rtbStatsVec = make([]RTBStats, ret.eventSize(), ret.eventSize())
	return ret
}

type plotLine struct {
	y     func(RTBStats) float64
	color color.RGBA
	title string
}

func (pl plotLine) buildXYs(rs []RTBStats) plotter.XYs {
	ret := make(plotter.XYs, len(rs))
	for i, rss := range rs {
		ret[i].X = float64(i)
		ret[i].Y = pl.y(rss)
	}
	return ret
}

func (pl plotLine) buildLinePoints(rs []RTBStats) (*plotter.Line, *plotter.Scatter, error) {
	pLine, pScatter, err := plotter.NewLinePoints(pl.buildXYs(rs))
	if err != nil {
		return nil, nil, err
	}
	pLine.LineStyle = draw.LineStyle{Color: pl.color, Width: vg.Points(1)}
	pScatter.GlyphStyle.Color = pl.color
	return pLine, pScatter, nil
}

func (rsv rtbStatsPlotter) eventSize() int {
	if rsv.plotConf == nil {
		return 0
	}
	return int(rsv.plotConf.X.Max - rsv.plotConf.X.Min + 1)
}

// SetRTBStats set rtbstats
func (rsv *rtbStatsPlotter) SetRTBStats(rss RTBStats) {
	rsv.rtbStatsVec = append(rsv.rtbStatsVec[1:], rss)
}

// LoadFile load file
func (rsv *rtbStatsPlotter) LoadFile(path string) error {
	fp, err := os.Open(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		rs, err := NewRTBStatsFromJson(scanner.Bytes())
		if err != nil {
			return err
		}
		rsv.SetRTBStats(*rs)
	}
	return nil
}

// SetLine set line
func (rsv *rtbStatsPlotter) SetLine(y func(RTBStats) float64, c color.RGBA, t string) {
	rsv.lines = append(rsv.lines, plotLine{y: y, color: c, title: t})
}

// ClearLines clear lines
func (rsv *rtbStatsPlotter) ClearLines(y func(RTBStats) float64, c color.RGBA, t string) {
	rsv.lines = []plotLine{}
}

func (rsv rtbStatsPlotter) newPlot() (*plot.Plot, error) {
	p, err := plot.New()
	if err != nil {
		return nil, err
	}
	p.Add(plotter.NewGrid())
	if rsv.plotConf == nil {
		return p, err
	}
	p.Title = rsv.plotConf.Title
	p.BackgroundColor = rsv.plotConf.BackgroundColor
	p.BackgroundColor = rsv.plotConf.BackgroundColor
	p.X = rsv.plotConf.X
	p.Y = rsv.plotConf.Y
	p.Legend = rsv.plotConf.Legend
	return p, nil
}

// Png is make PNG
func (rsv rtbStatsPlotter) Png(path string) error {
	p, err := rsv.newPlot()
	if err != nil {
		return err
	}

	for _, pl := range rsv.lines {
		line, point, err := pl.buildLinePoints(rsv.rtbStatsVec)
		if err != nil {
			return err
		}
		p.Add(line)
		p.Add(point)
		p.Legend.Add(pl.title, line)
	}

	return p.Save(6*vg.Inch, 6*vg.Inch, path)
}
