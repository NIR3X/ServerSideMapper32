package rawwrapper

import "unsafe"

// RawWrapper represents a wrapper for raw data.
type RawWrapper struct {
	Data []uint8
}

// NewRawWrapper creates a new instance of RawWrapper with the given data.
// It returns a pointer to the created RawWrapper.
func NewRawWrapper(data []uint8) *RawWrapper {
	return &RawWrapper{Data: data}
}

// At returns the unsafe.Pointer to the element at the specified position in the RawWrapper's data.
// The pos parameter represents the index of the element in the data slice.
// It is important to note that this method returns an unsafe.Pointer, which should be used with caution.
func (self *RawWrapper) At(pos uintptr) unsafe.Pointer {
	return unsafe.Pointer(&self.Data[pos])
}
