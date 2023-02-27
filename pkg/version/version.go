// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package version

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/gosuri/uitable"
)

var (
	// version NO.
	GitVersion = "0.0.0-main+$Format:%h$"
	// ISO8601 format time $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	BuildDate = "1970-01-01T00:00:00Z"
	// value about Git SHA1，$(git rev-parse HEAD)
	GitCommit = "$Format:%H$"
	// Git repo state when build：clean, dirty
	GitTreeState = ""
)

// version info
type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// string format version info
func (info Info) String() string {
	if s, err := info.Text(); err == nil {
		return string(s)
	}

	return info.GitVersion
}

// json format version info
func (info Info) ToJSON() string {
	s, _ := json.Marshal(info)

	return string(s)
}

// version info -> UTF-8 format
func (info Info) Text() ([]byte, error) {
	table := uitable.New()
	table.RightAlign(0)
	table.MaxColWidth = 80
	table.Separator = " "
	table.AddRow("gitVersion:", info.GitVersion)
	table.AddRow("gitCommit:", info.GitCommit)
	table.AddRow("gitTreeState:", info.GitTreeState)
	table.AddRow("buildDate:", info.BuildDate)
	table.AddRow("goVersion:", info.GoVersion)
	table.AddRow("compiler:", info.Compiler)
	table.AddRow("platform:", info.Platform)

	return table.Bytes(), nil
}

// return version info detail
func Get() Info {
	// -ldflags setting
	return Info{
		GitVersion:   GitVersion,
		GitCommit:    GitCommit,
		GitTreeState: GitTreeState,
		BuildDate:    BuildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
