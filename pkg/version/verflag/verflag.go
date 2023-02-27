// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package verflag

// Package verflag defines utility functions to handle command line flags
// related to version of IAM.


import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"

	"microblog/pkg/version"
)

type versionValue int

const (
	VersionFalse versionValue = 0
	VersionTrue  versionValue = 1
	VersionRaw   versionValue = 2
)

const (
	strRawVersion   = "raw"
	versionFlagName = "version"
)

var versionFlag = Version(versionFlagName, VersionFalse, "Print version information and quit.")

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() interface{} {
	return v
}

// String method for pflag.Value interface
func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}

	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

// Set method for pflag.Value interface
func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw

		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}

	return err
}

// Type method for pflag.Value interface
func (v *versionValue) Type() string {
	return "version"
}

// VersionVar define a flag which has a specify name and use method
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	pflag.Var(p, name, usage)
	// `--version` equal to `--version=true`
	pflag.Lookup(name).NoOptDefVal = "true"
}

// VersionVar function's encapsulation 
func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)

	return p
}

// Registers this package's flags on an arbitrary FlagSet 
// such that they point to the same values ​​as the global flags
func AddFlags(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(versionFlagName))
}

// Checks if the `--version` flag was passed, if so, prints the version and exits.
func PrintAndExitIfRequested() {
	if *versionFlag == VersionRaw {
		fmt.Printf("%#v\n", version.Get())
		os.Exit(0)
	} else if *versionFlag == VersionTrue {
		fmt.Printf("%s\n", version.Get())
		os.Exit(0)
	}
}