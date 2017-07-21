package ts

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/mh-orange/testutil"
)

func TestReadSync(t *testing.T) {
	inputFile := "tests/reader_read_sync.yaml"
	tests, err := testutil.GetTestData(inputFile)
	if err == nil {
		for name, test := range tests {
			reader := bufio.NewReader(bytes.NewBuffer(test.Input))
			err := readSync(reader)
			if test.Err != "" && err == nil || test.Err == "" && err != nil || (test.Err != "" && test.Err != err.Error()) {
				t.Errorf("Test %s expected error \"%v\" but got \"%v\"", name, test.Err, err)
			}
		}
	} else {
		t.Errorf("Failed to read test input file \"%s\": %v", inputFile, err)
	}
}
