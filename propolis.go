package propolis

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gitlab.com/catastrophic/assistance/logthis"
)

type Propolis struct {
	Path      string   `json:"path"`
	Checks    []*Check `json:"checks"`
	stdOutput bool
	buffer    bytes.Buffer
	Errors    int
	Warnings  int
}

func NewPropolis(path string) *Propolis {
	return &Propolis{Path: path, stdOutput: true}
}

func (p *Propolis) ToggleStdOutput(enabled bool) {
	p.stdOutput = enabled
	if !enabled {
		logthis.SetLevel(logthis.VERBOSESTEST)
		logthis.SetStdOutput(true)
		logthis.SetOutputWriter(&p.buffer)
	} else {
		logthis.SetOutputWriter(os.Stdout)
	}
}

func (p *Propolis) ParseResults() (int, int, int) {
	var numOK int
	for _, c := range p.Checks {
		switch c.Result {
		case OK:
			numOK++
		case Warning:
			p.Warnings++
		case KO:
			p.Errors++
		}
	}
	return numOK, p.Warnings, p.Errors
}

func (p *Propolis) Summary() string {
	ok, warn, ko := p.ParseResults()
	return fmt.Sprintf("%d checks OK, %d checks KO, and %d warnings.", ok, ko, warn)
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

// JSONOutput the complete log in JSON
func (p Propolis) JSONOutput() string {
	p.ParseResults()
	// marshallIndentint *p itself
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "could not generate JSON"
	}
	return string(data)
}