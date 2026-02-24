# Nox for Go
The official golang SDK for nox makes it easy and simple to write
nox APIs quicky and safely. As a high level wrapper all the low
level C/DLL hastle is gone away!

# How to use
Just import the library into your project as noxgo, and
call the functions!

```go
import "C"
import ...

func main() { }

//export NoxMain
func NoxMain() {
	nox := noxgo.InitNox()
	nox.CreateEndpoint("/text", func(resp *noxgo.HttpResponse, req *noxgo.HttpRequest) {
		resp.WriteText("Fortnite!")
	})
	nox.CreateEndpoint("/file", func(resp *noxgo.HttpResponse, req *noxgo.HttpRequest) {
		resp.WriteFile("./main.go")
	})
	nox.CreateEndpoint("/buff", func(resp *noxgo.HttpResponse, req *noxgo.HttpRequest) {
		buff := noxgo.NewBuffer(10)
		buff.Append([]byte{1,2,3,4,5,6,7,8,9,10})
		resp.WriteBuff(buff)
	})
	nox.CreateEndpoint("/json", func(resp *noxgo.HttpResponse, req *noxgo.HttpRequest) {
		t := &Test{ A: 5, B: "Hi!" }
		resp.WriteJson(t)
	})
}

type Test struct {
	A int `json:"a"`
	B string `json:"b"`
}
```

**NOTE: //export NoxMain is a REQUIRED HANDLE! It will be called as the entrypoint to the DLL. All setup must be done in that function.**
