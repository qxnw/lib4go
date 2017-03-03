// Copyright 2015 The WebServer Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// WebServer is a micro & pluggable web framework for Go language.

// 	package main

// 	import "github.com/lunny/WebServer"

// 	type Action struct {
// 	}

// 	func (Action) Get() string {
// 	    return "Hello WebServer!"
// 	}

// 	func main() {
// 	    t := WebServer.Classic()
// 	    t.Get("/", new(Action))
// 	    t.Run()
// 	}

// Middlewares allow you easily plugin/unplugin features for your WebServer applications.

// There are already many [middlewares](https://github.com/WebServer-contrib) to simplify your work:

// - recovery - recover after panic
// - compress - Gzip & Deflate compression
// - static - Serves static files
// - logger - Log the request & inject Logger to action struct
// - param - get the router parameters
// - return - Handle the returned value smartlly
// - ctx - Inject context to action struct

// - [session](https://github.com/WebServer-contrib/session) - Session manager, with stores support:
//   * Memory - memory as a session store
//   * [Redis](https://github.com/WebServer-contrib/session-redis) - redis server as a session store
//   * [nodb](https://github.com/WebServer-contrib/session-nodb) - nodb as a session store
//   * [ledis](https://github.com/WebServer-contrib/session-ledis) - ledis server as a session store)
// - [xsrf](https://github.com/WebServer-contrib/xsrf) - Generates and validates csrf tokens
// - [binding](https://github.com/WebServer-contrib/binding) - Bind and validates forms
// - [renders](https://github.com/WebServer-contrib/renders) - Go template engine
// - [dispatch](https://github.com/WebServer-contrib/dispatch) - Multiple Application support on one server
// - [tpongo2](https://github.com/WebServer-contrib/tpongo2) - Pongo2 teamplte engine support
// - [captcha](https://github.com/WebServer-contrib/captcha) - Captcha
// - [events](https://github.com/WebServer-contrib/events) - Before and After
// - [flash](https://github.com/WebServer-contrib/flash) - Share data between requests
// - [debug](https://github.com/WebServer-contrib/debug) - Show detail debug infomaton on log

package websrv
