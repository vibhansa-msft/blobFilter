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

func (cf *concurrentFilters) start(concurrency int, filter func(attr *BlobAttr) bool) {
	cf.work = make(chan filterKey, concurrency*2)
	cf.results = make(chan filterResult, concurrency*2)

	// Start worker threads that will process the keys
	cf.wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go cf.filterProcessor(filter)
	}
}

func (cf *concurrentFilters) stop() {
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
