package bleeding_edge

import (
	"github.com/tracehawk/tracehawkx/modules"
)

func init() {
	// Register all bleeding-edge modules
	modules.Register(&LLMFuzzerModule{})
	modules.Register(&AutoPatchModule{})
	modules.Register(&ShadowCloneModule{})
	modules.Register(&DepDriftModule{})
	modules.Register(&TimingMapModule{})
	modules.Register(&BlueTeamModule{})
}
