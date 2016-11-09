package db

import (
	"reflect"
	"testing"
)

type tqDbData struct {
	data []map[string]interface{}
	cols []string
	err  error
}

type teDBData struct {
	row int64
	err error
}
type tpData struct {
	query  string
	input  map[string]interface{}
	args   []interface{}
	result tqDbData
	tp     int
}

type tDB struct {
	qdata map[string]tqDbData
	edata map[string]teDBData
}

type tTPL struct {
	data map[string]tpData
}

func (t *tTPL) GetSQLContext(tpl string, input map[string]interface{}) (query string, args []interface{}) {
	return t.data[tpl].query, t.data[tpl].args
}
func (t *tTPL) GetSPContext(tpl string, input map[string]interface{}) (query string, args []interface{}) {
	return t.data[tpl].query, t.data[tpl].args
}
func (t *tTPL) Replace(sql string, args []interface{}) (r string) {
	return "REPLACE"
}

func (t *tDB) Query(q string, input ...interface{}) ([]map[string]interface{}, []string, error) {
	return t.qdata[q].data, t.qdata[q].cols, t.qdata[q].err
}
func (t *tDB) Execute(q string, input ...interface{}) (int64, error) {
	return t.edata[q].row, t.edata[q].err
}
func (t *tDB) Begin() (IDBTrans, error) {
	return nil, nil
}

func makeDB(q map[string]tqDbData, e map[string]teDBData, p map[string]tpData) *DB {
	d := &DB{}
	tdb := &tDB{}
	tdb.edata = e
	tdb.qdata = q
	d.db = tdb
	ttp := &tTPL{}
	ttp.data = p
	d.tpl = ttp
	return d
}

func TestDBQuery(t *testing.T) {
	queryMap := make(map[string]tqDbData)
	executeMap := make(map[string]teDBData)
	tplMap := make(map[string]tpData)

	tplMap = map[string]tpData{
		"select a from dual": tpData{
			tp:    1,
			query: "select 'a' from dual",
			args:  []interface{}{"a", 1},
			input: map[string]interface{}{
				"name": "colin",
			},
		},
		"update order set t=1 where id=2": tpData{
			tp:    2,
			query: "update order set t=1 where id='2'",
			args:  []interface{}{"a", 1},
			input: map[string]interface{}{
				"name": "colin",
			},
		},
	}

	queryMap["select 'a' from dual"] = tqDbData{
		err:  nil,
		cols: []string{"name"},
		data: []map[string]interface{}{
			map[string]interface{}{
				"name": "colin",
			},
		},
	}

	executeMap["update order set t=1 where id='2'"] = teDBData{
		err: nil,
		row: 2,
	}

	db := makeDB(queryMap, executeMap, tplMap)
	for k, v := range tplMap {
		switch v.tp {
		case 1:
			result, sql, input, err := db.Query(k, v.input)
			if !reflect.DeepEqual(result, queryMap[v.query].data) || sql != v.query || err != nil || !reflect.DeepEqual(input, v.args) {
				t.Error("ExecuteSP返回参数有误", len(result), len(queryMap[v.query].data))
			}
			dt, sql, input, err := db.Scalar(k, v.input)
			if dt != queryMap[v.query].data[0]["name"] || sql != v.query || err != nil || !reflect.DeepEqual(input, v.args) {
				t.Error("Scalar", len(result), len(queryMap[v.query].data))
			}

		case 2:
			row, sql, input, err := db.Execute(k, v.input)
			if row != executeMap[v.query].row || sql != v.query || err != nil || !reflect.DeepEqual(input, v.args) {
				t.Error("execute返回参数有误")
			}
			row, sql, input, err = db.ExecuteSP(k, v.input)
			if row != executeMap[v.query].row || sql != v.query || err != nil || !reflect.DeepEqual(input, v.args) {
				t.Error("ExecuteSP返回参数有误")
			}
		}
	}
}
