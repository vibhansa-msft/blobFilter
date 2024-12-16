package blobfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type nameTestSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *nameTestSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

func (suite *nameTestSuite) TestEqual() {
	atf := newNameFilter()
	err := atf.configure("name=hot*")

	filter := atf.(*nameFilter)

	suite.assert.Nil(err)
	suite.assert.True(filter.opr)
	suite.assert.NotNil(filter.exp)

	suite.assert.True(atf.isAcceptable(&BlobAttr{Name: "hot123"}))
	suite.assert.False(atf.isAcceptable(&BlobAttr{Tier: "ht123"}))
}

func (suite *nameTestSuite) TestNotEqual() {
	atf := newNameFilter()
	err := atf.configure("name!=hot*")

	filter := atf.(*nameFilter)

	suite.assert.Nil(err)
	suite.assert.False(filter.opr)
	suite.assert.NotNil(filter.exp)

	suite.assert.False(atf.isAcceptable(&BlobAttr{Name: "hot123"}))
	suite.assert.True(atf.isAcceptable(&BlobAttr{Tier: "ht123"}))
}

func (suite *nameTestSuite) TestSpaceConfig() {
	atf := newNameFilter()
	err := atf.configure(" name  =    hot*   ")

	filter := atf.(*nameFilter)

	suite.assert.Nil(err)
	suite.assert.True(filter.opr)
	suite.assert.NotNil(filter.exp)

	suite.assert.True(atf.isAcceptable(&BlobAttr{Name: "hot123"}))
	suite.assert.False(atf.isAcceptable(&BlobAttr{Tier: "ht123"}))
}

func (suite *nameTestSuite) TestInvalidConfig() {
	atf := newNameFilter()
	err := atf.configure(" name  >=   hot   ")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid operator")

	err = atf.configure(" name=  ")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid name filter")

	err = atf.configure("name=a[1-")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid regex filter")
}

func (suite *nameTestSuite) TestMultipleFilter() {
	testFileName := "mine1982.doc" // File name to be used as input to filtering api

	filters := map[string]bool{
		"m":                  true,  // contains m somewhere in the name
		"m+":                 true,  // contains m (one or more occurance) in the name
		"z":                  false, // contains z somewhere in the name
		"^m":                 true,  // name starts with m
		"^a":                 false, // name starts with a
		"c$":                 true,  // name ends with c
		"a$":                 false, // name ends with a
		"m*":                 true,  // One or more occurance of m in file name
		"^m.*":               true,  // starts with m followed by anything
		"^m.*\\.doc":         true,  // starts with m and has .doc in name
		"^m.*\\.doc$":        true,  // starts with m and ends with .doc in name
		"^*":                 true,  // * match
		"\\.":                true,  // expect a . in the file name
		"\\.doc":             true,  // expect .doc in the end of the name
		"^*z":                false, // z not at the first character of the file name
		"^*z$":               false, // z at the end of the file name
		"^a[a-z]*":           false, // starts with a and followed by any number of a-z characters
		"^a*z":               false, // starts with a and ends with z
		"^*.doc":             true,  // ends with .doc
		"^*.pdf":             false, // ends with .pdf
		"^a[a-z]*.*doc":      false, // starts with a and followed by any number of a-z characters and ends with .doc
		"^m[a-z]*.*doc":      true,  // starts with m and followed by any number of a-z characters and ends with .doc
		"^m[A-Z]*.doc":       false, // starts with m and followed by any number of A-Z characters and ends with .doc
		"^mine[0-1]\\d{3}.*": true,  // starts with mine followed by 4 digit number
	}

	for filter := range filters {
		atf := newNameFilter()
		err := atf.configure("name=" + filter)
		suite.assert.Nil(err)

		res := atf.isAcceptable(&BlobAttr{Name: testFileName})
		suite.assert.Equal(filters[filter], res)
	}
}

func TestNameFilterSuite(t *testing.T) {
	suite.Run(t, new(nameTestSuite))
}
