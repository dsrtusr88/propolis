package propolis

import (
	"gitlab.com/catastrophic/assistance/logthis"
)

const (
	OK = iota
	Info
	NeutralInfo
	Warning
	KO
)

type Result int

type Log struct {
	logthis.LogThis
	problemsOnly bool
}

func (l *Log) log(level Result, result string) {
	if !l.problemsOnly || (level == Warning || level == KO) {
		logthis.Info(result, logthis.NORMAL)
	}
}
