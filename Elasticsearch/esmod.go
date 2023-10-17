package Elasticsearch

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
)

type EsMod struct {
	DocData          map[string]interface{}
	Client           *elastic.Client
	Ctx              context.Context
	PageSize, PageNo int
	IndexName        string
}

func NewEsMod() *EsMod {
	that := new(EsMod)
	return that.Init()
}

func (that *EsMod) Init() *EsMod {
	that.PageSize = 1000
	that.PageNo = 1
	that.DocData = make(map[string]interface{})
	that.Ctx = context.Background()
	return that
}

func (that *EsMod) SetIndexName(indexname string) *EsMod {
	that.IndexName = indexname
	return that
}

func (that *EsMod) Connect(url string) bool {
	var err error
	that.Client, err = elastic.NewClient(elastic.SetURL(url))
	if err != nil {
		return false
	}
	return true
}

func (that *EsMod) Close() {
	that.Close()
}

func (that *EsMod) SetData(data map[string]interface{}) *EsMod {
	that.DocData = make(map[string]interface{})
	for key, val := range data {
		that.DocData[key] = val
	}
	return that
}

func (that *EsMod) SetPageNo(pageno int) *EsMod {
	that.PageNo = pageno
	return that
}

func (that *EsMod) SetPageSize(pagesize int) *EsMod {
	that.PageSize = pagesize
	return that
}

