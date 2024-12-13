package blobfilter

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type sizeSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *sizeSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

func (suite *sizeSuite) TestAllOperations() {
	size := 10
	opts := []string{"<=", ">=", "!=", "<", ">", "="}

	for _, opt := range opts {
		atf := newSizeFilter()
		err := atf.configure(fmt.Sprintf("size%s%d", opt, size))

		filter := atf.(*sizeFilter)

		suite.assert.Nil(err)
		suite.assert.Equal(opt, filter.opr)
		suite.assert.Equal(int64(10), filter.value)
	}

	atf := newSizeFilter()
	err := atf.configure("size@10")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid size filter operator")

	atf = newSizeFilter()
	err = atf.configure("size=abcd")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid size format")
}

func (suite *sizeSuite) TestEqual() {
	atf := newSizeFilter()
	err := atf.configure("size=10")

	filter := atf.(*sizeFilter)

	suite.assert.Nil(err)
	suite.assert.Equal("=", filter.opr)
	suite.assert.Equal(int64(10), filter.value)

	suite.assert.True(atf.isAcceptable(&BlobAttr{Size: 10}))
	suite.assert.False(atf.isAcceptable(&BlobAttr{Size: 12}))
}

func (suite *sizeSuite) TestNotEqual() {
	atf := newSizeFilter()
	err := atf.configure("size!=10")

	filter := atf.(*sizeFilter)

	suite.assert.Nil(err)
	suite.assert.Equal("!=", filter.opr)
	suite.assert.Equal(int64(10), filter.value)

	suite.assert.True(atf.isAcceptable(&BlobAttr{Size: 12}))
	suite.assert.False(atf.isAcceptable(&BlobAttr{Size: 10}))
}

func (suite *sizeSuite) TestSpaceConfig() {
	atf := newSizeFilter()
	err := atf.configure(" size  =    10   ")

	filter := atf.(*sizeFilter)

	suite.assert.Nil(err)
	suite.assert.Equal("=", filter.opr)
	suite.assert.Equal(int64(10), filter.value)

	suite.assert.True(atf.isAcceptable(&BlobAttr{Size: 10}))
	suite.assert.False(atf.isAcceptable(&BlobAttr{Size: 12}))
}

func TestSizeFilterSuite(t *testing.T) {
	suite.Run(t, new(sizeSuite))
}
