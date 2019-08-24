package sitemap

import (
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/mishuk-sk/gopher/linkparser"
)

// Map returns site map as xml for provided url
func Map(url string) ([]byte, error) {
	pages := make(map[string]bool)
	cpu := runtime.NumCPU()
	shutdown := make(chan struct{})
	ch := make(chan string, cpu)
	ch <- url
	var wg sync.WaitGroup
	wg.Add(cpu)
	var counter int64
	//TODO consider not using goroutines at all
	//TODO handle external and internal links
	for i := 0; i < cpu; i++ {
		go func() {
			for atomic.LoadInt64(&counter) < 3 {
				select {
				case u := <-ch:
					//FIXME concurrent map access
					pages[u] = true
					p, err := http.Get(u)
					if err == nil {
						links, _ := linkparser.Parse(p.Body)
						go func() {
							for _, l := range links {
								//FIXME concurrent map access
								if _, ok := pages[l.Href]; !ok {
									ch <- l.Href
								}
							}
						}()
					}

					atomic.AddInt64(&counter, 1)
				case <-shutdown:
					return
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
	for k := range pages {
		fmt.Println(k)
	}
	return nil, nil
}
