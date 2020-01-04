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

func (l *Log) CriticalResult(check bool, rule, commentOK, commentKO string) Result {
	res := ui.BlueBold("OK")
	comment := ui.BlueBold(commentOK)
	if !check {
		res = ui.RedBold("KO")
		comment = ui.RedBold(commentKO)
		l.log(res, rule, comment)
		return KO
	}
	l.log(res, rule, comment)
	return OK
}

func (l *Log) NonCriticalResult(check bool, rule, commentOK, commentKO string) Result {
	res := ui.BlueBold("OK")
	comment := ui.BlueBold(commentOK)
	if !check {
		res = ui.YellowBold("KO")
		comment = ui.YellowBold(commentKO)
		l.log(res, rule, comment)
		return Warning
	}
	l.log(res, rule, comment)
	return OK
}

func (l *Log) NeutralResult(check bool, rule, commentOK, commentKO string) Result {
	res := ui.BlueBold("OK")
	comment := ui.BlueBold(commentOK)
	if !check {
		res = ui.BlueBold("KO")
		comment = ui.BlueBold(commentKO)
	}
	l.log(res, rule, comment)
	return Info
}

func (l *Log) BadResultInfo(check bool, rule, commentOK, commentKO string) Result {
	res := ui.YellowBold("KO")
	comment := ui.YellowBold(commentOK)
	if !check {
		res = ui.RedBold("KO")
		comment = ui.RedBold(commentKO)
	}
	l.log(res, rule, comment)
	return Info
}

func (l *Log) log(res, rule, comment string) {
	logthis.Info(fmt.Sprintf(" %2s | %-10s | %s", res, rule, comment), logthis.NORMAL)
}
