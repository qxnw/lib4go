package types

func GetIMap(m map[string]string) map[string]interface{} {
	n := make(map[string]interface{})
	for k, v := range m {
		n[k] = v
	}
	return n
}
func GetMapValue(key string, m map[string]string, def ...string) string {
	if v, ok := m[key]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}
