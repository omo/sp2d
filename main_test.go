package main

import (
	"log"
	"net/http"
	"net/url"
	"testing"
)

func TestHello(t *testing.T) {
}

func makeTestRequest(host string, path string) *http.Request {
	u, e := url.Parse(path)
	if e != nil {
		log.Fatal(e)
	}

	return &http.Request{
		Method: "GET",
		Proto:  "http",
		URL:    u,
		Host:   host,
	}
}

func makeTestMapper() *URLMapper {
	return &URLMapper{
		Active: MustParse("http://active/"),
		Frontend:     MustParse("http://front/"),
		LastStore: MustParse("http://last.aws/"),
		ArchiveStore: MustParse("http://archive.aws/"),
	}
}

func TestMapFromTDiaryRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/?date=20091231")
	Expect(mapper.GetMapping(req).Front.String(), "http://front/bn/2009/12/31/", t)
	ExpectTrue(mapper.GetMapping(req).Stored == nil, req.URL.String(), t)
}

func TestMapppingShouldBeServedWithLiveDomain(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/foo/")
	Expect(mapper.GetMapping(req).Front.String(), "http://front/foo/", t)
	Expect(mapper.GetMapping(req).Stored.String(), "http://last.aws/foo/", t)
}

func TestMapppingShouldBeServedWithStagedDomain(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("s2p.flakiness.es", "/foo/")
	Expect(mapper.GetMapping(req).Front.String(), "http://front/foo/", t)
	Expect(mapper.GetMapping(req).Stored.String(), "http://last.aws/foo/", t)
}

func TestMapppingShouldBeRedirectedForNonLiveDomain(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("old.front", "/foo/")
	Expect(mapper.GetMapping(req).Front.String(), "http://front/foo/", t)
	ExpectTrue(mapper.GetMapping(req).Stored == nil, req.URL.String(), t)
}

func TestMapFromTDiaryAtomRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/index.rdf")
	Expect(mapper.GetMapping(req).Stored.String(), "http://active/index.xml", t)
}

func TestMapLastAtomRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/atom.xml")
	Expect(mapper.GetMapping(req).Stored.String(), "http://active/index.xml", t)
}

func TestMapFromOctopressRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/b/2014/07/05/life-of-touch/")
	Expect(mapper.GetMapping(req).Stored.String(), "http://last.aws/b/2014/07/05/life-of-touch/", t)
}
func TestMapFromOctopressImageRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "http://front/images/line-tile.png?1383981792")
	Expect(mapper.GetMapping(req).Front.String(), "http://front/images/line-tile.png?1383981792", t)
	Expect(mapper.GetMapping(req).Stored.String(), "http://last.aws/images/line-tile.png?1383981792", t)
}

func TestMapFromBacknumberRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/bn/2011/11/04/")
	Expect(mapper.GetMapping(req).Stored.String(), "http://archive.aws/bn/2011/11/04/", t)
}

func TestMapFromBacknumberAssetRequest(t *testing.T) {
	mapper := makeTestMapper()
	req := makeTestRequest("front", "/bn/stylesheets/style.css")
	Expect(mapper.GetMapping(req).Stored.String(), "http://archive.aws/bn/stylesheets/style.css", t)
}

func TestShouldServeDirectly(t *testing.T) {
	s := &DirectServer{}
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/")), "/", t)
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/foo/bar/")), "/", t)
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/atom.xml")), "atom.xml", t)
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/index.rdf")), "index.rdf", t)
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/index.html")), "index.html", t)
	ExpectTrue(s.ShouldServe(MustParse("http://example.com/hoge.css")), "css", t)
	ExpectTrue(!s.ShouldServe(nil), "nil", t)
	ExpectTrue(!s.ShouldServe(MustParse("http://example.com/fuga.jpeg")), "jpeg", t)
}
