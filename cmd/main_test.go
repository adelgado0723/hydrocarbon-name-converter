package main

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

func Test_getInput(t *testing.T) {
	tests := []struct {
		testName  string
		want      string
		errorCase bool
		errorMsg  string
		args      []string
	}{
		{"arg of 'C' returns 'methane'", "methane", false, "", []string{"cmd", "C"}},
		{"arg of 'CCC' returns 'methane'", "propane", false, "", []string{"cmd", "CCC"}},
		{
			"invalid input returns error",
			"",
			true,
			"SMILES string contains invalid characters",
			[]string{"cmd", "1231233-2"},
		},
		{
			"empty input returns error",
			"",
			true,
			"Cannot set data from empty SMILES string",
			[]string{"cmd", ""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			// Saving the original os.Args reference
			savedArgs := os.Args
			defer func() {
				// Restoring the original os.Args and flag.CommandLine
				os.Args = savedArgs
				flag.CommandLine = flag.NewFlagSet(
					os.Args[0],
					flag.ExitOnError,
				)
			}()

			os.Args = tt.args
			got, err := handleInput()
			if (err != nil) != tt.errorCase {
				t.Errorf("Error reading SMILES input string: %v, errorCase %v ", err, tt.errorCase)
				return
			}
			if err != nil && err.Error() != tt.errorMsg {
				t.Errorf("Error comparing received error: %v, to the expected errorMsg %v ", err.Error(), tt.errorMsg)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Error comparing result: %v, does not equal the expected %v", got, tt.want)
			}
		})
	}
}
