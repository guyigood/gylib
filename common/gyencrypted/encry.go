package GyAesEncrypted

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"log"
)

type GyAesEncrypt struct {
	key, iv []byte
}

func NewGyAesEncrypt(key, iv string) *GyAesEncrypt {
	this := new(GyAesEncrypt)
	this.key = []byte(key)
	this.iv = []byte(iv)
	return this
}

func (this *GyAesEncrypt) Encrypt(text []byte) (string, error) {
	//生成cipher.Block 数据块
	block, err := aes.NewCipher(this.key)
	if err != nil {
		log.Println("错误 -" + err.Error())
		return "", err
	}
	//填充内容，如果不足16位字符
	blockSize := block.BlockSize()
	originData := this.pad(text, blockSize)
	//加密方式
	blockMode := cipher.NewCBCEncrypter(block, this.iv)
	//加密，输出到[]byte数组
	crypted := make([]byte, len(originData))
	blockMode.CryptBlocks(crypted, originData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func (this *GyAesEncrypt) pad(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func (this *GyAesEncrypt) Decrypt(text string) (string, error) {
	decode_data, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", nil
	}
	//生成密码数据块cipher.Block
	block, _ := aes.NewCipher(this.key)
	//解密模式
	blockMode := cipher.NewCBCDecrypter(block, this.iv)
	//输出到[]byte数组
	origin_data := make([]byte, len(decode_data))
	blockMode.CryptBlocks(origin_data, decode_data)
	//去除填充,并返回
	return string(this.unpad(origin_data)), nil
}

func (this *GyAesEncrypt) unpad(ciphertext []byte) []byte {
	length := len(ciphertext)
	//去掉最后一次的padding
	unpadding := int(ciphertext[length-1])
	return ciphertext[:(length - unpadding)]
}
