// Copyright 2023 Francesca <https://github.com/Francescatai>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Francescatai/microblog_GOgRPC.

package err

var ErrPostNotFound = &Errno{HTTP: 404, Code: "ResourceNotFound.PostNotFound", Message: "Post was not found."}
