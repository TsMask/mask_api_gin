package ip2region

import (
	"mask_api_gin/src/framework/constants/common"
	"mask_api_gin/src/framework/logger"
	"strings"
	"time"
)

var searcher *Searcher

func init() {
	dbPath := "src/assets/ip2region.xdb"

	// 1、从 dbPath 加载整个 xdb 到内存
	cBuff, err := LoadContentFromFile(dbPath)
	if err != nil {
		logger.Errorf("failed to load content from `%s`: %s\n", dbPath, err)
		return
	}

	// 2、用全局的 cBuff 创建完全基于内存的查询对象。
	base, err := NewWithBuffer(cBuff)
	if err != nil {
		logger.Errorf("failed to create searcher with content: %s\n", err)
		return
	}
	searcher = base
}

// RegionSearchByIp 查询IP所在地
//
// 国家|区域|省份|城市|ISP
func RegionSearchByIp(ip string) (string, int, int64) {
	if ip == "::1" || strings.HasPrefix(ip, common.IP_INNER_ADDR) {
		return "0|0|0|内网IP|内网IP", 0, 0
	}
	tStart := time.Now()
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		logger.Errorf("failed to SearchIP(%s): %s\n", ip, err)
		return "0|0|0|0|0", 0, 0
	}
	return region, 0, time.Since(tStart).Milliseconds()
}

// RealAddressByIp 地址IP所在地
//
// 218.4.167.70 江苏省 苏州市
func RealAddressByIp(ip string) string {
	if ip == "::1" || strings.HasPrefix(ip, common.IP_INNER_ADDR) {
		return common.IP_INNER_LOCATION
	}
	region, err := searcher.SearchByStr(ip)
	if err != nil {
		logger.Errorf("failed to SearchIP(%s): %s\n", ip, err)
		return "未知"
	}
	parts := strings.Split(region, "|")
	province := parts[2]
	city := parts[3]
	if province == "0" && city != "0" {
		return city
	}
	return province + " " + city
}
