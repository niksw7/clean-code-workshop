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
			t.Log("totalDuplicates", *tt.args.dupeSize)
			fmt.Println(tt.name, " -->> ", *tt.args.dupeSize, tt.args.directory)
		})
	}
}

func Test_toReadableSize(t *testing.T) {
	type args struct {
		nbytes int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1000 Bytes",
			args: args{
				nbytes: 1000,
			},
			want: "1000 B",
		},
		{
			name: "1000 KB",
			args: args{
				nbytes: 1000000,
			},
			want: "1000 KB",
		},
		{
			name: "1000 MB",
			args: args{
				nbytes: 1000000000,
			},
			want: "1000 MB",
		},
		{
			name: "1000 GB",
			args: args{
				nbytes: 1000000000000,
			},
			want: "1000 GB",
		},
		{
			name: "1000 TB",
			args: args{
				nbytes: 1000000000000000,
			},
			want: "1000 TB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toReadableSize(tt.args.nbytes); got != tt.want {
				t.Errorf("toReadableSize() = %v, want %v", got, tt.want)
			}
		})
	}
}
