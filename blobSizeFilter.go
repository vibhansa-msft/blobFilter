package blobfilter

import (
	"fmt"
	"strconv"
	"strings"
)

// SizeFilter and its attributes
type SizeFilter struct {
	opr   string // Operator
	value int64  // Value
}

func newSizeFilter() AttrFilter {
	return &SizeFilter{}
}

func (sf *SizeFilter) Configure(filter string) error {
	keyLen := len(SizeFilterKey)
	str := strings.Map(StringConv, filter)

	if len(str) < (keyLen + 2) {
		return fmt.Errorf("invalid size filter")
	}

	allowedOpts := []string{"<=", ">=", "!=", "<", ">", "="}
	for _, opt := range allowedOpts {
		// Check if filter contains an allowed operation or not
		if strings.Contains(str, opt) {
			// Operation is allowed so we need to extract the value
			splitedParts := strings.Split(str, opt)
			val, err := strconv.ParseInt(splitedParts[1], 10, 64)
			if err != nil {
				return fmt.Errorf("invalid size format")
			}

			sf.opr = opt
			sf.value = val
			return nil
		}
	}

	return fmt.Errorf("invalid size filter operator")
}

func (sf *SizeFilter) IsAcceptable(fileInfo *BlobAttr) bool {
	switch sf.opr {
	case "<=":
		return sf.value <= fileInfo.Size
	case ">=":
		return sf.value >= fileInfo.Size
	case "!=":
		return sf.value != fileInfo.Size
	case ">":
		return sf.value > fileInfo.Size
	case "<":
		return sf.value < fileInfo.Size
	case "=":
		return sf.value == fileInfo.Size
	}

	return false
}
