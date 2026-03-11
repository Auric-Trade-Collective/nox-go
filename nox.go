package noxgo

/*
#cgo LDFLAGS: -ldl
#include "webapi.h"
#include <stdint.h>
*/
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path/filepath"
	"unsafe"
)

// var Construct func()
var nox *Nox = nil

type Nox struct {
	Endpoints map[string]map[string]func(*HttpResponse, *HttpRequest)
}


func InitNox() *Nox {
	if nox != nil {
		return nil
	}

	nx := &Nox {
		Endpoints: make(map[string]map[string]func(*HttpResponse, *HttpRequest)), 
	}

	nox = nx
	return nx
}

func (nx *Nox) CreateGet(path string, fn func(*HttpResponse, *HttpRequest)) {
	endp, ok := nx.Endpoints[path]
	if !ok {
		nx.Endpoints[path] = make(map[string]func(*HttpResponse, *HttpRequest))
		endp = nx.Endpoints[path]
	}
	_ = endp
	nx.Endpoints[path][http.MethodGet] = fn
}

func (nx *Nox) CreatePost(path string, fn func(*HttpResponse, *HttpRequest)) {
	endp, ok := nx.Endpoints[path]
	if !ok {
		nx.Endpoints[path] = make(map[string]func(*HttpResponse, *HttpRequest))
		endp = nx.Endpoints[path]
	}
	_ = endp
	nx.Endpoints[path][http.MethodPost] = fn
}

func (nx *Nox) CreatePut(path string, fn func(*HttpResponse, *HttpRequest)) {
	endp, ok := nx.Endpoints[path]
	if !ok {
		nx.Endpoints[path] = make(map[string]func(*HttpResponse, *HttpRequest))
		endp = nx.Endpoints[path]
	}
	_ = endp
	nx.Endpoints[path][http.MethodPut] = fn
}

func (nx *Nox) CreateDelete(path string, fn func(*HttpResponse, *HttpRequest)) {
	endp, ok := nx.Endpoints[path]
	if !ok {
		nx.Endpoints[path] = make(map[string]func(*HttpResponse, *HttpRequest))
		endp = nx.Endpoints[path]
	}
	_ = endp
	nx.Endpoints[path][http.MethodDelete] = fn
}

type HttpRequest struct {
	handle unsafe.Pointer
	Endpoint string
	Method string
}

func (req *HttpRequest) ReadBody(buff []byte) (int, error) {
	ptr := (*C.uint8_t)(unsafe.Pointer(&buff[0]))
	
	cBytesRead := C.ReadBody((*C.HttpRequest)(req.handle), ptr, C.size_t(len(buff)))
	goBytesRead := int(cBytesRead)

	if goBytesRead <= 0 {
		return 0, errors.New("Could not read bytes into buffer")
	}

	return goBytesRead, nil
}

func (req *HttpRequest) ReadAsJson(target any) error {
	wholeBuff := []byte{}
	for {
		buff := [255]byte{}
		read, err := req.ReadBody(buff[:])
		if err != nil || read <= 0 {
			break
		}

		wholeBuff = append(wholeBuff, buff[:read]...)
	}

	err := json.Unmarshal(wholeBuff, target)
	if err != nil {
		return err
	}

	return nil
}

type HttpResponse struct {
	handle unsafe.Pointer
}

func (resp *HttpResponse) GetHeader(key string, index int) string {
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))

	var outStr *C.char
	defer C.free(unsafe.Pointer(outStr))
	success := C.TryGetResponseHeader((*C.HttpResponse)(resp.handle), cKey, C.size_t(index), &outStr)

	if int(success) != 0 {
		goStr := C.GoString(outStr)
		return goStr
	}

	return ""
}

func (resp *HttpResponse) SetHeader(key string, val string) error {
	cKey := C.CString(key)
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cKey))
	defer C.free(unsafe.Pointer(cVal))

	success := C.TrySetResponseHeader((*C.HttpResponse)(resp.handle), cKey, cVal, C.int(0))
	if success == 0 {
		return errors.New("Could not set header " + key)
	}

	return nil
}

func (resp *HttpResponse) WriteBuff(bytes *NoxBuffer, contentType string) {
	count := len(bytes.data)
	if count == 0 && bytes.Length == 0 {
		return
	}

	ccType := C.CString(contentType)
	// defer C.free(unsafe.Pointer(ccType))

	buff := C.NoxBuffer((*C.uint8_t)(bytes.ptr), C.size_t(bytes.Length), ccType)
	C.WriteMove((*C.HttpResponse)(resp.handle), buff)
}

func (resp *HttpResponse) WriteText(text string) {
	cstr := C.CString(text)
	defer C.free(unsafe.Pointer(cstr))

	C.WriteText((*C.HttpResponse)(resp.handle), cstr, C.int(len(text)))
}

func (resp *HttpResponse) WriteFile(path string) {
	abs, err := filepath.Abs(path)
	if err != nil {
		fmt.Println("nox-go: Couldn't filed file specified!")
	}

	str := C.CString(abs)
	defer C.free(unsafe.Pointer(str))

	dat := C.NoxFile(str)
	C.WriteFile((*C.HttpResponse)(resp.handle), dat);
	C.FreeData(dat)
}

func (resp *HttpResponse) WriteJson(dat any) {
	json, err := json.Marshal(dat)
	if err != nil {
		fmt.Println("nox-go: Couldn't serialize JSON")
	}

	buff := NewBuffer(uintptr(len(json)))
	buff.Append(json)

	resp.WriteBuff(buff, "application/json")
}

//export CreateNoxApi
func CreateNoxApi(endp *C.NoxEndpointCollection) {
	C.NoxMain()

	if nox == nil {
		panic("nox-go: Could not load API, was nil")
	}

	for k, v := range nox.Endpoints {
		cstr := C.CString(k)

		for k2, _ := range v {
			var method int
			switch k2 {
			case http.MethodGet: method = 0
			case http.MethodPost: method = 1
			case http.MethodPut: method = 2
			case http.MethodDelete: method = 3
			}

			C.InvokeEndp(C.int(method), endp, cstr, (*[0]byte)(unsafe.Pointer(C.GetInvokeGo())))
		}

		C.free(unsafe.Pointer(cstr))
	}
}

//export EndpointHandler
func EndpointHandler(resp *C.HttpResponse, req *C.HttpRequest) {
	goPath := C.GoString(req.endpoint)

	mReq := &HttpRequest{
		handle: unsafe.Pointer(req),
		Endpoint: goPath,
		Method: C.GoString(req.method),
	}
	mResp := &HttpResponse{
		handle: unsafe.Pointer(resp),	
	}

	if mp, ok := nox.Endpoints[goPath]; ok {
		if fn, ok := mp[mReq.Method]; ok {
			fn(mResp, mReq)
		}
	}
}
