package blobfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type tagFilterSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *tagFilterSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

var tagValues map[string]string

func (suite *tagFilterSuite) TestEqual() {
	btf := newBlobTagFilter()
	err := btf.configure("tag=key1:val1")

	filter := btf.(*blobTagFilter)

	suite.assert.Nil(err)
	suite.assert.True(filter.opr)
	suite.assert.Equal("key1", filter.key)
	suite.assert.Equal("val1", filter.value)

	suite.assert.True(btf.isAcceptable(&BlobAttr{Tags: tagValues}))
}

func (suite *tagFilterSuite) TestNotEqual() {
	btf := newBlobTagFilter()
	err := btf.configure("tag!=key1:val1")

	filter := btf.(*blobTagFilter)

	suite.assert.Nil(err)
	suite.assert.False(filter.opr)
	suite.assert.Equal("key1", filter.key)
	suite.assert.Equal("val1", filter.value)

	suite.assert.False(btf.isAcceptable(&BlobAttr{Tags: tagValues}))
}

func (suite *tagFilterSuite) TestSpaceConfig() {
	btf := newBlobTagFilter()
	err := btf.configure(" tag  =    key1   :  val2   ")

	filter := btf.(*blobTagFilter)

	suite.assert.Nil(err)
	suite.assert.True(filter.opr)

	suite.assert.False(btf.isAcceptable(&BlobAttr{Tags: tagValues}))
}

func (suite *tagFilterSuite) TestInvalidConfig() {
	btf := newBlobTagFilter()
	err := btf.configure(" tag  >=   key1 : val2   ")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid operator")

	err = btf.configure(" tag=  ")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid blob tag filter format")

	err = btf.configure(" tag=key1  ")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid tag filter format")

}

func TestTagFilterSuite(t *testing.T) {
	tagValues = make(map[string]string)
	tagValues["key1"] = "val1"
	tagValues["key2"] = "val2"

	suite.Run(t, new(tagFilterSuite))
}
