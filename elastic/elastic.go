package elastic

import (
	"encoding/json"
	"errors"
	"fmt"
)

type ElasticSearch struct {
	host []string
	conn *elastigo.Conn
}
type elasticConfig struct {
	Host []string `json:"hosts"`
}

func New(config string) (host *ElasticSearch, err error) {
	var elasticonfig elasticConfig
	err = json.Unmarshal([]byte(config), &elasticonfig)
	if err != nil {
		return
	}
	host = &ElasticSearch{}
	host.conn = elastigo.NewConn()
	host.host = elasticonfig.Host
	host.conn.SetHosts(elasticonfig.Host)
	return
}

func (host *ElasticSearch) Create(name string, typeName string, jsonData string) (id string, err error) {
	response, err := host.conn.Index(name, typeName, "", nil, jsonData)
	if err != nil {
		return
	}
	id = response.Id
	host.conn.Flush()
	if response.Created {
		return
	}
	err = errors.New(fmt.Sprintf("/%s/%s create error:%+v", name, typeName, response))
	return
}
func (host *ElasticSearch) Update(name string, typeName string, id string, jsonData string) (err error) {
	response, err := host.conn.Index(name, typeName, id, nil, jsonData)
	if err != nil {
		return
	}
	host.conn.Flush()
	if response.Ok {
		return
	}
	err = errors.New(fmt.Sprintf("/%s/%s update error:%+v", name, typeName, response))
	return
}

func (host *ElasticSearch) Search(name string, typeName string, query string) (result string, err error) {
	out, err := host.conn.Search(name, typeName, nil, query)
	if err != nil {
		return
	}
	var resultLst []*json.RawMessage
	for i := 0; i < len(out.Hits.Hits); i++ {
		resultLst = append(resultLst, (out.Hits.Hits[i].Source))
	}
	buffer, err := json.Marshal(&resultLst)
	if err != nil {
		return
	}
	result = string(buffer)
	return
}
