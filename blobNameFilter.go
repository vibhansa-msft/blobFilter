package blobfilter

import (
	"fmt"
	"regexp"
	"strings"
)

// RegexFilter and its attributes
type NameFilter struct {
	opr bool           // true means equal to , false means not equal to
	exp *regexp.Regexp // Regex to be matched
}

func newNameFilter() AttrFilter {
	return &NameFilter{}
}

func (nf *NameFilter) Configure(filter string) error {
	keyLen := len(NameFilterKey)
	str := strings.Map(StringConv, filter)

	if len(str) < (keyLen + 2) {
		return fmt.Errorf("invalid size filter")
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

func (nf *NameFilter) IsAcceptable(fileInfo *BlobAttr) bool {
	return (nf.opr == (nf.exp.MatchString((*fileInfo).Name)))
}
