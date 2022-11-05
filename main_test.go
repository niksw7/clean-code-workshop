package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

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

func Test_traverseDir1(t *testing.T) {
	entry, _ := ioutil.ReadDir("./fixtures")
	type args struct {
		d         DuplicatesInformation
		entries   []os.FileInfo
		directory string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				d:         DuplicatesInformation{make(map[string]string), make(map[string]string), new(int64)},
				entries:   make([]os.FileInfo, 0),
				directory: ".",
			},
		},
		{
			name: "test2",
			args: args{
				d:         DuplicatesInformation{make(map[string]string), make(map[string]string), new(int64)},
				entries:   entry,
				directory: "./fixtures",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := traverseDir(&tt.args.d, tt.args.entries, tt.args.directory); (err != nil) != tt.wantErr {
				t.Errorf("traverseDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			fmt.Println(tt.name, " -->> ", *tt.args.d.DupeSize, tt.args.directory)
		})
	}
}
