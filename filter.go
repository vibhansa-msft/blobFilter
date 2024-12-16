package blobfilter

import (
	"fmt"
	"strings"
)

type BlobFilter struct {
	filters         [][]attrFilter
	parallelFilters *concurrentFilters
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
		var individualFilter []attrFilter // This array will store all filters separated by && at each index

		splitAnd := strings.Split(strings.TrimSpace(andFilters), "&&") // Split the sub filter on basis of logical AND

		for _, singleFilter := range splitAnd {
			filterName := extractName(strings.TrimSpace(singleFilter))

			if filterName == "invalid" {
				return fmt.Errorf("invalid filter, token %s", singleFilter)
			}

			// Check the filter name is a valid and supported name or not
			if constructor, exists := filterFactory[filterName]; exists {
				// Create the filter object
				filterObj := constructor()
				if filterObj == nil {
					return fmt.Errorf("failed to create filter object: %s", singleFilter)
				}

				// Configure the filter object
				err := filterObj.configure(singleFilter)
				if err != nil {
					return fmt.Errorf("failed to configure filter %s [%s]", singleFilter, err.Error())
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

// --------------------------------------------------------------------------------------

// Check a given file attributes pass the configured filter or not
func (bf *BlobFilter) IsFileAcceptable(attr *BlobAttr) bool {
	for _, filterSet := range bf.filters {
		andFilterFailed := false
		for _, individualFilter := range filterSet {
			// One filterSet is composition of AND checks
			if !individualFilter.isAcceptable(attr) {
				// If one of the filter fails then its going to be FALSE result only so
				// no need to check further filters
				andFilterFailed = true
				break
			}
		}

		// all filterSets are composition of OR checks
		if !andFilterFailed {
			// As one of the subfilter is pass then return true no need to check further filters
			return true
		}
	}

	// Nothing matched so return FALSE
	return false
}

// --------------------------------------------------------------------------------------
// Parallel Filtering logic below

// Enable parallel filtering for the application
func (bf *BlobFilter) StartParallelProcessing(concurrency int) {
	// Create work and results channels for the application
	bf.parallelFilters = &concurrentFilters{
		work:    make(chan filterKey, concurrency*2),
		results: make(chan filterResult, concurrency*2),
	}

	// Start worker threads that will process the keys
	bf.parallelFilters.wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go bf.parallelFilters.filterProcessor(bf.IsFileAcceptable)
	}
}

// Stop parallel filtering
func (bf *BlobFilter) StopParallelProcessing() {
	// Close the work channel and wait for results to flush out
	if bf.parallelFilters != nil {
		bf.parallelFilters.close()
	}
	bf.parallelFilters = nil
}

// Add one item for filtering
func (bf *BlobFilter) AddItem(key string, attr *BlobAttr) error {
	if bf.parallelFilters != nil {
		bf.parallelFilters.addWork(key, attr)
		return nil
	}
	return fmt.Errorf("parallel filtering is not enabled")
}

// Get result of the next filtered item
func (bf *BlobFilter) NextResult() (string, bool) {
	if bf.parallelFilters != nil {
		return bf.parallelFilters.getNextResult()
	}

	return "", false
}
