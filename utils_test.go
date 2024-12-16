package blobfilter

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type utilsSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *utilsSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

func (suite *utilsSuite) TestExtractKeyword() {
	testData := map[string]string{
		"":            "invalid",
		"a":           "invalid",
		"a=":          "invalid",
		"a=b":         "a",
		"size=10":     "size",
		"name=a b c ": "name",
		"abc@123":     "abc",
	}

	for test := range testData {
		res := extractName(test)
		suite.assert.Equal(res, testData[test])
	}
}

func (suite *utilsSuite) TestStringConv() {
	testData := map[string]string{
		"":                  "",
		"a":                 "a",
		"a=":                "a=",
		"a=b":               "a=b",
		"size=10":           "size=10",
		" size   =  10    ": "size=10",
		"name=a b c ":       "name=abc",
		"abc@123":           "abc@123",
	}

	for test := range testData {
		res := strings.Map(stringConv, test)
		suite.assert.Equal(res, testData[test])
	}
}

func TestUtilsSuite(t *testing.T) {
	suite.Run(t, new(utilsSuite))
}
