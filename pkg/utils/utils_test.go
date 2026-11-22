/*
Copyright 2026 The N42 Gateway Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"testing"
)

func TestConfigValue(t *testing.T) {
	testCases := map[string]struct {
		input  string
		output string
	}{
		"should match single line": {
			input:  "cmd1",
			output: "cmd1",
		},
		"should match line with spaces": {
			input:  "cmd1 cmd2",
			output: "cmd1 cmd2",
		},
		"should trim spaces": {
			input:  "\tcmd1 ",
			output: "cmd1",
		},
		"should trim post lf": {
			input:  "cmd1\ncmd2",
			output: "cmd1",
		},
		"should trim post cr": {
			input:  "cmd1\rcmd2",
			output: "cmd1",
		},
		"should trim post crlf": {
			input:  "cmd1\r\ncmd2",
			output: "cmd1",
		},
		"should trim empty first line": {
			input:  "\ncmd1\rcmd2",
			output: "",
		},
	}

	for name, test := range testCases {
		t.Run(name, func(t *testing.T) {
			input := StringLine(test.input)
			if input != test.output {
				t.Errorf("output differs - actual: %s - expected: %s", input, test.output)
			}
		})
	}
}
