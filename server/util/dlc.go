package util

import (
	"goload/server/models"
	"net/http"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"io/ioutil"
	"encoding/xml"
)

const api = "http://service.jdownloader.org/dlcrypt/service.php?srcType=dlc&destType=%s&data=%s"

func DecodeDLC(dlc string) (files[] *models.File) {
	KEY   := []byte("cb99b5cbc24db398")
	IV    := []byte("9bc24cb995cb8db3")
	urls := decrypt(dlc,KEY,IV)

	for _,url:= range urls {
		files = append(files,&models.File{Url:url})
	}
	return files
}

func decrypt(dlc string, key []byte, iv []byte) (urls[] string) {
	dlckey := dlc[len(dlc)-88:]
	data, _ := base64.StdEncoding.DecodeString(dlc[:len(dlc)-88])
	resp, _ := http.Get("http://service.jdownloader.org/dlcrypt/service.php?srcType=dlc&destType=pylo&data="+dlckey)
	body, _ := ioutil.ReadAll(resp.Body)
	stringBody := string(body)
	dlc_key :=stringBody[4:len(stringBody)-5]
	println(dlc_key)
	block, err := aes.NewCipher(key)
	if err != nil{
		print(err.Error())
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	decode,_ := base64.StdEncoding.DecodeString(dlc_key)
	mode.CryptBlocks(decode, decode)
	block, _ = aes.NewCipher(decode)
	mode = cipher.NewCBCDecrypter(block, decode)
	mode.CryptBlocks(data, data)
	links,_ := base64.StdEncoding.DecodeString(string(data))

	type Result struct {
		XMLName xml.Name `xml:"dlc"`
		URLs   []string `xml:"content>package>file>url"`
	}
	v := Result{}
	xmlerr := xml.Unmarshal(links, &v)
	if xmlerr!= nil{
		print(xmlerr.Error())
	}
	for _,url := range v.URLs {
		value,_:= base64.StdEncoding.DecodeString(url)
		urls = append(urls,(string(value)))
	}
	return urls
}