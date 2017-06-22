package transform

const (
	SymbolEqual = iota
	SymbolNotEqual
	SymbolMore
	SymbolLess
)

//Expression 表达式
type Expression struct {
	Left   string
	Symbol int
	Right  string
}

/*
func parseQuery(query string) (err error) {
	for query != "" {
		key := query
		if i := strings.IndexAny(key, "&;"); i >= 0 {
			key, query = key[:i], key[i+1:]
		} else {
			query = ""
		}
		if key == "" {
			continue
		}
		value := ""
		if i := strings.Index(key, "="); i >= 0 {
			key, value = key[:i], key[i+1:]
		}
		m[key] = append(m[key], value)
	}
	return err
}
*/
