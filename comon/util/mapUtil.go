package util

// 将 sourceMap 合并到 targetMap 中
func MapMerge(targetMap, sourceMap map[string]func([]byte)interface{}) {
	for k, v := range sourceMap {
		targetMap[k] = v
	}
}