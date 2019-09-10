# rtbstats

## get

`go get github.com/K-Honkawa/rtbstats/cmd/rtbstats`

`go get github.com/K-Honkawa/rtbstats/cmd/statsmonitor`

`go get github.com/K-Honkawa/rtbstats/cmd/statser`

## exec

### rtbstats
`echo "${stat1};${stat2}...." | rtbstats`

### statsmonitor
`tail -f $file | $(2rtbstats) | statsmonitor`

### statser
`cat -f $file | $(2rtbstats) | statser`
