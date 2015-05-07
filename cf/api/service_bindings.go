package api

import (
	"fmt"
	"strings"

	"github.com/cloudfoundry/cli/cf/configuration/core_config"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/cf/net"
)

type ServiceBindingRepository interface {
	Create(instanceGuid, appGuid string) (apiErr error)
	Delete(instance models.ServiceInstance, appGuid string) (found bool, apiErr error)
}

type CloudControllerServiceBindingRepository struct {
	config  core_config.Reader
	gateway net.Gateway
}

func NewCloudControllerServiceBindingRepository(config core_config.Reader, gateway net.Gateway) (repo CloudControllerServiceBindingRepository) {
	repo.config = config
	repo.gateway = gateway
	return
}

func (repo CloudControllerServiceBindingRepository) Create(instanceGuid, appGuid string) (apiErr error) {
	path := "/v2/service_bindings"
	body := fmt.Sprintf(
		`{"app_guid":"%s","service_instance_guid":"%s"}`,
		appGuid, instanceGuid,
	)
	return repo.gateway.CreateResource(repo.config.ApiEndpoint(), path, strings.NewReader(body))
}

func (repo CloudControllerServiceBindingRepository) Delete(instance models.ServiceInstance, appGuid string) (found bool, apiErr error) {
	var path string

	for _, binding := range instance.ServiceBindings {
		if binding.AppGuid == appGuid {
			path = binding.Url
			break
		}
	}

	if path == "" {
		return
	} else {
		found = true
	}

	apiErr = repo.gateway.DeleteResource(repo.config.ApiEndpoint(), path)
	return
}
