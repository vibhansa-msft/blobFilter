package blobfilter

import (
	"fmt"
	"strings"
	"time"
)

// modTimeFilter and its attributes
type ModTimeFilter struct {
	opr   string    // Operator
	value time.Time // Value
}

func newModTimeFilter() AttrFilter {
	return &ModTimeFilter{}
}

func (mtf *ModTimeFilter) Configure(filter string) error {
	keyLen := len(ModTimeFilterKey)
	str := strings.Map(StringConv, filter)

	if len(str) < (keyLen + 2) {
		return fmt.Errorf("invalid modified time filter")
	}

	allowedOpts := []string{"<=", ">=", "<", ">", "="}
	for _, opt := range allowedOpts {
		// Check if filter contains an allowed operation or not
		if strings.Contains(str, opt) {
			// Operation is allowed so we need to extract the value
			splitedParts := strings.Split(str, opt)
			timeRFC1123str := strings.TrimSpace(splitedParts[1])
			timeRFC1123, err := time.Parse(time.RFC1123, timeRFC1123str)
			if err != nil {
				return fmt.Errorf("invalid time format")
			}

			mtf.opr = opt
			mtf.value = timeRFC1123
			return nil
		}
	}

	return fmt.Errorf("invalid time filter operator")
}

func (mtf *ModTimeFilter) IsAcceptable(fileInfo *BlobAttr) bool {
	switch mtf.opr {
	case "<=":
		return fileInfo.Mtime.Before(mtf.value) || fileInfo.Mtime.Equal(mtf.value)
	case ">=":
		return fileInfo.Mtime.After(mtf.value) || fileInfo.Mtime.Equal(mtf.value)
	case ">":
		return fileInfo.Mtime.After(mtf.value)
	case "<":
		return fileInfo.Mtime.Before(mtf.value)
	case "=":
		return fileInfo.Mtime.Equal(mtf.value)
	}

	return false
}
