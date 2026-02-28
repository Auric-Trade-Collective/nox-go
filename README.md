# Nox for Go
The official golang SDK for nox makes it easy and simple to write
nox APIs quicky and safely. As a high level wrapper all the low
level C/DLL hastle is gone away!

# How to use
Just import the library into your project as noxgo, and
call the functions!

```go
package main

import "C"
import (
	noxgo "github.com/Auric-Trade-Collective/nox-go"
)

func main() { }

//export NoxMain
func NoxMain() {
	nox := noxgo.InitNox()
	nox.CreateGet("/text", func(resp *noxgo.HttpResponse, req *noxgo.HttpRequest) {
		resp.WriteText("Get")
	})
	nox.CreatePost("/text", func(resp *noxgo.HttpResponse, req *noxgo.HttpRequest) {
		dat := &Test{ A: 1, B: "Post" }
		resp.WriteJson(dat)
	})
}

type Test struct {
	A int `json:"a"`
	B string `json:"b"`
}

```

**NOTE: //export NoxMain is a REQUIRED HANDLE! It will be called as the entrypoint to the DLL. All setup must be done in that function.**
