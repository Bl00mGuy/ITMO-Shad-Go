//go:build !solution

package blowfish

// #cgo pkg-config: libcrypto
// #cgo CFLAGS: -Wno-deprecated-declarations
// #include <openssl/blowfish.h>
import "C"
import (
	"unsafe"
)

type Blowfish struct {
	key C.BF_KEY
}

func New(key []byte) *Blowfish {
	if len(key) == 0 || len(key) > 56 {
		panic("key length must be between 1 and 56 bytes")
	}
	bf := &Blowfish{}
	cKey := (*C.uchar)(unsafe.Pointer(&key[0]))
	C.BF_set_key(&bf.key, C.int(len(key)), cKey)
	return bf
}

func (bf *Blowfish) BlockSize() int {
	return 8
}

func (bf *Blowfish) Encrypt(dst, src []byte) {
	if len(src) != 8 || len(dst) != 8 {
		panic("input and output slices must be exactly 8 bytes")
	}
	cSrc := (*C.uchar)(unsafe.Pointer(&src[0]))
	cDst := (*C.uchar)(unsafe.Pointer(&dst[0]))
	C.BF_ecb_encrypt(cSrc, cDst, &bf.key, C.BF_ENCRYPT)
}

func (bf *Blowfish) Decrypt(dst, src []byte) {
	if len(src) != 8 || len(dst) != 8 {
		panic("input and output slices must be exactly 8 bytes")
	}
	cSrc := (*C.uchar)(unsafe.Pointer(&src[0]))
	cDst := (*C.uchar)(unsafe.Pointer(&dst[0]))
	C.BF_ecb_encrypt(cSrc, cDst, &bf.key, C.BF_DECRYPT)
}
