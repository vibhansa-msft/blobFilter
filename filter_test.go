package blobfilter

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type filerSuite struct {
	suite.Suite
	assert *assert.Assertions
}

func (suite *filerSuite) SetupTest() {
	suite.assert = assert.New(suite.T())
}

func (suite *filerSuite) TestCountSubFilters() {
	type result struct {
		count int
		err   string
	}
	filters := map[string]result{
		"a":   {0, "invalid filter, token"},
		"a=":  {0, "invalid filter, token"},
		"a=b": {0, "unsupported filter"},

		"size=10":          {1, ""},
		" size  =   10   ": {1, ""},
		"size=acbd":        {0, "failed to configure filter"},

		"size < 10 && size > 5 ": {1, ""},
		"size = 10 || size = 5 ": {2, ""},

		"size > 1000 && tag=key1:val1":                                                {1, ""},
		"size > 1000 && tag=key1:val1 || size > 2000 && tag=key2:val2":                {2, ""},
		"size > 1000 && tag=key1:val1 || size > 2000 && tag=key2:val2 || tier=hot":    {3, ""},
		"size > 1000 && tag=key1:val1 || size > 2000 && tag=key2:val2 && tier != hot": {2, ""},
	}

	for filter := range filters {
		bf := BlobFilter{}

		err := bf.Configure(filter)
		if filters[filter].count == 0 {
			suite.assert.NotNil(err)
			suite.assert.Contains(err.Error(), filters[filter].err)
		} else {
			suite.assert.Nil(err)
			suite.assert.Equal(len(bf.filters), filters[filter].count)
		}
	}
}

func (suite *filerSuite) TestFilterExecution() {
	filter := "size > 1000 && tag=key1:val1 || size > 2000 && tag=key2:val2 || tier=hot || name=^mine[0-1]\\d{3}.*"

	bf := BlobFilter{}
	err := bf.Configure(filter)
	suite.assert.Nil(err)

	type filterTests struct {
		attr   BlobAttr
		result bool
	}

	// Acceptable files are
	// size > 1000 && tag=key1:val
	// size > 2000 && tag=key2:val2
	// tier=hot
	// name=^mine[0-1]\d{3}.*

	tests := []filterTests{
		{
			// Nothing matches
			attr: BlobAttr{
				Size: 1, Tags: map[string]string{"key1": "val3"}, Tier: "cold", Name: "nine1982.doc",
			},
			result: false,
		},
		{
			// for filter 1 size matches but tag does not
			attr: BlobAttr{
				Size: 1500, Tags: map[string]string{"key1": "val3"}, Tier: "cold", Name: "nine1982.doc",
			},
			result: false,
		},
		{
			// for filter 2 size matches but tag does not
			attr: BlobAttr{
				Size: 2500, Tags: map[string]string{"key1": "val3"}, Tier: "cold", Name: "nine1982.doc",
			},
			result: false,
		},
		{
			// match just based on name
			attr: BlobAttr{
				Size: 2500, Tags: map[string]string{"key1": "val3"}, Tier: "cold", Name: "mine1982.doc",
			},
			result: true,
		},
		{
			// match just based on tier
			attr: BlobAttr{
				Size: 2500, Tags: map[string]string{"key1": "val3"}, Tier: "hot", Name: "nine1982.doc",
			},
			result: true,
		},
		{
			// match based on size and tag
			attr: BlobAttr{
				Size: 2500, Tags: map[string]string{"key2": "val2"}, Tier: "cold", Name: "nine1982.doc",
			},
			result: true,
		},
		{
			// match based on size and tag
			attr: BlobAttr{
				Size: 1500, Tags: map[string]string{"key1": "val1"}, Tier: "cold", Name: "nine1982.doc",
			},
			result: true,
		},
	}

	for _, test := range tests {
		res := bf.IsAcceptable(&test.attr)
		suite.assert.Equal(res, test.result)
	}

}

func (suite *filerSuite) TestFilterExecutionParallel() {
	filter := "size > 1000 && tag=key1:val1 || size > 2000 && tag=key2:val2 || tier=hot || name=^mine[0-1]\\d{3}.*"

	bf := BlobFilter{}
	err := bf.Configure(filter)
	suite.assert.Nil(err)

	type filterTests struct {
		attr   BlobAttr
		result bool
	}

	// Acceptable files are
	// size > 1000 && tag=key1:val
	// size > 2000 && tag=key2:val2
	// tier=hot
	// name=^mine[0-1]\d{3}.*

	tests := []filterTests{
		{
			// Nothing matches
			attr: BlobAttr{
				Size: 1, Tags: map[string]string{"key1": "val3"}, Tier: "cold", Name: "nine1982.doc",
			},
			result: false,
		},
		{
			// for filter 1 size matches but tag does not
			attr: BlobAttr{
				Size: 1500, Tags: map[string]string{"key1": "val3"}, Tier: "cold", Name: "nine1982.doc",
			},
			result: false,
		},
		{
			// for filter 2 size matches but tag does not
			attr: BlobAttr{
				Size: 2500, Tags: map[string]string{"key1": "val3"}, Tier: "cold", Name: "nine1982.doc",
			},
			result: false,
		},
		{
			// match just based on name
			attr: BlobAttr{
				Size: 2500, Tags: map[string]string{"key1": "val3"}, Tier: "cold", Name: "mine1982.doc",
			},
			result: true,
		},
		{
			// match just based on tier
			attr: BlobAttr{
				Size: 2500, Tags: map[string]string{"key1": "val3"}, Tier: "hot", Name: "nine1982.doc",
			},
			result: true,
		},
		{
			// match based on size and tag
			attr: BlobAttr{
				Size: 2500, Tags: map[string]string{"key2": "val2"}, Tier: "cold", Name: "nine1982.doc",
			},
			result: true,
		},
		{
			// match based on size and tag
			attr: BlobAttr{
				Size: 1500, Tags: map[string]string{"key1": "val1"}, Tier: "cold", Name: "nine1982.doc",
			},
			result: true,
		},
	}

	// Start parallel processing
	bf.EnableAsyncFilter(2)

	// Enqueue all items
	for idx, test := range tests {
		go bf.AddItem(strconv.FormatInt(int64(idx), 10), &test.attr)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	count := 0

	go func() {
		defer wg.Done()
		for ; count < len(tests); count++ {
			idStr, res := bf.NextResult()
			id, err := strconv.ParseInt(idStr, 10, 32)
			suite.assert.Nil(err)
			suite.assert.Equal(tests[int(id)].result, res)
		}
	}()

	wg.Wait()
	bf.TerminateAsyncFilter()
	suite.assert.Equal(count, len(tests))
}

func TestFilterSuite(t *testing.T) {
	suite.Run(t, new(filerSuite))
}
