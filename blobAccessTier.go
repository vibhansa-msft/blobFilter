package blobfilter

import (
	"fmt"
	"strings"
)

type accessTierFilter struct {
	opr  bool   // true means equal to , false means not equal to
	tier string // allowed access tier
}

func newAccessTierFilter() attrFilter {
	return &accessTierFilter{}
}

func (atf *accessTierFilter) configure(filter string) error {
	keyLen := len(accessTierFilterKey)
	str := strings.Map(stringConv, filter)

	if len(str) < (keyLen + 2) {
		return fmt.Errorf("invalid access tier filter format")
	}

	operand := str[keyLen : keyLen+2] // single char after tier (ex- tier=hot , here sinChk will be "=")

	if operand[0] == '=' {
		atf.opr = true
		atf.tier = str[keyLen+1:]
	} else if operand == "!=" {
		atf.opr = false
		atf.tier = str[keyLen+2:]
	} else {
		return fmt.Errorf("invalid operator %s", operand)
	}

	return nil
}

func (atf *accessTierFilter) isAcceptable(fileInfo *BlobAttr) bool {
	// if opr is true then return true if tier is equal to fileInfo.Tier
	// if opr is false then return true if tier is not equal to fileInfo.Tier
	return (atf.opr == (atf.tier == strings.ToLower(fileInfo.Tier)))
}
