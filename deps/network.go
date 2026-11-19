package deps

import (
	"net"
	"time"
)

// CheckTileServer 快速探测在线地图瓦片服务是否可达。
//
// 飞行路径地图依赖在线瓦片，离线环境下渲染必然失败。
// 这里用 3 秒超时做 TCP 连通性探测，提前给出明确提示，
// 而不是让用户等待后拿到一个晦涩的渲染错误。
func CheckTileServer() bool {
	conn, err := net.DialTimeout("tcp", "tile.openstreetmap.org:443", 3*time.Second)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}
