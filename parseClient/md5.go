package parseClient

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
)
//解密
func AESDecrypt(crypted, key []byte) []byte {
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()
	a := make([]byte, 16)
	blockMode := cipher.NewCBCDecrypter(block, a[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData
}

//去补码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:length-unpadding]
}

//加密
func AESEncrypt(origData, key []byte) []byte {
	//获取block块
	block, _ := aes.NewCipher(key)
	//补码
	origData = PKCS7Padding(origData, block.BlockSize())

	a := make([]byte, 16)
	//加密模式，
	//blockMode := cipher.NewCBCEncrypter(block, a[:block.BlockSize()])
	blockMode := cipher.NewCBCEncrypter(block, a)
	//创建明文长度的数组
	crypted := make([]byte, len(origData))
	//加密明文
	blockMode.CryptBlocks(crypted, origData)
	return crypted
}

//补码
func PKCS7Padding(origData []byte, blockSize int) []byte {
	//计算需要补几位数
	padding := blockSize - len(origData)%blockSize
	//d := make([]byte, 16)
	//在切片后面追加char数量的byte(char)
	//padtext := bytes.Repeat(d, padding)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(origData, padtext...)
}

func tes() {
	//设置是要
	key := []byte("cqyshkgsetermsdk")
	//明文
	//origData := []byte(`{"OtherOfficeNo":"CKG177","TicketNo":"6183397350753","User":{"AppCaptcha":"","AppName":"refund_ticket","AppPwd":"refund_ticket","ConfigGroup":0,"Department":"技术中心","UserName":"退票中心"}}`)
	data := map[string]interface{}{
		"OtherOfficeNo": "CKG177", "TicketNo": "6183397350753",
		"User": map[string]interface{}{"AppCaptcha": "", "AppName": "refund_ticket", "AppPwd": "refund_ticket", "ConfigGroup": 0, "Department": "技术中心", "UserName": "退票中心"}}
	dataB, err := json.Marshal(data)
	fmt.Println("err:", err)
	//加密
	en := AESEncrypt(dataB, key)

	//RoK18en/NWUhdZpc8rABIA==
	s := base64.StdEncoding.EncodeToString(en)
	fmt.Println(s)
	ss := []byte("FdyNdO18F7mWLQwm4teYn/ia50CylL8xsaHgu5t5hSo6NARzsg5b72oBY9GmsaSZ9A9T7ybtT586SLW08KlRgI1oH0lWwH1V1xFdMN480mPmMHYggQvnsoJt/cCqCowLlkTuhvXzMYHwKv+yH/ipjGZdtISh05zLDS7lFHHvNsn1lDaGfVoKB0TelKGypkDilBTcu0OxS2P5bMS0wqJ6D2NSOO3dmqDnDqZsjTpWWM0=")
	//aa := []byte("")

	by, err := base64.StdEncoding.DecodeString(string(ss))
	fmt.Println(ss, by, err)
	////解密
	//de := AESDecrypt(by, key)
	//fmt.Println(string(de))
}

