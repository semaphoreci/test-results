package parsers

import (
	"bytes"
	"testing"

	"github.com/semaphoreci/test-results/pkg/fileloader"
	"github.com/semaphoreci/test-results/pkg/parser"
	"github.com/stretchr/testify/assert"
)

func TestFindParser(t *testing.T) {
	tests := []struct {
		desc    string
		name    string
		path    string
		reader  *bytes.Reader
		want    parser.Parser
		wantErr bool
	}{
		{
			desc: "finds parser automatically",
			name: "auto",
			path: "/some/path",
			reader: bytes.NewReader([]byte(`
				<?xml version="1.0"?>
					<testsuites name="foo" time="0.1234" tests="10" failures="5" errors="1">
						<testsuite>
							<testcase name="bar">
							</testcase>
							<testcase name="baz">
							</testcase>
						</testsuite>
					</testsuites>
			`)),
			want:    Generic{},
			wantErr: false,
		},
		{
			desc: "finds parser automatically",
			name: "auto",
			path: "/some/path",
			reader: bytes.NewReader([]byte(`
				<?xml version="1.0"?>
					<testsuites name="rspec" time="0.1234" tests="10" failures="5" errors="1">
						<testsuite>
							<testcase name="bar">
							</testcase>
							<testcase name="baz">
							</testcase>
						</testsuite>
					</testsuites>
			`)),
			want:    RSpec{},
			wantErr: false,
		},
		{
			desc: "finds parser automatically",
			name: "auto",
			path: "/some/path",
			reader: bytes.NewReader([]byte(`
				<?xml version="1.0"?>
					<testsuite name="rspec">
						<testcase name="bar">
						</testcase>
						<testcase name="baz">
						</testcase>
					</testsuite>
			`)),
			want:    RSpec{},
			wantErr: false,
		},
		{
			desc: "finds parser automatically",
			name: "auto",
			path: "/some/path",
			reader: bytes.NewReader([]byte(`
				<?xml version="1.0"?>
					<testsuite name="rspec">
						<testcase name="bar">
						</testcase>
						<testcase name="baz">
						</testcase>
					</testsuite>
			`)),
			want:    ExUnit{},
			wantErr: false,
		},
		{
			desc: "finds parser automatically",
			name: "auto",
			path: "/some/path",
			reader: bytes.NewReader([]byte(`
				<?xml version="1.0"?>
					<testsuite name="rspec">
						<testcase name="bar">
						</testcase>
						<testcase name="baz">
						</testcase>
					</testsuite>
			`)),
			want:    Mocha{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			fileloader.Load(tt.path, tt.reader)

			got, err := FindParser(tt.name, tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindParser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.IsType(t, tt.want, got) {
				t.Errorf("Type of FindParser() = %v, want %v", got, tt.want)
			}
		})
	}
}
