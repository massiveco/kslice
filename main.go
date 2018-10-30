package main

import (
	"flag"
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
	"regexp"
)

var shortenKind *bool
var findLowerCase =  regexp.MustCompile("[a-z]*")

func init()  {
	shortenKind = flag.Bool("shorten", false, "Shorten the Kind used in the filename")
	flag.Parse()
}

func main() {
	reader := bufio.NewScanner(os.Stdin)
	reader.Split(splitOnDashes)

	for reader.Scan() {
		manifest := reader.Text()
		filename, _ := buildFilename(manifest)
		ioutil.WriteFile(filename, []byte(manifest), 0644)
	}
}

func buildFilename(manifest string) (string, error) {
	obj := manifestStub{}
	yaml.Unmarshal([]byte(manifest), &obj)
	kind := obj.Kind
	if *shortenKind {
		kind = findLowerCase.ReplaceAllString(kind, "")
	}

	return strings.ToLower(fmt.Sprintf("%s-%s.yaml", kind, obj.Metadata.Name)), nil
}

func splitOnDashes(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := strings.Index(string(data), "---"); i >= 0 {
		return i + 3, data[0:i], nil
	}
	if atEOF {
		return len(data), data, nil
	}
	return 0, nil, nil
}

type manifestStub struct {
	Kind     string
	Metadata manifestMetadata
}

type manifestMetadata struct {
	Name string
}
