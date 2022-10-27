package nebula

/*
#cgo LDFLAGS: -L ./lib/ -lcrypto
#cgo LDFLAGS: -L ./lib/ -lssl
#cgo CFLAGS: -I ./include/
#include "openssl/evp.h"
#include "openssl/aes.h"
*/
import "C"

import (
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"unsafe"

	//"fmt"
	//"unsafe"

	"github.com/cipherloc/noise"
)

type endianness interface {
	PutUint64(b []byte, v uint64)
}

var noiseEndianness endianness = binary.BigEndian

type NebulaCipherState struct {
	c noise.Cipher
	//k [32]byte
	//n uint64
}

func NewNebulaCipherState(s *noise.CipherState) *NebulaCipherState {
	return &NebulaCipherState{c: s.Cipher()}
}

type (
	Ctx *C.EVP_CIPHER_CTX
)

func Get_Ctx() Ctx {
	return C.EVP_CIPHER_CTX_new()
}

// EncryptDanger encrypts and authenticates a given payload.
//
// out is a destination slice to hold the output of the EncryptDanger operation.
// - ad is additional data, which will be authenticated and appended to out, but not encrypted.
// - plaintext is encrypted, authenticated and appended to out.
// - n is a nonce value which must never be re-used with this key.
// - nb is a buffer used for temporary storage in the implementation of this call, which should
// be re-used by callers to minimize garbage collection.
func (s *NebulaCipherState) EncryptDanger(out, ad, plaintext []byte, n uint64, nb []byte, ct string) ([]byte, error) {
	if s != nil {
		// TODO: Is this okay now that we have made messageCounter atomic?
		// Alternative may be to split the counter space into ranges
		//if n <= s.n {
		//	return nil, errors.New("CRITICAL: a duplicate counter value was used")
		//}
		//s.n = n
		nb[0] = 0
		nb[1] = 0
		nb[2] = 0
		nb[3] = 0
		noiseEndianness.PutUint64(nb[4:], n)

		c := s.c.(cipher.AEAD)

		// override encrypt for AESGCMFIPS
		if ct == "AESGCMFIPS" {

			//fmt.Printf("CIPHER WITH STACK: %s\n", c.name)
			var tempLength int = 0
			var output []byte = make([]byte, 8096)
			var outputLength int = 0
			var inputArray []byte = []byte(plaintext)
			var inputLength int = len(inputArray)
			var key [32]byte = s.c.Key()

			//fmt.Println("********* ENCRYPT *********")

			pInput := (*C.uchar)(unsafe.Pointer(&inputArray[0]))
			// fmt.Printf("Original: %s\n", string(inputArray))
			// fmt.Printf("Original Length: %d\n", len(inputArray))

			pKey := (*C.uchar)(unsafe.Pointer(C.CString(string(key[:]))))
			defer C.free((unsafe.Pointer)(pKey))
			// fmt.Printf("UChar* key = %s", key)

			pIv := (*C.uchar)(unsafe.Pointer(C.CString(string(key[0:15]))))
			defer C.free((unsafe.Pointer)(pIv))
			//fmt.Printf("UChar* iv =  %\n", pIv)

			var ctx Ctx = Get_Ctx()
			//fmt.Printf("Context made\n")

			C.EVP_EncryptInit_ex(ctx, C.EVP_aes_128_gcm(), nil, pKey, pIv)
			//fmt.Printf("Encrypt Init\n")

			_ = C.EVP_EncryptUpdate(ctx, (*C.uchar)(&output[0]), (*C.int)(unsafe.Pointer(&outputLength)), pInput, (C.int)(inputLength))
			//fmt.Printf("Update Value: %d\n", value)

			_ = C.EVP_EncryptFinal_ex(ctx, (*C.uchar)(&output[outputLength]), (*C.int)(unsafe.Pointer(&tempLength)))
			//fmt.Printf("Final Value: %d\n", value)

			// fmt.Printf("TempLength: %d\nTotalLength: %d\n", tempLength, outputLength+tempLength)
			C.EVP_CIPHER_CTX_free(ctx)
			//fmt.Printf("Freed\n")

			output = output[0 : outputLength+tempLength]

			//fmt.Printf("FIPSTEXT:\n%s\n", string(output))

			ciphertext := c.Seal(out, nb, output, ad)

			//fmt.Printf("CIPHERTEXT:\n%s\n", string(ciphertext))

			//fmt.Printf("ENCRYPTION: %s\n", c.name)

			return ciphertext, nil
		}

		// out = s.c.Encrypt(out, n, ad, plaintext)
		out = c.Seal(out, nb, plaintext, ad)
		//l.Debugf("Encryption: outlen: %d, nonce: %d, ad: %s, plainlen %d", len(out), n, ad, len(plaintext))
		return out, nil
		//}
	} else {
		return nil, errors.New("no cipher state available to encrypt")
	}
}

func (s *NebulaCipherState) DecryptDanger(out, ad, ciphertext []byte, n uint64, nb []byte, ct string) ([]byte, error) {
	if s != nil {
		nb[0] = 0
		nb[1] = 0
		nb[2] = 0
		nb[3] = 0
		noiseEndianness.PutUint64(nb[4:], n)

		c := s.c.(cipher.AEAD)

		ctext, err := c.Open(out, nb, ciphertext, ad)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return ctext, err
		}

		if ct == "AESGCMFIPS" {
			var inputLength int = len(ctext)
			var tempLength int = 0
			var output []byte = make([]byte, 8096)
			var outputLength int = 0
			var key [32]byte = s.c.Key()

			// TODO: Need error detection
			// fmt.Println("********* DECRYPT *********")

			pInput := (*C.uchar)(unsafe.Pointer(&ctext[0]))

			pKey := (*C.uchar)(unsafe.Pointer(C.CString(string(key[:]))))
			defer C.free((unsafe.Pointer)(pKey))
			//fmt.Printf("UChar* key = %v", pKey)

			pIv := (*C.uchar)(unsafe.Pointer(C.CString(string(key[0:15]))))
			defer C.free((unsafe.Pointer)(pIv))
			//fmt.Printf("UChar* iv =  %v\n", pIv)

			var ctx Ctx = Get_Ctx()
			//fmt.Printf("Context made\n")

			C.EVP_DecryptInit_ex(ctx, C.EVP_aes_128_gcm(), nil, pKey, pIv)
			//fmt.Printf("Decrypt Init\n")

			_ = C.EVP_DecryptUpdate(ctx, (*C.uchar)(&output[0]), (*C.int)(unsafe.Pointer(&outputLength)), pInput, (C.int)(inputLength))
			// fmt.Printf("Input Length = %d\nOutput Length = %d\n", inputLength, outputLength)
			//fmt.Printf("Update Value: %d\n", value)

			_ = C.EVP_DecryptFinal_ex(ctx, (*C.uchar)(&output[outputLength]), (*C.int)(unsafe.Pointer(&tempLength)))
			//fmt.Printf("Final Value: %d\n", value)

			// fmt.Printf("TempLength: %d\nTotalLength: %d\n", tempLength, outputLength+tempLength)

			C.EVP_CIPHER_CTX_free(ctx)
			//fmt.Printf("Freed\n")

			output = output[0 : outputLength+tempLength]

			//fmt.Printf("PLAINTEXT:\n%s\n", string(output))

			return output, nil
		}

		return ctext, nil
		// return s.c.Decrypt(out, n, ad, ciphertext)
	} else {
		return []byte{}, nil
	}
}

func (s *NebulaCipherState) Overhead() int {
	if s != nil {
		return s.c.(cipher.AEAD).Overhead()
	}
	return 0
}
