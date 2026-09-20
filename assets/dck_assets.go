package assets

import "embed"

// DCKAssetFiles shares an embedded resource with the optional DCK version.
func DCKAssetFiles() embed.FS { return Files }
