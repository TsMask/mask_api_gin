package parse

import "strconv"

// 解析数值型
func Number(str interface{}) int {
	switch str := str.(type) {
	case string:
		if str == "" {
			return 0
		}
		if num, err := strconv.Atoi(str); err == nil {
			return num
		}
	case int:
		return str
	}
	return 0
}

// 数组内字符串去重
func RemoveDuplicates(ids []string) []string {
	uniqueIDs := make(map[string]bool)
	uniqueIDSlice := make([]string, 0, len(ids))

	for _, id := range ids {
		if _, ok := uniqueIDs[id]; !ok {
			uniqueIDs[id] = true
			uniqueIDSlice = append(uniqueIDSlice, id)
		}
	}

	return uniqueIDSlice
}
