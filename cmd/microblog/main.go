// Copyright 2023 Innkeeper Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package main

import (
	"fmt"
	"os"

	"microblog/internal/microblog"
)

func main() {
	fmt.Println("hello World")

	command := microblog.NewMicroBlogCommand()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
