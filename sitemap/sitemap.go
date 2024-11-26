package sitemap

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/Junior-Green/gophercises/hlp"
	"github.com/Junior-Green/gophercises/queue"
	"github.com/Junior-Green/gophercises/set"
)

type bfsNode struct {
	url   url.URL
	level int
}

type xmlUrl struct {
	Url string `xml:"loc"`
}

type siteMap struct {
	XMLName xml.Name `xml:"urlset"`
	Urls    []xmlUrl `xml:"url"`
}

type SiteMapBuilder struct {
	host url.URL
	urls *set.Set[url.URL]
}

func NewBuilder(host url.URL) *SiteMapBuilder {
	builder := &SiteMapBuilder{
		host: host,
		urls: set.NewSet[url.URL](),
	}

	return builder
}

func (b *SiteMapBuilder) BuiltSiteMap(depth int, out io.Writer) error {
	b.urls.Clear()

	b.crawl(depth, b.host)

	if err := generateXMLSiteMap(b.urls.Slice(), out); err != nil {
		return err
	}
	return nil
}

func (b *SiteMapBuilder) crawl(depth int, start url.URL) {
	var q queue.Queue[bfsNode]
	q.Push(bfsNode{start, 0})

	for q.Length() > 0 {
		n := q.Dequeue()

		if (n.level > depth && depth > -1) || b.urls.Has(n.url) {
			continue
		}
		b.urls.Add(n.url)

		html, err := fetchHTML(n.url)

		if err != nil {
			fmt.Fprintf(os.Stdout, "Error occured fetching at \"%v\". Skipping...\n", n.url)
			continue
		}

		links, err := hlp.ExtractLinks(bytes.NewReader(html))
		if err != nil {
			fmt.Fprintf(os.Stdout, "Error extracting links at \"%v\". Skipping...\n", n.url)
			continue
		}

		for _, link := range links {
			url := b.buildURLFromLink(link)
			if url == nil {
				fmt.Fprintf(os.Stdout, "Invalid href: \"%s\". Skipping...\n", link.Href)
				continue
			}
			q.Push(bfsNode{*url, n.level + 1})
		}
	}

}

func (b *SiteMapBuilder) buildURLFromLink(link hlp.Link) *url.URL {
	parsedLink, err := url.Parse(link.Href)
	if err != nil {
		return nil
	}
	if parsedLink.Host == "" {
		parsedLink.Host, parsedLink.Scheme = b.host.Host, b.host.Scheme
		return parsedLink
	} else if parsedLink.Host == b.host.Host && parsedLink.Scheme == b.host.Scheme {
		return parsedLink
	}

	return nil
}

func fetchHTML(url url.URL) ([]byte, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func generateXMLSiteMap(urls []url.URL, out io.Writer) error {
	encoder := xml.NewEncoder(out)
	defer encoder.Close()

	xmlUrls := make([]xmlUrl, 0)

	for _, url := range urls {
		xmlUrls = append(xmlUrls, xmlUrl{url.String()})
	}

	out.Write([]byte(xml.Header))
	encoder.Indent("", "    ")
	if err := encoder.Encode(siteMap{Urls: xmlUrls}); err != nil {
		return err
	}

	return nil
}
