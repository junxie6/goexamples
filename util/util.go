package util

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func ReadInput() (string, error) {
	return bufio.NewReader(os.Stdin).ReadString('\n')
}

func WriteFile(fileName string, data []byte) error {
	return ioutil.WriteFile(fileName, data, 0644)
}

func ReadFile(fileName string) ([]byte, error) {
	return ioutil.ReadFile(fileName)
}

func ReadWebContent(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	// We must close resp.Body on all execution paths.
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("StatusCode: %s", resp.Status)
	}

	defer resp.Body.Close()

	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		return body, nil
	}
}
