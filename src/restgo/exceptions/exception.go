package exceptions

import (
	"fmt"
)

type NotFoundException struct {
	ResourceType string
	ResourceId   string
}

func (self NotFoundException) Error() string {
	return fmt.Sprintf("%s with %s cannot be found",
		self.ResourceType, self.ResourceId)
}

func NewNotFoundException(resource_type string,
	resource_id string) NotFoundException {
	return NotFoundException{
		ResourceType: resource_type,
		ResourceId:   resource_id,
	}
}
