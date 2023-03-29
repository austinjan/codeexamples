package sauth

import (
	"crypto/rand"
	"encoding/base32"
	"net/url"
	"os"
	"syscall"
	"testing"
	"unsafe"

	qr "rsc.io/qr"
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

// Test2FA tests 2FA machenism
func Test2FA(t *testing.T) {
	// Generate a random secret
	secret := make([]byte, 10)
	_, err := rand.Read(secret)
	if err != nil {
		t.Fatal(err)
	}

	// base32 encode use 32 characters (A-Z and 2-7) it is case insensitive, base32
	// base32 has lower encoding efficiency than base64, encoded text will be 40% longer than base64
	secretBase32 := base32.StdEncoding.EncodeToString(secret)

	account := "austin.jan@gmail.com"
	issuer := "austin"

	// Deal with request

	URL, err := url.Parse("otpauth://totp")
	if err != nil {
		t.Fatal(err)
	}

	URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(account)
	params := url.Values{}
	params.Add("secret", secretBase32)
	params.Add("issuer", issuer)

	URL.RawQuery = params.Encode()
	t.Logf("URL: %v", URL.String())

	//generate QR code
	code, err := qr.Encode(URL.String(), qr.Q)
	if err != nil {
		t.Fatal(err)
	}

	qrfilename := "./qr.png"
	err = os.WriteFile(qrfilename, code.PNG(), 0600)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("QR code saved to %v", qrfilename)

	// otpc := &dgoogauth.OTPConfig{
	// 	Secret:      secretBase32,
	// 	WindowSize:  3,
	// 	HotpCounter: 0,
	// 	// UTC:         true,
	// }
	// for {
	// 	fmt.Printf("Please enter the token value (or q to quit): ")

	// 	var token string
	// 	fmt.Scanln(&token)
	// 	if token == "q" {
	// 		break
	// 	}

	// 	val, err := otpc.Authenticate(token)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		continue
	// 	}

	// 	if !val {
	// 		fmt.Println("Sorry, Not Authenticated")
	// 		continue
	// 	}

	// 	fmt.Println("Authenticated!")
	// }
}
