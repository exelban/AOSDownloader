package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Node struct {
	URI   string
	Name  string
	Nodes *[]Node
}

func parseURI(uri string, path string) (*Node, int, error) {
	log.Printf("%s - %s", uri, path)

	node := &Node{
		URI:   path,
		Nodes: &[]Node{},
	}

	b, err := fetchPage(uri)
	if err != nil {
		log.Printf("[ERROR] fetch page: %v", err)
		return nil, 0, err
	}

	l, f, err := parsePage(b)
	if err != nil {
		log.Printf("[ERROR] getLinks: %v", err)
		return nil, 0, err
	}
	log.Printf("[DEBUG] %s has %d folders and %d files", path, len(f), len(l))

	for _, link := range l {
		*node.Nodes = append(*node.Nodes, Node{
			URI:   fmt.Sprintf("%s/%s", uri, link),
			Name:  link,
			Nodes: nil,
		})
	}

	for _, folder := range f {
		folder = strings.TrimSuffix(folder, "/")
		log.Printf("[DEBUG] going to loop with: %s", folder)
		links, n, err := parseURI(fmt.Sprintf("%s/%s", uri, folder), folder)
		if err != nil {
			return nil, n, err
		}
		*node.Nodes = append(*node.Nodes, *links)
	}

	return node, len(l), nil
}

func fetchPage(link string) ([]byte, error) {
	resp, err := http.Get(link)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = errors.New(fmt.Sprintf("GET %s - wrong status code: %d (%s)", link, resp.StatusCode, http.StatusText(resp.StatusCode)))
		return nil, err
	}

	log.Printf("[DEBUG] page %s downloaded", link)

	return ioutil.ReadAll(resp.Body)
}

func parsePage(b []byte) ([]string, []string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		log.Print("[DEBUG] failed to parse document by goquery")
		return nil, nil, err
	}

	links := []string{}
	folders := []string{}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		text := s.Text()
		if len(s.ParentFiltered("td").Nodes) == 0 || text == "" || text == "Parent Directory" {
			return
		}

		if text[len(text)-1] == 47 {
			folders = append(folders, text)
		} else {
			links = append(links, text)
		}
	})

	return links, folders, nil
}
