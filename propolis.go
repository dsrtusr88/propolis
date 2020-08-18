package propolis

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type Propolis struct {
	Path           string   `json:"path"`
	Checks         []*Check `json:"checks"`
	output         OutputType
	CriticalErrors int
}

func NewPropolis(path string, out OutputType) *Propolis {
	return &Propolis{Path: path, output: out}
}

func (p *Propolis) AddChecks(c ...*Check) {
	p.Checks = append(p.Checks, c...)
}

func (p *Propolis) Summary() string {
	var numOK, numWarning int
	for _, c := range p.Checks {
		switch c.Result {
		case OK:
			numOK++
		case Warning:
			numWarning++
		case KO:
			p.CriticalErrors++
		}
	}
	return fmt.Sprintf("%d checks OK, %d checks KO, and %d warnings.", numOK, p.CriticalErrors, numWarning)
}

func (p *Propolis) ConditionCheck(level Level, rule, OKString, KOString string, condition bool) {
	check := NewCheck(rule, level, OKString, KOString)
	check.EvaluateCondition(condition)
	p.Checks = append(p.Checks, check)
}

func (p *Propolis) ErrorCheck(level Level, rule, OKString, KOString string, err error, appendError bool) {
	check := NewCheck(rule, level, OKString, KOString)
	check.EvaluateErr(err, appendError)
	p.Checks = append(p.Checks, check)
}

func (p *Propolis) ListErrors() string {
	var errors []string
	for _, c := range p.Checks {
		if c.Result == KO {
			errors = append(errors, c.ResultComment)
		}
	}
	return strings.Join(errors, " | ")
}

func (p *Propolis) ListWarnings() string {
	var warnings []string
	for _, c := range p.Checks {
		if c.Result == Warning {
			warnings = append(warnings, c.ResultComment)
		}
	}
	return strings.Join(warnings, " | ")
}

// Output the complete log
func (p *Propolis) Output() string {
	var output string
	for _, c := range p.Checks {
		// no color support, use case is writing log files
		output += c.RawString() + "\n"
	}
	return output
}

func (p *Propolis) SaveOuput(dir, version string) error {
	// TODO check dir
	outputFile := filepath.Join(dir, fmt.Sprintf("propolis_%s.log", version))
	return ioutil.WriteFile(outputFile, []byte(p.Output()), 0600)
}
