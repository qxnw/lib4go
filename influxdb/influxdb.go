package influxdb

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/qxnw/lib4go/utility"
)

type influxDbConfig struct {
	Address   string `json:"address"`
	DbName    string `json:"db"`
	UserName  string `json:"user"`
	Password  string `json:"password"`
	RowFormat string `json:"row"`
}
type InfluxDB struct {
	config *influxDbConfig
}

func New(config string) (i *InfluxDB, err error) {
	fmt.Println("new influxdb:", config)
	i = &InfluxDB{}
	i.config = &influxDbConfig{}
	err = json.Unmarshal([]byte(config), &i.config)
	if err != nil {
		return
	}
	if strings.EqualFold(i.config.Address, "") ||
		strings.EqualFold(i.config.DbName, "") ||
		strings.EqualFold(i.config.RowFormat, "") {
		err = errors.New("influxDbConfig必须参数不能为空")
		return
	}
	return
}
func (db *InfluxDB) SaveString(rows string) (err error) {
	var data []map[string]interface{}
	err = json.Unmarshal([]byte(rows), &data)
	if err != nil {
		return
	}
	return db.Save(data)
}

func (db *InfluxDB) Save(rows []map[string]interface{}) (err error) {
	url := fmt.Sprintf("%s/write?db=%s", db.config.Address, db.config.DbName)
	var datas []string
	for i := 0; i < len(rows); i++ {
		d := utility.NewDataMaps(rows[i])
		datas = append(datas, d.Translate(db.config.RowFormat))
	}
	data := strings.Join(datas, "\n")
	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data))
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 204 {
		return nil
	}
	err = errors.New(fmt.Sprintf("influxdb save error:%d", resp.StatusCode))
	return
}
