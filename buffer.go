package noxgo

import "C"
import "unsafe"

type NoxBuffer struct {
	data []byte
	ptr unsafe.Pointer
	Length uintptr
}

func NewBuffer(size uintptr) *NoxBuffer {
	ptr := C.malloc(C.size_t(size))
	slice := unsafe.Slice((*byte)(ptr), size)

	return &NoxBuffer{
		data: slice,
		ptr: ptr,
		Length: 0,
	}
}

func (buff *NoxBuffer) Append(slice []byte) {
	n := copy(buff.data[buff.Length:], slice)
	buff.Length += uintptr(n)
}
