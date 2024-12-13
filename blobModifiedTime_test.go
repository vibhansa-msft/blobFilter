package blobfilter

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LMTSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *LMTSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

var timeRFC1123 time.Time
var timeStr string

func (suite *LMTSuite) TestConfigureAllOperations() {
	opts := []string{"<=", ">=", "<", ">", "="}

	for _, opt := range opts {
		atf := newModTimeFilter()
		err := atf.configure(fmt.Sprintf("modtime%s%s", opt, timeStr))

		filter := atf.(*modTimeFilter)

		suite.assert.Nil(err)
		suite.assert.Equal(opt, filter.opr)
		suite.assert.Equal(timeRFC1123, filter.value)
	}

	atf := newModTimeFilter()
	err := atf.configure("modtime@10")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid time filter operator")

	atf = newModTimeFilter()
	err = atf.configure("modtime=24/01/1982 13:00:00 IST")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid time format")

	atf = newModTimeFilter()
	err = atf.configure("modtime=")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid modified time filter")
}

func (suite *LMTSuite) TestAllOperations() {
	opts := []string{"<=", ">=", "<", ">", "="}
	times := []string{
		"Mon, 24 Jan 1982 13:00:00 IST",
		"Mon, 24 Jan 1982 13:00:00 IST",
		"Mon, 24 Jan 1982 12:00:00 IST",
		"Mon, 24 Jan 1982 14:00:00 IST",
		"Mon, 24 Jan 1982 13:00:00 IST",
	}

	for idx, opt := range opts {
		atf := newModTimeFilter()
		err1 := atf.configure(fmt.Sprintf("modtime%s%s", opt, timeStr))

		filter := atf.(*modTimeFilter)

		suite.assert.Nil(err1)
		suite.assert.Equal(opt, filter.opr)

		tempTime, err := time.Parse(time.RFC1123, strings.TrimSpace(times[idx]))
		suite.assert.Nil(err)
		suite.assert.Equal(timeRFC1123, filter.value)

		suite.assert.True(filter.isAcceptable(&BlobAttr{Mtime: tempTime}))
	}

	atf := newModTimeFilter()
	err1 := atf.configure(fmt.Sprintf("modtime<%s", timeStr))

	filter := atf.(*modTimeFilter)
	suite.assert.Nil(err1)

	tempTime, err := time.Parse(time.RFC1123, strings.TrimSpace(times[0]))
	suite.assert.Nil(err)
	suite.assert.Equal(timeRFC1123, filter.value)

	suite.assert.False(filter.isAcceptable(&BlobAttr{Mtime: tempTime}))
}

func TestLMTFilterSuite(t *testing.T) {
	var err error
	timeStr = "Mon, 24 Jan 1982 13:00:00 IST"
	timeRFC1123, err = time.Parse(time.RFC1123, strings.TrimSpace(timeStr))
	if err != nil {
		fmt.Println("Failed to parse time")
	}

	suite.Run(t, new(LMTSuite))
}
