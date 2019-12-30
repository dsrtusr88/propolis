package main

import (
	"fmt"

	"gitlab.com/catastrophic/assistance/logthis"
	"gitlab.com/catastrophic/assistance/ui"
)

type Log struct {
	logthis.LogThis
}

func (l *Log) CriticalResult(check bool, rule, commentOK, commentKO string) {
	res := ui.BlueBold("OK")
	comment := ui.BlueBold(commentOK)
	if !check {
		res = ui.RedBold("KO")
		comment = ui.RedBold(commentKO)
	}
	l.log(res, rule, comment)
}

func (l *Log) NonCriticalResult(check bool, rule, commentOK, commentKO string) {
	res := ui.BlueBold("OK")
	comment := ui.BlueBold(commentOK)
	if !check {
		res = ui.YellowBold("KO")
		comment = ui.YellowBold(commentKO)
	}
	l.log(res, rule, comment)
}

func (l *Log) NeutralResult(check bool, rule, commentOK, commentKO string) {
	res := ui.BlueBold("OK")
	comment := ui.BlueBold(commentOK)
	if !check {
		res = ui.BlueBold("KO")
		comment = ui.BlueBold(commentKO)
	}
	l.log(res, rule, comment)
}

func (l *Log) BadResult(check bool, rule, commentOK, commentKO string) {
	res := ui.YellowBold("KO")
	comment := ui.YellowBold(commentOK)
	if !check {
		res = ui.RedBold("KO")
		comment = ui.RedBold(commentKO)
	}
	l.log(res, rule, comment)
}

func (l *Log) log(res, rule, comment string) {
	logthis.Info(fmt.Sprintf(" %2s | %-10s | %s", res, rule, comment), logthis.NORMAL)
}
