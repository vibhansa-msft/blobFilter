package blobfilter

import (
	"strings"
	"time"
	"unicode"
)

type BlobAttr struct {
	Name  string            // Name of the blob
	Mtime time.Time         // Last Modified time
	Size  int64             // Size of the file
	Tier  string            // Access tier of blob
	Tags  map[string]string // Blob-tags (key value pair)
}

type AttrFilter interface {
	// Configure the filter
	Configure(filter string) error

	// Apply the filter
	IsAcceptable(fileInfo *BlobAttr) bool
}

// Create map of filters and their object creator functions
type filterConstructor func() AttrFilter

func extractName(str string) string {
	// Filter name shall be with alphabets only
	// First non alphabet character will be the end of filter name

	for i := range str {
		if !((str[i] >= 'a' && str[i] <= 'z') || (str[i] >= 'A' && str[i] <= 'Z')) {
			return strings.ToLower(str[0:i])
		}
	}

	return "invalid"
}

// Used for converting string given by user to ideal string so that it becomes easy to process
func StringConv(r rune) rune {
	if unicode.IsSpace(r) {
		return -1 // Remove space
	}

	if r >= 'A' && r <= 'Z' {
		return unicode.ToLower(r) // Convert uppercase to lowercase
	}

	return r
}

// List of allowed filters
const (
	NameFilterKey       = "name"
	FormatFilterKey     = "format"
	SizeFilterKey       = "size"
	ModTimeFilterKey    = "modtime"
	AccessTierFilterKey = "tier"
	TagFilterKey        = "tag"
)

// Factory to hold constructors of each allowed filter
var filterFactory map[string]filterConstructor

func init() {
	// Init the factory so that based on string we can create filters
	filterFactory[NameFilterKey] = newNameFilter             // Filter on basis of blob name
	filterFactory[FormatFilterKey] = newFormatFilter         // Filter on basis of blob format
	filterFactory[SizeFilterKey] = newSizeFilter             // Filter on basis of blob size
	filterFactory[ModTimeFilterKey] = newModTimeFilter       // Filter on basis of blob modification time
	filterFactory[AccessTierFilterKey] = newAccessTierFilter // Filter on basis of blob tier
	filterFactory[TagFilterKey] = newBlobTagFilter           // Filter on basis of blob tags
}
