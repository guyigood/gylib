package common

import (
	"encoding/base64"
	"io/ioutil"
)

func GetBase64FromFile(filename string) (string, error) {
	fileBytes, err := ioutil.ReadFile(filename) // 读取file
	if err != nil {
		return "", err
	}
	bs64 := base64.StdEncoding.EncodeToString(fileBytes) // 加密成base64字符串
	return bs64, nil
}

func SaveBase64ToFile(bs64, filename string) error {
	bs64Bytes, err := base64.StdEncoding.DecodeString(bs64) // 解密base64字符串
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, bs64Bytes, 0666) // 写入file
	return err
}
