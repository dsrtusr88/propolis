package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlacUtils(t *testing.T) {
	fmt.Println("+ Testing Flac/Utils...")
	check := assert.New(t)

	var err error
	config, err = NewConfig("../../test/config.yaml")
	check.Nil(err)

	t1 := "NAO - Woman (2020) - WEB FLAC|2940820|7wj612|https://privatebin.url/?30d86620ab3ceb12#DY8PLjbWjVrSKpWme9Dh1QvVQDSX6tZdYuER58fv92Em|35 checks OK, 0 checks KO, and 1 warnings."
	path, idOrErr, parsed, err := ParseNodeMessage("sender", t1)
	check.Nil(err)
	check.Equal("NAO - Woman (2020) - WEB FLAC", path)
	check.Equal("2940820", idOrErr)
	fmt.Println(parsed)

	t2 := "NAO - Woman (2020) - WEB FLAC|Error: Tracks seem incorrectly organized: track number is not a number"
	path, idOrErr, parsed, err = ParseNodeMessage("sender", t2)
	check.Nil(err)
	check.Equal("NAO - Woman (2020) - WEB FLAC", path)
	check.Equal("Error: Tracks seem incorrectly organized: track number is not a number", idOrErr)
	fmt.Println(parsed)
}
