package blobfilter

import (
	"fmt"
	"regexp"
	"strings"
)

// RegexFilter and its attributes
type nameFilter struct {
	opr bool           // true means equal to , false means not equal to
	exp *regexp.Regexp // Regex to be matched
}

func newNameFilter() attrFilter {
	return &nameFilter{}
}

func (nf *nameFilter) configure(filter string) error {
	keyLen := len(nameFilterKey)
	str := strings.Map(stringConv, filter)

	if len(str) < (keyLen + 2) {
		return fmt.Errorf("invalid name filter")
	}

	operand := str[keyLen : keyLen+2] // single char after tier (ex- tier=hot , here sinChk will be "=")
	value := ""

	if operand[0] == '=' {
		nf.opr = true
		value = str[keyLen+1:]
	} else if operand == "!=" {
		nf.opr = false
		value = str[keyLen+2:]
	} else {
		return fmt.Errorf("invalid operator %s", operand)
	}

	var err error
	nf.exp, err = regexp.Compile(value)
	if err != nil {
		return fmt.Errorf("invalid regex filter")
	}

	return nil
}

func (nf *nameFilter) isAcceptable(fileInfo *BlobAttr) bool {
	return (nf.opr == (nf.exp.MatchString((*fileInfo).Name)))
}
