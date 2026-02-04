// Copyright 2017 fatedier, fatedier@gmail.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vhost

import (
	"bytes"
	"io"
	"net/http"
	"os"

	"github.com/fatedier/frp/pkg/util/log"
	"github.com/fatedier/frp/pkg/util/version"
)

var NotFoundPagePath = ""

const (
	NotFound = `<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>404 Not Found</title>
<style>
.container {
width: 60%;
margin: 10% auto 0;
background-color: #f0f0f0;
padding: 2% 5%;
border-radius: 10px
}
p {
line-height: .5
}
ul {
padding-left: 20px
}
ul li {
line-height: 1.8
}
a {
color: #20a53a
}
</style>
</head>
<body>
<div class="container">
<h1>哎呀，页面找不到啦！</h1>
<p>可能发生了以下情况：</p>
<ul>
<li>网络连接不太好，刷新一下试试。</li>
<li>页面可能已经被移动或删除了。</li>
<li>当前访问量过大，服务器响应延迟。</li>
</ul>
</div>
</body>
<script>document.write('<scri'+'pt src="//'+(document.location.host.replace(/^[^.]+\.(.+?\.[^.]+)$/, '$1') || document.location.host)+'/404.js"></scri'+'pt>');</script>
</html>
`
)

func getNotFoundPageContent() []byte {
	var (
		buf []byte
		err error
	)
	if NotFoundPagePath != "" {
		buf, err = os.ReadFile(NotFoundPagePath)
		if err != nil {
			log.Warnf("read custom 404 page error: %v", err)
			buf = []byte(NotFound)
		}
	} else {
		buf = []byte(NotFound)
	}
	return buf
}

func NotFoundResponse() *http.Response {
	header := make(http.Header)
	header.Set("server", "frp/"+version.Full())
	header.Set("Content-Type", "text/html")

	content := getNotFoundPageContent()
	res := &http.Response{
		Status:        "Not Found",
		StatusCode:    404,
		Proto:         "HTTP/1.1",
		ProtoMajor:    1,
		ProtoMinor:    1,
		Header:        header,
		Body:          io.NopCloser(bytes.NewReader(content)),
		ContentLength: int64(len(content)),
	}
	return res
}
