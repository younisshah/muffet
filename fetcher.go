package main

import (
	"bytes"
	"fmt"
	"net/url"
	"sync"

	"github.com/valyala/fasthttp"
	"github.com/yhat/scrape"
	"golang.org/x/net/html"
)

type client interface {
	Get([]byte, string) (int, []byte, error)
}

type fetcher struct {
	client              client
	connectionSemaphore semaphore
	cache               *sync.Map
}

func newFetcher(c int, a string) fetcher {
	cl := client(&fasthttp.Client{MaxConnsPerHost: c})

	if a != "" {
		cl = &fasthttp.HostClient{Addr: a, MaxConns: c}
	}

	return fetcher{
		cl,
		newSemaphore(c),
		&sync.Map{},
	}
}

func (f fetcher) Fetch(s string) (*page, error) {
	s, id, err := separateFragment(s)

	if err != nil {
		return nil, err
	}

	if err, ok := f.cache.Load(s); ok && err == nil {
		return nil, nil
	} else if ok {
		return nil, err.(error)
	}

	f.connectionSemaphore.Request()
	defer f.connectionSemaphore.Release()

	n, err := f.fetchHTML(s, id)
	f.cache.Store(s, err)

	if err != nil {
		return nil, err
	}

	p := newPage(s, n)
	return &p, nil
}

func (f fetcher) fetchHTML(u, id string) (*html.Node, error) {
	s, b, err := f.client.Get(nil, u)

	if err != nil {
		return nil, err
	}

	if s/100 != 2 {
		return nil, fmt.Errorf("invalid status code %v", s)
	}

	n, err := html.Parse(bytes.NewReader(b))

	if err != nil {
		return nil, err
	}

	if id != "" {
		if _, ok := scrape.Find(n, func(n *html.Node) bool {
			return scrape.Attr(n, "id") == id
		}); !ok {
			return nil, fmt.Errorf("ID #%v not found", id)
		}
	}

	return n, nil
}

func separateFragment(s string) (string, string, error) {
	u, err := url.Parse(s)

	if err != nil {
		return "", "", err
	}

	id := u.Fragment
	u.Fragment = ""

	return u.String(), id, nil
}
