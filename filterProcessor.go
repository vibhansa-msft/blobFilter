package blobfilter

import "sync"

// Key to filter the blobs
type filterKey struct {
	key  string
	attr *BlobAttr
}

// Filtering results
type filterResult struct {
	key  string
	pass bool
}

type concurrentFilters struct {
	wg      sync.WaitGroup
	work    chan filterKey
	results chan filterResult
}

func (cf *concurrentFilters) close() {
	close(cf.work)
	cf.wg.Wait()
	close(cf.results)
}

func (cf *concurrentFilters) addWork(key string, attr *BlobAttr) {
	cf.work <- filterKey{key: key, attr: attr}
}

func (cf *concurrentFilters) getNextResult() (string, bool) {
	result := <-cf.results
	return result.key, result.pass
}

func (cf *concurrentFilters) filterProcessor(filter func(attr *BlobAttr) bool) {
	for work := range cf.work {
		cf.results <- filterResult{key: work.key, pass: filter(work.attr)}
	}

	cf.wg.Done()
}
