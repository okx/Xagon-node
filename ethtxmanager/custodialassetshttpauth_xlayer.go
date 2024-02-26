package ethtxmanager

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

const (
	algoSha256      = "sha256"
	algoMd5         = "md5"
	headerSignKey   = "sign"
	headerAccessKey = "accessKey"
)

// auth is SecretKey or AccessKey is empty then no auth is performed
func (c *Client) auth(ctx context.Context, req *http.Request) error {
	if c.cfg.CustodialAssets.SecretKey == "" || c.cfg.CustodialAssets.AccessKey == "" {
		return nil
	}
	req.Header[headerAccessKey] = []string{c.cfg.CustodialAssets.AccessKey}
	// Generate the signature
	signature, err := c.genAuth(ctx, req, algoSha256)
	if err != nil {
		return fmt.Errorf("failed to generate signature: %v", err)
	}
	req.Header[headerSignKey] = []string{signature}

	return nil
}

func (c *Client) genAuth(ctx context.Context, req *http.Request, algorithm string) (string, error) {
	params := req.URL.Query()

	treeMap := make(map[string][]string)
	for _, v := range params {
		var key string
		for _, vv := range v {
			key += vv
		}
		treeMap[key] = v
	}

	var body strings.Builder
	if req.Body != nil {
		readCloser, err := req.GetBody()
		if err != nil {
			return "", fmt.Errorf("get body error: %v", err)
		}
		defer readCloser.Close()

		buffer := make([]byte, 0)
		// Read the request body into the 'buffer' variable
		_, err = io.Copy(&bufferWriter{&buffer}, readCloser)
		if err != nil {
			return "", fmt.Errorf("read body error: %v", err)
		}

		// Append the body if present
		if len(buffer) != 0 {
			body.Write(buffer)
		}
	}

	// Calculate the hash based on the selected algorithm
	return c.generateSignature(ctx, treeMap, body.String(), algorithm)
}

func (c *Client) generateSignature(ctx context.Context, treeMap map[string][]string, body, algorithm string) (string, error) {
	// Sort the map by values
	var keys []string
	for key := range treeMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Construct the content string
	var content strings.Builder
	for _, key := range keys {
		value := treeMap[key]
		for _, v := range value {
			content.WriteString(v)
		}
	}
	// Append the body if present
	if body != "" {
		content.WriteString(body)
	}

	// Calculate the hash based on the selected algorithm
	var hash []byte
	switch strings.ToLower(algorithm) {
	case "sha256":
		hash = signBySha256(content.String())
	case "md5":
		hash = signByMd5(content.String())
	default:
		// Handle unsupported algorithm
		return "", fmt.Errorf("unsupported algorithm: %v", algorithm)
	}

	// Convert the hash to a hexadecimal string
	hashString := hex.EncodeToString(hash)

	// Encrypt the hash using AES
	return encryptAES(hashString, c.cfg.CustodialAssets.SecretKey)
}

// bufferWriter is a simple implementation of io.Writer to write to a buffer
type bufferWriter struct {
	buffer *[]byte
}

func (bw *bufferWriter) Write(p []byte) (n int, err error) {
	*bw.buffer = append(*bw.buffer, p...)
	return len(p), nil
}

func signByMd5(content string) []byte {
	hash := md5.New() // golint:ignore
	hash.Write([]byte(content))
	return []byte(hex.EncodeToString(hash.Sum(nil)))
}

func signBySha256(content string) []byte {
	hash := sha256.New()
	hash.Write([]byte(content))
	return hash.Sum(nil)
}

func encryptAES(src, key string) (string, error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", fmt.Errorf("failed to create AES cipher: %v", err)
	}
	ecbEncrypt := newECBEncrypter(block)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	des := make([]byte, len(content))
	err = ecbEncrypt.(*ecbEncrypter).cryptBlocksWithError(des, content)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt AES: %v", err)
	}

	return base64.StdEncoding.EncodeToString(des), nil
}

type ecb struct {
	b         cipher.Block
	blockSize int
}

func newECB(b cipher.Block) *ecb {
	return &ecb{
		b:         b,
		blockSize: b.BlockSize(),
	}
}

type ecbEncrypter ecb

// newECBEncrypt returns a BlockMode which encrypts in electronic code book
// mode, using the given Block.
func newECBEncrypter(b cipher.Block) cipher.BlockMode {
	return (*ecbEncrypter)(newECB(b))
}
func (x *ecbEncrypter) BlockSize() int { return x.blockSize }
func (x *ecbEncrypter) CryptBlocks(dst, src []byte) {
	if len(src)%x.blockSize != 0 {
		panic("crypto/cipher: input not full blocks")
	}
	if len(dst) < len(src) {
		panic("crypto/cipher: output smaller than input")
	}
	for len(src) > 0 {
		x.b.Encrypt(dst, src[:x.blockSize])
		src = src[x.blockSize:]
		dst = dst[x.blockSize:]
	}
}

func (x *ecbEncrypter) cryptBlocksWithError(dst, src []byte) error {
	if len(src)%x.blockSize != 0 {
		return fmt.Errorf("crypto/cipher: input not full blocks %v, %v", len(src), x.blockSize)
	}
	if len(dst) < len(src) {
		return fmt.Errorf("crypto/cipher: output smaller than input %v, %v", len(dst), len(src))
	}

	x.CryptBlocks(dst, src)

	return nil
}

// PKCS5Padding padding ciphertext
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}
