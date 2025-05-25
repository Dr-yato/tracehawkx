package stable

import (
	"github.com/tracehawk/tracehawkx/modules"
)

func init() {
	// Register all stable modules
	modules.Register(&SubfinderModule{})
	modules.Register(&NmapModule{})
	modules.Register(&HttpxModule{})
	modules.Register(&NucleiModule{})
}
