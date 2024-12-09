package blobfilter

import (
	"fmt"
	"strings"
)

type BlobFilter struct {
	filters [][]AttrFilter
}

func NewBlobFilter() *BlobFilter {
	return &BlobFilter{}
}

// This function parses the input string and stores filters in the BlobFilter object
func (bf *BlobFilter) Configure(filterStr string) error {
	// LOGIC:
	// Parse the input string
	// Split the input string on basis of logical OR
	// For each part split by OR
	// Split the part on basis of logical AND
	// For each filter in the part
	// Get the filter name
	// Create the filter object
	// Configure the filter object
	// Append the filter object to the filter array
	// Return nil if everything went well, else return an error

	splitOr := strings.Split(filterStr, "||") // Split the string on basis of logical OR

	for _, andFilters := range splitOr {
		var individualFilter []AttrFilter // This array will store all filters separated by && at each index

		splitAnd := strings.Split(andFilters, "&&") // Split the sub filter on basis of logical AND

		for _, singleFilter := range splitAnd {
			filterName := extractName(strings.TrimSpace(singleFilter))

			if filterName == "invalid" {
				return fmt.Errorf("invalid filter, token %s", singleFilter)
			}

			// Check the filter name is a valid and supported name or not
			if constructor, exists := filterFactory[filterName]; exists {
				// Create the filter object
				filterObj := constructor()
				if filterObj != nil {
					return fmt.Errorf("failed to create filter object: %s", singleFilter)
				}

				// Configure the filter object
				err := filterObj.Configure(singleFilter)
				if err != nil {
					return fmt.Errorf("failed to configure filter: %s", singleFilter)
				}

				// Append the filter object to the filter array
				individualFilter = append(individualFilter, filterObj) // inner array (splitted by &&) is being formed

			} else {
				return fmt.Errorf("unsupported filter: %s", filterName)
			}
		}

		// Append the filter array to the BlobFilter object
		bf.filters = append(bf.filters, individualFilter) // outer array (splitted by ||) is being formed
	}

	return nil
}
