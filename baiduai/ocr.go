package baiduai

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type OcrReturn struct {
	LogId          uint64       `json:"log_id"`           //	log_id	是	uint64	唯一的log id，用于问题定位
	WordsResultNum uint32       `json:"words_result_num"` //words_result_num	是	uint32	识别结果数，表示words_result的元素个数
	WordsResult    []WordReulst `json:"words_result"`     //words_result	是	array[]	定位和识别结果数组
	Chars          []CharResult `json:"chars"`            //chars	否	array[]	单字符结果，recognize_granularity=small时存在
	Probability    float32      `json:"probability"`      //probability	否	float	识别结果中每一行的置信度值，包含average：行置信度平均值，variance：行置信度方差，min：行置信度最小值
}

type WordReulst struct {
	Location OcrLocation `json:"location"` // location	是	object{}	位置数组（坐标0点为左上角）
	Words    string      `json:"words"`    //words	是	string	识别结果字符串
}

type OcrLocation struct {
	Left   uint32 `json:"left"`   //left	是	uint32	表示定位位置的长方形左上顶点的水平坐标
	Top    uint32 `json:"top"`    //top	是	uint32	表示定位位置的长方形左上顶点的垂直坐标
	Width  uint32 `json:"width"`  //width	是	uint32	表示定位位置的长方形的宽度
	Heigth uint32 `json:"heigth"` //height	是	uint32	表示定位位置的长方形的高度
}

type CharResult struct {
	Location OcrLocation `json:"location"` // location	是	object{}	位置数组（坐标0点为左上角）
	Char     string      `json:"char"`     //words	是	string	识别结果字符串
}

func (this *BaiDuAI) SetImage(img_path string) *BaiDuAI {
	image, err := ioutil.ReadFile(img_path)
	if err != nil {
		fmt.Println(err)
		return this
	}
	this.imgBase64 = base64.StdEncoding.EncodeToString(image)
	return this
}

func (this *BaiDuAI) OcrReceipt() bool {
	if this.AccessToken == "" {
		fmt.Println("token未获取")
		return false
	}
	if this.imgBase64 == "" {
		fmt.Println("未设置图片信息")
		return false
	}
	url := "https://aip.baidubce.com/rest/2.0/ocr/v1/receipt?access_token=" + this.AccessToken
	result := this.curlweb.Https_post(url, map[string]interface{}{"image": this.imgBase64})
	if result == "" {
		return false
	}
	//fmt.Println(result)
	this.OcrResult = new(OcrReturn)
	err := json.Unmarshal([]byte(result), this.OcrResult)
	if err != nil {
		fmt.Println(err)
		return false
	}

	//fmt.Println(this.OcrResult)
	return true

}
