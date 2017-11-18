package deploy

import (
	"k8sdeploy/utils"
)

func CreateCA() error {
	return utils.InitCA()
}
