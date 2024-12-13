package blobfilter

import (
	"fmt"
	"strings"
)

type blobTagFilter struct {
	opr   bool   // true means equal to , false means not equal to
	key   string // Blob tag key
	value string // Blob tag value
}

func newBlobTagFilter() attrFilter {
	return &blobTagFilter{}
}

func (btf *blobTagFilter) configure(filter string) error {
	keyLen := len(tagFilterKey)
	str := strings.Map(stringConv, filter)

	if len(str) < (keyLen + 2) {
		return fmt.Errorf("invalid blob tag filter format")
	}

	operand := str[keyLen : keyLen+2]
	if operand[0] == '=' {
		btf.opr = true
	} else if operand == "!=" {
		btf.opr = false
	} else {
		return fmt.Errorf("invalid operator %s", operand)
	}

	values := strings.Split(filter, "=")
	if len(values) == 2 {
		tagFilter := strings.Split(values[1], ":")
		if len(tagFilter) == 2 {
			btf.key = tagFilter[0]
			btf.value = tagFilter[1]
		} else {
			return fmt.Errorf("invalid tag filter format")
		}
	} else {
		return fmt.Errorf("invalid tag filter format")
	}

	return nil
}

func (btf *blobTagFilter) isAcceptable(fileInfo *BlobAttr) bool {
	// Search the configured tag in fileinfo tags
	// If tag is found then match the value
	if val, ok := fileInfo.Tags[btf.key]; ok {
		return (btf.opr == (btf.value == strings.ToLower(val)))
	}

	return false
}
