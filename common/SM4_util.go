package common

import (
	"bytes"
	"crypto/cipher"
	"github.com/tjfoc/gmsm/sm4"
)

//func main(){
//	// 128比特密钥
//	key := []byte("1234567890abcdef")
//	// 128比特iv
//	iv := make([]byte, sm4.BlockSize)
//	data := []byte("Tongji Fintech Research Institute")
//	ciphertxt,err := sm4Encrypt(key,iv, data)
//	if err != nil{
//		log.Fatal(err)
//	}
//	fmt.Printf("加密结果: %x\n", ciphertxt)
//}

func Sm4Encrypt(plainText []byte) ([]byte, error) {
	// 128比特密钥
	key := []byte("p9qlhaoeio608kwa")
	// 128比特iv
	iv := make([]byte, sm4.BlockSize)
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData := pkcs5Padding(plainText, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, iv)
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return cryted, nil
}

func Sm4Decrypt(cipherText []byte) ([]byte, error) {
	// 128比特密钥
	key := []byte("p9qlhaoeio608kwa")
	// 128比特iv
	iv := make([]byte, sm4.BlockSize)

	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(cipherText))
	blockMode.CryptBlocks(origData, cipherText)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

// pkcs5填充
func pkcs5Padding(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(src, padtext...)
}

func pkcs5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
