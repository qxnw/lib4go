package tpl

import (
	"fmt"
	"strings"
)

const (
	cOracle = "oracle"
	cSqlite = "sqlite"
)

var (
	tpls map[string]ITPLContext
)

//ITPLContext 模板上下文
type ITPLContext interface {
	GetSQLContext(tpl string, input map[string]interface{}) (query string, args []interface{})
	GetSPContext(tpl string, input map[string]interface{}) (query string, args []interface{})
	Replace(sql string, args []interface{}) (r string)
}

func init() {
	tpls = make(map[string]ITPLContext)
	tpls[cOracle] = OracleTPLContext{}
	tpls[cSqlite] = SqliteTPLContext{}
}

//GetDBContext 获取数据库上下文操作
func GetDBContext(name string) (ITPLContext, error) {
	if v, ok := tpls[strings.ToLower(name)]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("不支持的数据库类型:%s", name)
}
