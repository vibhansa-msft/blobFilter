package blobfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type formatFilterSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *formatFilterSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

func (suite *formatFilterSuite) TestEqual() {
	atf := newFormatFilter()
	err := atf.configure("format=pdf")

	filter := atf.(*formatFilter)

	suite.assert.Nil(err)
	suite.assert.True(filter.opr)
	suite.assert.Equal(".pdf", filter.ext)

	suite.assert.True(atf.isAcceptable(&BlobAttr{Name: "abcd.pdf"}))
	suite.assert.False(atf.isAcceptable(&BlobAttr{Name: "abcd.doc"}))
}

func (suite *formatFilterSuite) TestNotEqual() {
	atf := newFormatFilter()
	err := atf.configure("format!=pdf")

	filter := atf.(*formatFilter)

	suite.assert.Nil(err)
	suite.assert.False(filter.opr)
	suite.assert.Equal(".pdf", filter.ext)

	suite.assert.False(atf.isAcceptable(&BlobAttr{Name: "abcd.pdf"}))
	suite.assert.True(atf.isAcceptable(&BlobAttr{Name: "abcd.doc"}))
}

func (suite *formatFilterSuite) TestSpaceConfig() {
	atf := newFormatFilter()
	err := atf.configure(" format  =    pdf   ")

	filter := atf.(*formatFilter)

	suite.assert.Nil(err)
	suite.assert.True(filter.opr)
	suite.assert.Equal(".pdf", filter.ext)

	suite.assert.True(atf.isAcceptable(&BlobAttr{Name: "abcd.pdf"}))
	suite.assert.False(atf.isAcceptable(&BlobAttr{Name: "abcd.doc"}))
}

func (suite *formatFilterSuite) TestInvalidConfig() {
	atf := newFormatFilter()
	err := atf.configure(" format  >=   doc   ")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid operator")

	err = atf.configure(" format=  ")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid format filter")

}

func TestFormatFilterSuite(t *testing.T) {
	suite.Run(t, new(formatFilterSuite))
}
