package rename

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const workers = 5

func Rename(rootDir, exp, template string) error {
	regex := regexp.MustCompile(exp)

	ch := make(chan string, workers)
	done := make(chan bool)
	defer close(done)

	for i := 0; i < workers; i++ {
		go worker(ch, done, getConvertFn(regex, template))
	}

	walk(rootDir, ch)

	return nil
}

func walk(root string, ch chan<- string) {
	filepath.WalkDir(root, func(path string, d fs.DirEntry, _ error) error {
		if !d.IsDir() {
			ch <- path
		}
		return nil
	})
}

func worker(ch <-chan string, done <-chan bool, convert func(string) []byte) {
	for {
		select {
		case <-done:
			return
		case path := <-ch:
			filename := filepath.Base(path)
			if newName := convert(filename); newName != nil {
				newPath := strings.Replace(path, filename, string(newName), 1)
				if err := os.Rename(path, newPath); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
		}
	}
}

func getConvertFn(reg *regexp.Regexp, template string) func(string) []byte {
	return func(src string) []byte {
		if !reg.MatchString(src) {
			return nil
		}
 		matches := reg.FindSubmatchIndex([]byte(src))
		return reg.ExpandString([]byte{}, template, src, matches)
	}
}