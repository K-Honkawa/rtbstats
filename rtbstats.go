package rtbstats

import (
	"encoding/json"
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
	StatusStats map[int]int `json:"statusstats"`
	SUMPrice    int         `json:"sumprice"`
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
	DSPStats map[int]*stats `json:"stats"`
}

// NewRTBStats is
func NewRTBStats() *RTBStats {
	return &RTBStats{DSPStats: map[int]*stats{}}
}

// ToJSON return json
func (rss RTBStats) ToJSON() (string, error) {
	jsonBytes, err := json.Marshal(rss)
	return string(jsonBytes), err
}

// Stack is
func (rss *RTBStats) Stack(statInt int) {
	stat := NewRTBStat(statInt)
	if _, exist := rss.DSPStats[stat.DSPID]; !exist {
		rss.DSPStats[stat.DSPID] = newStats()
	}
	rss.DSPStats[stat.DSPID].stack(stat)
}
