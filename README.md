# rtbstats

## get

`go get github.com/K-Honkawa/rtbstats/cmd/rtbstats`
`go get github.com/K-Honkawa/rtbstats/cmd/statsmonitor`

## exec

### rtbstats
`echo "${stat1};${stat2}...." | rtbstats`

### statsmonitor
`tail -f $file | $(2rtbstats) | statsmonitor`
