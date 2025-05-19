package pip

import (
	"fmt"
	"github.com/sfc-gh-bprosnitz/go-embed-python/embed_util"
	"github.com/sfc-gh-bprosnitz/go-embed-python/pip/internal/data"
)

func NewPipLib(name string) (*embed_util.EmbeddedFiles, error) {
	return embed_util.NewEmbeddedFiles(data.Data, fmt.Sprintf("pip-%s", name))
}
