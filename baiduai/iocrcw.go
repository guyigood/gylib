package baiduai

import (
	"fmt"
	"gylib/common/datatype"
)

type IocrInput struct {
	Image    string `json:"image"`                  //image	true	string	-	图像数据，base64编码后进行urlencode，要求base64编码和urlencode后大小不超过4M，最短边至少15px，最长边最大4096px,支持jpg/jpeg/png/bmp格式
	Template string `json:"templateSign,omitempty"` //false	string	-	自定义模板的ID，举例：Nsdax2424asaAS791823112
	Dectid   int    `json:"detectorId,omitempty"`   //"	false	int	0	检测器ID，可选值仅有0。detectorId = 0时，启用混贴票据识别功能，可对发票粘贴单上的多张不同种类票据进行分类检测和识别
	ClassId  int    `json:"classifierId"`           //"	false	int	-	分类器Id，用于指定使用哪个分类器。
}

type IocrRF struct {
	ErrorCode int            `json:"error_code"` //error_code	int	0代表成功，如果有错误码返回可以参考下方错误码列表排查问题
	ErrorMsg  string         `json:"error_msg"`  //error_msg	string	如果error_code具体的失败信息，可以参考下方错误码列表排查问题
	LogId     string         `json:"logId"`      //logId	string	调用的日志id
	data      IocrJsonObject `json:"data"`
	//data	jsonObject	识别返回的结果
}

type IocrJsonObject struct {
	/*ret []	jsonArray	识别出来的字段数组，每一个单元里包含以下几个元素
	probability	jsonObject	字段的置信度，包括最大，最小和方差
	location	jsonObject	字段在原图上对应的矩形框位置，通过上边距、左边距、宽度、高度表示
	word_name	string	isStructured 为 true 时存在，表示字段的名字；如果 isStructured 为 false 时，不存在
	word	string	识别的字符串或单字
	templateSign	string	图片分类结果对应的模板id或指定使用的模版id。
	detectorId = 0时，对上传的发票粘贴单中的多张不同票据进行检测分类，返回每张发票的类别，templateSign的对应关系为：
	- vat_invoice：增值税发票；
	- taxi：出租车票；
	- roll_ticket：卷票；
	- train_ticket：火车票；
	- quota_invoice：定额发票；
	- travel_itinerary：行程单；
	- printed_invoice：机打发票。
	scores	float	分类置信度，如果指定templateSign，则该值为1
	isStructured	string	表示是否结构化成功，true为成功，false为失败；成功时候，返回结构化的识别结果；失败时，如果能识别，按行返回结果，如果不能识别，返回空*/
}

func (this *BaiDuAI) IocrNewInput() {
	this.IorcIn = new(IocrInput)
}

func (this *BaiDuAI) IocrRecogniseFinance() bool {
	if this.AccessToken == "" {
		fmt.Println("token未获取")
		return false
	}
	if this.imgBase64 == "" {
		fmt.Println("未设置图片信息")
		return false
	}
	this.IorcIn.Image = this.imgBase64
	url := "https://aip.baidubce.com/rest/2.0/solution/v1/iocr/recognise/finance?access_token=" + this.AccessToken
	result := this.curlweb.Https_post(url, datatype.Struct2DBMap(this.IorcIn))
	fmt.Println(result)
	return true
}
