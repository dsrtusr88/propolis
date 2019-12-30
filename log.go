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
	logthis.Info(fmt.Sprintf("%s | %s | %s", res, rule, comment), logthis.NORMAL)
}

func (l *Log) NonCriticalResult(check bool, rule, commentOK, commentKO string) {
	res := ui.BlueBold("OK")
	comment := ui.BlueBold(commentOK)
	if !check {
		res = ui.YellowBold("KO")
		comment = ui.YellowBold(commentKO)
	}
	logthis.Info(fmt.Sprintf("%s | %s | %s", res, rule, comment), logthis.NORMAL)
}

func (l *Log) NeutralResult(check bool, rule, commentOK, commentKO string) {
	res := ui.BlueBold("OK")
	comment := ui.BlueBold(commentOK)
	if !check {
		res = ui.BlueBold("KO")
		comment = ui.BlueBold(commentKO)
	}
	logthis.Info(fmt.Sprintf("%s | %s | %s", res, rule, comment), logthis.NORMAL)
}

func (l *Log) BadResult(check bool, rule, commentOK, commentKO string) {
	res := ui.YellowBold("KO")
	comment := ui.YellowBold(commentOK)
	if !check {
		res = ui.RedBold("KO")
		comment = ui.RedBold(commentKO)
	}
	logthis.Info(fmt.Sprintf("%s | %s | %s", res, rule, comment), logthis.NORMAL)
}
