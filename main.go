package main

import (
	"bufio"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewScanner(os.Stdin)
	reader.Split(splitOnDashes)

	for reader.Scan() {
		manifest := reader.Text()
		kind, name, _ := extractDetails(manifest)
		ioutil.WriteFile(fmt.Sprintf("%s-%s.yaml", kind, name), []byte(manifest), 0644)
	}
}

func extractDetails(manifest string) (string, string, error) {
	obj := ManifestStub{}
	yaml.Unmarshal([]byte(strings.ToLower(manifest)), &obj)

	return obj.Kind, obj.Metadata.Name, nil
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

type ManifestStub struct {
	Kind     string
	Metadata ManifestMetadata
}

type ManifestMetadata struct {
	Name string
}
