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
	"gitlab.com/catastrophic/assistance/music"
)

type Propolis struct {
	Path         string   `json:"path"`
	Checks       []*Check `json:"checks"`
	release      *music.Release
	stdOutput    bool
	problemsOnly bool
	buffer       bytes.Buffer
	Passed       int
	Errors       int
	Warnings     int
}

func NewPropolis(path string, release *music.Release, problemsOnly bool) *Propolis {
	return &Propolis{Path: path, release: release, stdOutput: true, problemsOnly: problemsOnly}
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

func (p *Propolis) ParseResults() {
	p.Passed, p.Warnings, p.Errors = 0, 0, 0
	for _, c := range p.Checks {
		switch c.Result {
		case OK:
			p.Passed++
		case Warning:
			p.Warnings++
		case KO:
			p.Errors++
		}
	}
}

func (p *Propolis) Summary() string {
	p.ParseResults()
	return fmt.Sprintf("%d checks OK, %d checks KO, and %d warnings.", p.Passed, p.Errors, p.Warnings)
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

// Output the complete log.
func (p *Propolis) Output() string {
	var output string
	for _, c := range p.Checks {
		if p.problemsOnly && (c.Result == OK || c.Result == Info || c.Result == NeutralInfo) {
			continue
		}
		// no color support, use case is writing log files
		output += c.RawString() + "\n"
	}
	return output
}

// Tags of the flacs.
func (p *Propolis) Tags() string {
	return p.release.GetRawTags()
}

func (p *Propolis) SaveOuput(dir, version string) error {
	// TODO check dir
	outputFile := filepath.Join(dir, fmt.Sprintf("propolis_%s.log", version))
	if err := ioutil.WriteFile(outputFile, []byte(p.Output()), 0600); err != nil {
		return err
	}
	tagsOutputFile := filepath.Join(dir, "tags.txt")
	return ioutil.WriteFile(tagsOutputFile, []byte(p.release.GetRawTags()), 0600)
}

// JSONOutput the complete log in JSON.
func (p Propolis) JSONOutput() string {
	p.ParseResults()
	// marshallIndentint *p itself
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return "could not generate JSON"
	}
	return string(data)
}
