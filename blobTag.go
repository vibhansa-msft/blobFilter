package blobfilter

import (
	"fmt"
	"strings"
)

type BlobTagFilter struct {
	key   string // Blob tag key
	value string // Blob tag value
}

func newBlobTagFilter() AttrFilter {
	return &BlobTagFilter{}
}

func (btf *BlobTagFilter) Configure(filter string) error {
	keyLen := len(TagFilterKey)
	str := strings.Map(StringConv, filter)

	if len(str) < (keyLen + 2) {
		return fmt.Errorf("invalid blob tag filter format")
	}

	opr := str[keyLen : keyLen+1]
	if !(opr == "=") {
		return fmt.Errorf("invalid operator %s", opr)
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

func (btf *BlobTagFilter) IsAcceptable(fileInfo *BlobAttr) bool {
	// Search the configured tag in fileinfo tags
	// If tag is found then match the value
	if val, ok := fileInfo.Tags[btf.key]; ok {
		return (btf.value == strings.ToLower(val))
	}

	return false
}
