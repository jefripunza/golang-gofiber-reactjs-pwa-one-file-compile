package variable

import (
	"backend/server/env"
	"path/filepath"
)

var TempPath = filepath.Join(env.GetPwd(), "temp")
