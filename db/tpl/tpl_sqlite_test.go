package tpl

import "testing"

func TestSqliteTPLGetContext(t *testing.T) {
	sqlite := SqliteTPLContext{}
	input := make(map[string]interface{})
	input["id"] = 1
	input["name"] = "colin"

	//正确参数解析
	tpl := "select seq_wxaccountmenu_auto_id.nextval from where id=@id and name=@name"
	except := "select seq_wxaccountmenu_auto_id.nextval from where id=? and name=?"
	actual, params := sqlite.GetSQLContext(tpl, input)
	if actual != except || len(params) != 2 || params[0] != input["id"] || params[1] != input["name"] {
		t.Error("GetSQLContext解析参数有误")
	}

	//正确参数解析o
	tpl = "select seq_wxaccountmenu_auto_id.nextval from where id=@id \r\nand name=@name"
	except = "select seq_wxaccountmenu_auto_id.nextval from where id=? \r\nand name=?"
	actual, params = sqlite.GetSQLContext(tpl, input)
	if actual != except || len(params) != 2 || params[0] != input["id"] || params[1] != input["name"] {
		t.Error("GetSQLContext解析参数有误")
	}

}

func TestSqliteTPLGetSPContext(t *testing.T) {
	sqlite := SqliteTPLContext{}
	input := make(map[string]interface{})
	input["id"] = 1
	input["name"] = "colin"

	//正确参数解析
	tpl := "order_create(@id,@name,@colin)"
	except := "order_create(?,?,?)"
	actual, params := sqlite.GetSPContext(tpl, input)
	if actual != except || len(params) != 3 || params[0] != input["id"] || params[1] != input["name"] || params[2] != nil {
		t.Error("GetSPContext解析参数有误")
	}
}

func TestSqliteTPLReplace(t *testing.T) {
	orcl := SqliteTPLContext{}
	input := make([]interface{}, 0, 2)

	tpl := "begin order_create(?,?,?);end;"
	except := "begin order_create(NULL,NULL,NULL);end;"
	actual := orcl.Replace(tpl, input)
	if actual != except {
		t.Error("Replace解析参数有误", actual)
	}

	tpl = ""
	except = ""
	actual = orcl.Replace(tpl, input)
	if actual != except {
		t.Error("Replace解析参数有误", actual)
	}

	input = append(input, 1)
	input = append(input, "colin")

	tpl = "begin order_create(?,?,?);end;"
	except = "begin order_create('1','colin',NULL);end;"
	actual = orcl.Replace(tpl, input)
	if actual != except {
		t.Error("Replace解析参数有误", actual)
	}

	tpl = "begin order_create(?);end;"
	except = "begin order_create('1');end;"
	actual = orcl.Replace(tpl, input)
	if actual != except {
		t.Error("Replace解析参数有误", actual)
	}

	tpl = "begin order_create(?,?,?);end;"
	except = "begin order_create('1','colin',NULL);end;"
	actual = orcl.Replace(tpl, input)
	if actual != except {
		t.Error("Replace解析参数有误", actual)
	}

	tpl = "begin order_create(?,'?234');end;"
	except = "begin order_create('1','?234');end;"
	actual = orcl.Replace(tpl, input)
	if actual != except {
		t.Error("Replace解析参数有误", actual)
	}

}
