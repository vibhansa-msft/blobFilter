package blobfilter

import (
	"fmt"
	"path/filepath"
	"strings"
)

// formatFilter and its attributes
type formatFilter struct {
	opr bool   // true means equal to , false means not equal to
	ext string // allowed file extension
}

func newFormatFilter() attrFilter {
	return &formatFilter{}
}

func (ff *formatFilter) configure(filter string) error {
	keyLen := len(formatFilterKey)
	str := strings.Map(stringConv, filter)

	if len(str) < (keyLen + 2) {
		return fmt.Errorf("invalid format filter")
	}

	operand := str[keyLen : keyLen+2] // single char after tier (ex- tier=hot , here sinChk will be "=")

	if operand[0] == '=' {
		ff.opr = true
		ff.ext = "." + str[keyLen+1:]
	} else if operand == "!=" {
		ff.opr = false
		ff.ext = "." + str[keyLen+2:]
	} else {
		return fmt.Errorf("invalid operator %s", operand)
	}

	return nil
}
func (ff *formatFilter) isAcceptable(fileInfo *BlobAttr) bool {
	// if opr is true then return true if ext matches
	// if opr is false then return true if ext does not match
	fileExt := filepath.Ext(fileInfo.Name)
	return (ff.opr == (ff.ext == strings.ToLower(fileExt)))
}
