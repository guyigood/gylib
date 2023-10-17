package xmltomap

import (
	"bytes"
	"container/list"
	"encoding/xml"
	"io"
)

func UnmarshalXml(xmlbody []byte, v map[string]interface{}) error {
	var (
		t        xml.Token
		err      error
		curName  string
		pmapName string
		curValue string
		mapobj   map[string]interface{}
	)
	decoder := xml.NewDecoder(bytes.NewReader(xmlbody))
	mapStack := list.New()
	for t, err = decoder.Token(); err == nil; t, err = decoder.Token() {
		switch token := t.(type) {
		case xml.StartElement:
			var parentmap map[string]interface{}
			pmapName = curName
			curName = token.Name.Local
			ele := mapStack.Back()
			if ele != nil {
				parentmap = ele.Value.(map[string]interface{})
			} else {
				parentmap = v
			}
			mapobj = make(map[string]interface{})
			parentmap[curName] = mapobj
			if pmapName != "" {
				if mapvalue, ok := parentmap[pmapName]; ok {
					switch mapvalue.(type) {
					case map[string]interface{}:
					default:
						//删除值
						delete(parentmap, pmapName)
					}
				}
			}
			mapStack.PushBack(mapobj)
		case xml.CharData:
			if curName != "" {
				curValue = string([]byte(token))
				vmap := mapStack.Back().Value.(map[string]interface{})
				vmap[curName] = curValue
			} else {
				curValue = ""
			}
		case xml.EndElement:
			ele := mapStack.Back()
			mapobj := ele.Value.(map[string]interface{})
			mapStack.Remove(ele)
			if len(mapobj) == 1 {
				ele = mapStack.Back()
				if ele != nil {
					parentmap := ele.Value.(map[string]interface{})
					for k, v := range mapobj {
						parentmap[k] = v
					}
				}
			}
			curName = ""
		}
	}
	if err == io.EOF {
		err = nil
	}
	return err
}
