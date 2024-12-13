package blobfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type accessTierSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *accessTierSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

func (suite *accessTierSuite) TestEqual() {
	atf := newAccessTierFilter()
	err := atf.configure("tier=hot")

	filter := atf.(*accessTierFilter)

	suite.assert.Nil(err)
	suite.assert.True(filter.opr)
	suite.assert.Equal("hot", filter.tier)

	suite.assert.True(atf.isAcceptable(&BlobAttr{Tier: "hot"}))
	suite.assert.False(atf.isAcceptable(&BlobAttr{Tier: "cool"}))
}

func (suite *accessTierSuite) TestNotEqual() {
	atf := newAccessTierFilter()
	err := atf.configure("tier!=hot")

	filter := atf.(*accessTierFilter)

	suite.assert.Nil(err)
	suite.assert.False(filter.opr)
	suite.assert.Equal("hot", filter.tier)

	suite.assert.False(atf.isAcceptable(&BlobAttr{Tier: "hot"}))
	suite.assert.True(atf.isAcceptable(&BlobAttr{Tier: "cool"}))
}

func (suite *accessTierSuite) TestSpaceConfig() {
	atf := newAccessTierFilter()
	err := atf.configure(" tier  =    hot   ")

	filter := atf.(*accessTierFilter)

	suite.assert.Nil(err)
	suite.assert.True(filter.opr)
	suite.assert.Equal("hot", filter.tier)

	suite.assert.True(atf.isAcceptable(&BlobAttr{Tier: "hot"}))
	suite.assert.False(atf.isAcceptable(&BlobAttr{Tier: "cool"}))
}

func (suite *accessTierSuite) TestInvalidConfig() {
	atf := newAccessTierFilter()
	err := atf.configure(" tier  >=   hot   ")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid operator")

	err = atf.configure(" tier=  ")
	suite.assert.NotNil(err)
	suite.assert.Contains(err.Error(), "invalid access tier filter format")

}

func TestAccessTierFilterSuite(t *testing.T) {
	suite.Run(t, new(accessTierSuite))
}
