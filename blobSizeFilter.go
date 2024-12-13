package blobfilter

import (
	"fmt"
	"strconv"
	"strings"
)

// sizeFilter and its attributes
type sizeFilter struct {
	opr   string // Operator
	value int64  // Value
}

func newSizeFilter() attrFilter {
	return &sizeFilter{}
}

func (sf *sizeFilter) configure(filter string) error {
	keyLen := len(sizeFilterKey)
	str := strings.Map(stringConv, filter)

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

func (sf *sizeFilter) isAcceptable(fileInfo *BlobAttr) bool {
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
