package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
)

//func main() {
//	testAes()
//	Test_B_2()
//}

type Login struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

// CBC模式
func TestAes() {
	// AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
	str := `{"username":"admin","password":"Cdp12345"}`
	orig := []byte(str) // 待加密的数据
	key := []byte("silvanware123456")
	result, err := AesEncrypt(orig, key)
	if err != nil {
		panic(err)
	}
	log.Println("------------------ CBC模式 --------------------")
	log.Println(result)
	log.Println("密文(hex)：", hex.EncodeToString(result))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(result))
	origData, err := AesDecrypt(result, key)
	if err != nil {
		panic(err)
	}
	log.Println("解密结果：", string(origData))
}

// ECB模式
func Test_B_2() {
	str := fmt.Sprintf(`{"username":"%s","password":"%s"}`, "admin", "Cdp12345")
	//data := url.Values{}
	//data.Set("username", "admin")
	//data.Set("password", "Cdp12345")
	origData := []byte(str)           // 待加密的数据
	key := []byte("silvanware123456") // 加密的密钥
	log.Println("原文：", string(origData))

	log.Println("------------------ ECB模式 --------------------")
	encrypted := AesEncryptECB(origData, key)
	log.Println(encrypted)
	log.Println("密文(hex)：", hex.EncodeToString(encrypted))
	log.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := AesDecryptECB(encrypted, key)
	log.Println("解密结果：", string(decrypted))
}

func AesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

// 填充
func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// 去除填充
func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

// ECB加密
func AesEncryptECB(origData []byte, key []byte) (encrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	length := (len(origData) + aes.BlockSize) / aes.BlockSize
	plain := make([]byte, length*aes.BlockSize)
	copy(plain, origData)
	pad := byte(len(plain) - len(origData))
	for i := len(origData); i < len(plain); i++ {
		plain[i] = pad
	}
	encrypted = make([]byte, len(plain))
	// 分组分块加密
	for bs, be := 0, cipher.BlockSize(); bs <= len(origData); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Encrypt(encrypted[bs:be], plain[bs:be])
	}

	return encrypted
}

//ECB解密
func AesDecryptECB(encrypted []byte, key []byte) (decrypted []byte) {
	cipher, _ := aes.NewCipher(generateKey(key))
	decrypted = make([]byte, len(encrypted))
	//
	for bs, be := 0, cipher.BlockSize(); bs < len(encrypted); bs, be = bs+cipher.BlockSize(), be+cipher.BlockSize() {
		cipher.Decrypt(decrypted[bs:be], encrypted[bs:be])
	}

	trim := 0
	if len(decrypted) > 0 {
		trim = len(decrypted) - int(decrypted[len(decrypted)-1])
	}

	return decrypted[:trim]
}

func generateKey(key []byte) (genKey []byte) {
	genKey = make([]byte, 16)
	copy(genKey, key)
	for i := 16; i < len(key); {
		for j := 0; j < 16 && i < len(key); j, i = j+1, i+1 {
			genKey[j] ^= key[i]
		}
	}
	return genKey
}
