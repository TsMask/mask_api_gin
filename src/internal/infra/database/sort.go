package database

import (
	"strings"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SortOrder 排序项
type SortOrder struct {
	Field string // 数据库列名
	Asc   bool   // true=ASC, false=DESC
}

// SortParse 解析统一排序字符串
// 格式: sort=field1,field2:asc,field3:desc
// 不加方向默认 DESC，方向不区分大小写
// fieldMapping 用于将前端字段名映射为数据库列名，未匹配的字段名原样传递
func SortParse(sortStr string, fieldMapping map[string]string) []SortOrder {
	if sortStr == "" {
		return nil
	}
	parts := strings.Split(sortStr, ",")
	orders := make([]SortOrder, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		// 拆分 field:direction
		field, direction, _ := strings.Cut(part, ":")
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}
		// 字段映射
		col := field
		if mapped, ok := fieldMapping[field]; ok {
			col = mapped
		}
		// 默认 DESC
		asc := false
		if dir := strings.ToLower(strings.TrimSpace(direction)); dir == "asc" {
			asc = true
		}
		orders = append(orders, SortOrder{Field: col, Asc: asc})
	}
	return orders
}

// SortOrDefault 将排序项应用到 gorm 查询，无排序项时使用默认排序
func SortOrDefault(tx *gorm.DB, orders []SortOrder, defaultOrder string) *gorm.DB {
	if len(orders) > 0 {
		columns := make([]clause.OrderByColumn, 0, len(orders))
		for _, o := range orders {
			columns = append(columns, clause.OrderByColumn{
				Column: clause.Column{Name: o.Field},
				Desc:   !o.Asc,
			})
		}
		return tx.Order(clause.OrderBy{Columns: columns})
	}
	return tx.Order(defaultOrder)
}
