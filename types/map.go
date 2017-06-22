package types

func GetIMap(m map[string]string) map[string]interface{} {
	n := make(map[string]interface{})
	for k, v := range m {
		n[k] = v
	}
	return n
}
