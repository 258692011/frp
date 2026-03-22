//go:build !noweb

package frps

import (
	"embed"

	assetsfrps "github.com/fatedier/frp/assets/frps"
)

//go:embed dist
var EmbedFS embed.FS

func init() {
	assetsfrps.Register(EmbedFS)
}
