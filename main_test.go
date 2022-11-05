package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func Test_traverseDir(t *testing.T) {
	entry, _ := ioutil.ReadDir("./fixtures")
	type args struct {
		hashes     map[string]string
		duplicates map[string]string
		dupeSize   *int64
		entries    []os.FileInfo
		directory  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "testSuccessPath",
			args: args{
				hashes:     make(map[string]string),
				duplicates: make(map[string]string),
				dupeSize:   new(int64),
				entries:    make([]os.FileInfo, 0),
				directory:  ".",
			},
		},
		{
			name: "test1",
			args: args{
				hashes:     make(map[string]string),
				duplicates: make(map[string]string),
				dupeSize:   new(int64),
				entries:    entry,
				directory:  "./fixtures",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			traverseDir(tt.args.hashes, tt.args.duplicates, tt.args.dupeSize, tt.args.entries, tt.args.directory)
			t.Log("totalDuplicates", tt.args.dupeSize)
			fmt.Println(tt.name, " -->> ", *tt.args.dupeSize, tt.args.directory)
		})
	}
}