func (that *EsMod) InsertDoc() (*elastic.IndexResponse, error) {
	data, err := that.Client.Index().
		Index(that.IndexName).
		BodyJson(that.DocData).
		Do(that.Ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (that *EsMod) UpdateDoc(id string) (*elastic.UpdateResponse, error) {
	data, err := that.Client.Update().
		Index(that.IndexName).
		Id(id).
		Doc(that.DocData).
		Do(that.Ctx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (that *EsMod) CreateIndex() bool {
	_, err := that.Client.CreateIndex(that.IndexName).Do(that.Ctx)
	if err != nil {
		return false
	}
	return true
}

func (that *EsMod) QueryAll() []map[string]string {
	searchResult, err := that.Client.Search().
		Index(that.IndexName).
		Query(elastic.NewMatchAllQuery()).
		From((that.PageNo - 1) * that.PageSize).
		Size(that.PageSize).
		Do(that.Ctx)
	if err != nil {
		return nil
	}
	if searchResult.TotalHits() == 0 {
		return nil
	}
	ReData := make([]map[string]string, 0)
	//fmt.Printf("Found %d documents\n", searchResult.Hits.TotalHits.Value)
	for _, hit := range searchResult.Hits.Hits {
		temp := make(map[string]string)
		temp["data"] = fmt.Sprintf("%s", hit.Source)
		ReData = append(ReData, temp)
		//fmt.Printf("Document ID: %s\n", hit.Id)
		//fmt.Printf("Document Source: %s\n", hit.Source)
	}
	return ReData
}

func (that *EsMod) LikeQuery(fdname, fdvalue string) []map[string]string {
	// 创建模糊查询
	query1 := elastic.NewWildcardQuery(fdname, fdvalue)
	searchResult, err := that.Client.Search().
		Index(that.IndexName).
		Query(query1).
		From((that.PageNo - 1) * that.PageSize).
		Size(that.PageSize).
		Do(that.Ctx)
	if err != nil {
		return nil
	}
	if searchResult.TotalHits() == 0 {
		return nil
	}
	ReData := make([]map[string]string, 0)
	//fmt.Printf("Found %d documents\n", searchResult.Hits.TotalHits.Value)
	for _, hit := range searchResult.Hits.Hits {
		temp := make(map[string]string)
		temp["data"] = fmt.Sprintf("%s", hit.Source)
		ReData = append(ReData, temp)
		//fmt.Printf("Document ID: %s\n", hit.Id)
		//fmt.Printf("Document Source: %s\n", hit.Source)
	}
	return ReData
}

func (that *EsMod) MultiMatchQuery(searchKeyword string, fdname ...string) []map[string]string {
	query1 := elastic.NewMultiMatchQuery(searchKeyword, fdname...)
	searchResult, err := that.Client.Search().
		Index(that.IndexName).
		Query(query1).
		From((that.PageNo - 1) * that.PageSize).
		Size(that.PageSize).
		Do(that.Ctx)
	if err != nil {
		return nil
	}
	if searchResult.TotalHits() == 0 {
		return nil
	}
	ReData := make([]map[string]string, 0)
	//fmt.Printf("Found %d documents\n", searchResult.Hits.TotalHits.Value)
	for _, hit := range searchResult.Hits.Hits {
		temp := make(map[string]string)
		temp["data"] = fmt.Sprintf("%s", hit.Source)
		ReData = append(ReData, temp)
		//fmt.Printf("Document ID: %s\n", hit.Id)
		//fmt.Printf("Document Source: %s\n", hit.Source)
	}
	return ReData
}

func (that *EsMod) IndexExists() bool {
	exists, err := that.Client.IndexExists(that.IndexName).Do(that.Ctx)
	if err != nil {
		return false
	}

	if exists {
		return true
	} else {
		return false
	}
}

func (that *EsMod) DelIndex() bool {
	_, err := that.Client.DeleteIndex(that.IndexName).Do(that.Ctx)
	if err != nil {
		return false
	}
	return true
}

func (that *EsMod) ComplexQuery(match_data, wildcard_data []map[string]string, MinimumShouldMatch int) ([]map[string]string, int64) {
	var exactMatch []elastic.Query
	var fuzzyMatch []elastic.Query
	for _, val := range match_data {
		temp := elastic.NewMatchQuery(val["fd_name"], val["fd_value"])
		exactMatch = append(exactMatch, temp)
	}
	for _, val := range wildcard_data {
		temp := elastic.NewWildcardQuery(val["fd_name"], val["fd_value"])
		fuzzyMatch = append(fuzzyMatch, temp)
	}
	// 使用 BoolQuery 构建复合查询
	boolQuery := elastic.NewBoolQuery().
		Must(exactMatch...).
		Should(fuzzyMatch...).
		MinimumNumberShouldMatch(MinimumShouldMatch)
	searchResult, err := that.Client.Search().
		Index(that.IndexName).
		Query(boolQuery).
		From((that.PageNo - 1) * that.PageSize).
		Size(that.PageSize).
		Do(that.Ctx)
	if err != nil {
		return nil, 0
	}
	total := searchResult.TotalHits()
	if searchResult.TotalHits() == 0 {
		return nil, 0
	}

	ReData := make([]map[string]string, 0)
	// 处理搜索结果
	for _, hit := range searchResult.Hits.Hits {
		temp := make(map[string]string)
		temp["data"] = fmt.Sprintf("%s", hit.Source)
		ReData = append(ReData, temp)
		// 处理其他字段数据
	}
	return ReData, total
}

func (that *EsMod) MultiWildcardQuery(wildcard_data []map[string]string, MinimumShouldMatch int) ([]map[string]string, int64) {
	var fuzzyMatch []elastic.Query
	for _, val := range wildcard_data {
		temp := elastic.NewWildcardQuery(val["fd_name"], val["fd_value"])
		fuzzyMatch = append(fuzzyMatch, temp)
	}
	// 使用 BoolQuery 构建复合查询
	boolQuery := elastic.NewBoolQuery().
		Should(fuzzyMatch...).
		MinimumNumberShouldMatch(MinimumShouldMatch)
	searchResult, err := that.Client.Search().
		Index(that.IndexName).
		Query(boolQuery).
		From((that.PageNo - 1) * that.PageSize).
		Size(that.PageSize).
		Do(that.Ctx)
	if err != nil {
		return nil, 0
	}
	total := searchResult.TotalHits()
	if total == 0 {
		return nil, 0
	}

	ReData := make([]map[string]string, 0)
	// 处理搜索结果
	for _, hit := range searchResult.Hits.Hits {
		temp := make(map[string]string)
		temp["data"] = fmt.Sprintf("%s", hit.Source)
		ReData = append(ReData, temp)
		// 处理其他字段数据
	}
	return ReData, total
}
