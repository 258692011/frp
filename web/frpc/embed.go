package frpc

import (
	"embed"

	assetsfrpc "github.com/fatedier/frp/assets/frpc"
)

//go:embed dist
var EmbedFS embed.FS

func init() {
	assetsfrpc.Register(EmbedFS)
}
