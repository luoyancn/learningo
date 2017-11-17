package deploy

import (
	"k8sdeploy/utils"
)

func PrepareK8SBinary(nodes ...string) bool {
	return utils.SCPBinary("k8s", nodes...)
}
