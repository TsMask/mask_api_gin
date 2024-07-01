package ip2region

import (
	"mask_api_gin/src/framework/logger"
	"strings"
	"time"
)

// LocalHost 网络地址(内网)
const LocalHost = "127.0.0.1"

// 全局查询对象
var searcher *Searcher

func init() {
	dbPath := "assets/ip2region.xdb"

	// 从 dbPath 加载整个 xdb 到内存
	cBuff, err := LoadContentFromFile(dbPath)
	if err != nil {
		logger.Fatalf("failed error load xdb from : %s\n", err)
		return
	}

	// 用全局的 cBuff 创建完全基于内存的查询对象。
	base, err := NewWithBuffer(cBuff)
	if err != nil {
		logger.Errorf("failed error create searcher with content: %s\n", err)
		return
	}

	// 赋值到全局查询对象
	searcher = base
}

// RegionSearchByIp 查询IP所在地
//
// 国家|区域|省份|城市|ISP
func RegionSearchByIp(ip string) (string, int, int64) {
	ip = ClientIP(ip)
	if ip == LocalHost {
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
	ip = ClientIP(ip)
	if ip == LocalHost {
		return "内网IP"
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

// ClientIP 处理客户端IP地址显示iPv4
//
// 转换 ip2region.ClientIP(c.ClientIP())
func ClientIP(ip string) string {
	if strings.HasPrefix(ip, "::ffff:") {
		ip = strings.Replace(ip, "::ffff:", "", 1)
	}
	if ip == LocalHost || ip == "::1" {
		return LocalHost
	}
	return ip
}
