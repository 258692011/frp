package frps

import (
	"io/fs"
	"net/http"
)

var (
	// 只给 frps 使用的嵌入静态资源 FS
	content fs.FS

	// 暴露给上层的 http.FileSystem
	FileSystem http.FileSystem
)

// Register 由 web/frps/embed.go 在 init 中调用，用于注册嵌入的 dist 目录。
func Register(fileSystem fs.FS) {
	subFs, err := fs.Sub(fileSystem, "dist")
	if err != nil {
		return
	}
	content = subFs
	FileSystem = http.FS(content)
}

