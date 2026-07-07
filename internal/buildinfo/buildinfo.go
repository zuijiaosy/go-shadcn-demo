// Package buildinfo 保存构建时通过 ldflags 注入的版本信息
package buildinfo

var (
	// Version 版本号（git describe）
	Version = "dev"
	// GitCommit Git 提交哈希
	GitCommit = "unknown"
	// BuildTime 构建时间（UTC）
	BuildTime = "unknown"
)
