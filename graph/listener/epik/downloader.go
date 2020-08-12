package epik

import (
	"context"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"

	"github.com/epik-protocol/gateway/clog"
	"golang.org/x/sync/semaphore"
)

type downloader struct {
	mu     sync.Mutex
	wg     sync.WaitGroup
	sem    *semaphore.Weighted
	result map[string][]byte
	failed int
}

func newDownloader() *downloader {
	return &downloader{
		sem:    semaphore.NewWeighted(int64(runtime.NumCPU()) * 2),
		result: make(map[string][]byte),
	}
}

func (d *downloader) download(ctx context.Context, key, url string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.result[key]; ok {
		return nil
	}
	if err := d.sem.Acquire(ctx, 1); err != nil {
		clog.Errorf("failed to acquire semaphore error: %v", err)
		return err
	}
	d.wg.Add(1)
	go func() {
		defer d.sem.Release(1)
		defer d.wg.Done()

		resp, err := http.Get(url)
		if err != nil {
			clog.Errorf("failed to download %s(%s), error: %v", key, url, err)
			d.failed++
			return
		}
		defer resp.Body.Close()
		if resp.StatusCode == http.StatusOK {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				clog.Errorf("failed to read response %s(%s), error: %v", key, url, err)
				d.failed++
				return
			}
			body, err = preprocess(body)
			if err != nil {
				clog.Errorf("failed to preprocess data from %s(%s), error: %v", key, url, err)
				d.failed++
				return
			}
			d.mu.Lock()
			d.result[key] = body
			d.mu.Unlock()
		} else {
			d.failed++
		}
	}()
	return nil
}

func (d *downloader) wait() {
	d.wg.Wait()
}
