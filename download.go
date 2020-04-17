package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
)

func download(node *Node, folder string) error {
	if node.Nodes != nil {
		if _, err := os.Stat(folder); err == nil {
			if err = os.RemoveAll(folder); err != nil {
				return err
			}
		}

		if err := os.Mkdir(folder, 0777); err != nil {
			return err
		}

		log.Printf("[DEBUG] folder %s created", folder)

		wg := sync.WaitGroup{}
		wg.Add(len(*node.Nodes))
		errList := []error{}
		for _, n := range *node.Nodes {
			go func(n Node) {
				f := path.Clean(fmt.Sprintf("%s/%s", folder, n.URI))
				if n.Nodes == nil {
					f = folder
				}
				if err := download(&n, f); err != nil {
					log.Printf("[ERROR] download node %s", n.URI)
					errList = append(errList, err)
				}
				wg.Done()
			}(n)
		}
		wg.Wait()

		return nil
	}

	f := path.Clean(fmt.Sprintf("%s/%s", folder, node.Name))
	log.Printf("[DEBUG] downloading file %s to %s", node.URI, f)
	if err := downloadFile(f, node.URI); err != nil {
		return err
	}

	return nil
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		b, _ := ioutil.ReadAll(resp.Body)
		return errors.New(fmt.Sprintf("unable to fetch remote file: (%d: %s). Remote error: \n%v", resp.StatusCode, http.StatusText(resp.StatusCode), string(b)))
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	_, err = io.Copy(out, resp.Body)
	return err
}
