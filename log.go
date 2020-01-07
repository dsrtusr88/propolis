package main

import (
	"fmt"

	"gitlab.com/catastrophic/assistance/logthis"
	"gitlab.com/catastrophic/assistance/ui"
)

const (
	OK = iota
	Warning
	KO
	Info
	NeutralInfo
)

type Result int

type Results struct {
	ok      int
	warning int
	ko      int
}

func (r *Results) Add(res Result) {
	switch res {
	case OK:
		r.ok++
	case Warning:
		r.warning++
	case KO:
		r.ko++
	}
}

func (r *Results) String() string {
	return fmt.Sprintf("%d checks OK, %d checks KO, and %d warnings.", r.ok, r.ko, r.warning)
}

type Log struct {
	logthis.LogThis
}

func (l *Log) Critical(check bool, rule, commentOK, commentKO string) Result {
	if check {
		l.log(OK, OKString, rule, commentOK)
		return OK
	}
	l.log(KO, KOString, rule, commentKO)
	return KO
}

func (l *Log) NonCritical(check bool, rule, commentOK, commentKO string) Result {
	if check {
		l.log(OK, OKString, rule, commentOK)
		return OK
	}
	l.log(Warning, WarningString, rule, commentKO)
	return Warning
}

func (l *Log) Info(check bool, rule, commentOK, commentKO string) Result {
	if check {
		l.log(NeutralInfo, NeutralString, rule, commentOK)
	} else {
		l.log(NeutralInfo, WarningString, rule, commentKO)
	}
	return NeutralInfo
}

func (l *Log) BadResult(check bool, rule, commentOK, commentKO string) Result {
	if check {
		l.log(Warning, WarningString, rule, commentOK)
	} else {
		l.log(KO, KOString, rule, commentKO)
	}
	return Info
}

func (l *Log) log(level Result, res, rule, comment string) {
	switch {
	case level == OK || level == NeutralInfo:
		res = ui.BlueBold(res)
		comment = ui.BlueBold(comment)
	case level == Warning:
		res = ui.YellowBold(res)
		comment = ui.YellowBold(comment)
	case level == KO:
		res = ui.RedBold(res)
		comment = ui.RedBold(comment)
	}
	logthis.Info(fmt.Sprintf(" %2s | %-10s | %s", res, rule, comment), logthis.NORMAL)
}
