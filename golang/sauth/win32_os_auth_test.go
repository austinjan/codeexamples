package sauth

import (
	"syscall"
	"testing"
	"unsafe"
)

var (
	modadvapi32    = syscall.NewLazyDLL("Advapi32.dll")
	procLogonUserW = modadvapi32.NewProc("LogonUserW")
)

func TestWin32Auth(t *testing.T) {
	var token syscall.Token
	username := "austin.jan@hotmail.com"
	password := "Ponbi@1102"

	r, _, err := procLogonUserW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(username))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr("."))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(password))),
		3, 0, uintptr(unsafe.Pointer(&token)))
	if r == 0 {
		t.Errorf("LogonUserW Error: %v", err)
		return
	}

	defer token.Close()
	t.Log("Authenticated!")
}
