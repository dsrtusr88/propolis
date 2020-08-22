package propolis

import (
	"fmt"

	"gitlab.com/catastrophic/assistance/ui"
)

const (
	LevelInfo = iota
	LevelWarning
	LevelCritical
	LevelAwful
	LevelTrulyAwful
)

type Level int

const (
	AppendError      = true
	DoNotAppendError = false
)

type Check struct {
	Rule          string `json:"rule"`
	ConditionOK   string `json:"-"`
	ConditionKO   string `json:"-"`
	Result        Result `json:"result"`
	Level         Level  `json:"level"`
	ResultComment string `json:"result_comment"`
	Bullet        string `json:"-"`
}

func NewCheck(rule string, level Level, ok, ko string) *Check {
	return &Check{Rule: rule, Level: level, ConditionKO: ko, ConditionOK: ok}
}

func (c *Check) EvaluateCondition(condition bool) {
	c.evaluate(condition)
	c.Log()
}

func (c *Check) EvaluateErr(err error, appendError bool) {
	c.evaluate(err == nil)
	if err != nil && appendError {
		c.ResultComment += ": " + err.Error()
	}
	c.Log()
}

func (c *Check) evaluate(condition bool) {
	switch c.Level {
	case LevelInfo:
		if condition {
			c.Result = NeutralInfo
			c.ResultComment = c.ConditionOK
			c.Bullet = NeutralString
		} else {
			c.Result = Warning
			c.ResultComment = c.ConditionKO
			c.Bullet = WarningString
		}
	case LevelWarning:
		if condition {
			c.Result = OK
			c.ResultComment = c.ConditionOK
			c.Bullet = OKString
		} else {
			c.Result = Warning
			c.ResultComment = c.ConditionKO
			c.Bullet = WarningString
		}
	case LevelCritical:
		if condition {
			c.Result = OK
			c.ResultComment = c.ConditionOK
			c.Bullet = OKString
		} else {
			c.Result = KO
			c.ResultComment = c.ConditionKO
			c.Bullet = KOString
		}
	case LevelAwful:
		if condition {
			c.Result = Warning
			c.ResultComment = c.ConditionOK
			c.Bullet = WarningString
		} else {
			c.Result = KO
			c.ResultComment = c.ConditionKO
			c.Bullet = KOString
		}
	case LevelTrulyAwful:
		if condition {
			c.Result = KO
			c.ResultComment = c.ConditionOK
			c.Bullet = KOString
		} else {
			c.Result = KO
			c.ResultComment = c.ConditionKO
			c.Bullet = KOString
		}
	}
}

func (c *Check) RawString() string {
	return fmt.Sprintf(" %2s | %-10s | %s", c.Bullet, c.Rule, c.ResultComment)
}

func (c *Check) String() string {
	var res, comment string
	switch {
	case c.Result == OK || c.Result == NeutralInfo:
		res = ui.BlueBold(c.Bullet)
		comment = ui.BlueBold(c.ResultComment)
	case c.Result == Warning:
		res = ui.YellowBold(c.Bullet)
		comment = ui.YellowBold(c.ResultComment)
	case c.Result == KO:
		res = ui.RedBold(c.Bullet)
		comment = ui.RedBold(c.ResultComment)
	}
	return fmt.Sprintf(" %2s | %-10s | %s", res, c.Rule, comment)
}

func (c *Check) Log() {
	log.log(c.Result, c.String())
}
