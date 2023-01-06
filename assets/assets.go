package assets

import "embed"

//go:embed chart chart/.* chart/**/_*
var Assets embed.FS
