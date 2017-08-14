package psip

import (
	"testing"
	"time"

	"github.com/mh-orange/testutil"
)

func TestTables(t *testing.T) {
	time.Local = time.UTC
	tests := []struct {
		constructor func([]byte) interface{}
		file        string
	}{
		{func(data []byte) interface{} { return &Table{data} }, "tests/table.yaml"},
		{func(data []byte) interface{} { return newSTT(data) }, "tests/stt.yaml"},
		{func(data []byte) interface{} { return newMGT(data) }, "tests/mgt.yaml"},
		{func(data []byte) interface{} { return channel(data) }, "tests/channel.yaml"},
	}

	for _, testCase := range tests {
		err := testutil.IterateTests(testCase.file, func(name string, test testutil.Test) {
			table := testCase.constructor(test.Input)
			result := testutil.Compare(test.Expected, table)
			if result.Failed() {
				t.Errorf("Test %s:%s\n%s", testCase.file, name, result)
			}
		})

		if err != nil {
			t.Errorf("Tests failed: %v", err)
		}
	}
}
