package resource

import (
	"cli/pkg/manifest"
	"cli/pkg/util"
	"encoding/json"
	"fmt"

	"cli/pkg/config"
	"types"
)

func ApplyResource(manifest *manifest.Manifest) error {

	apiEndpoint := config.GetAPIEndpoint()

	resources, err := config.GetDiscoveryResources()
	if err != nil {
		return err
	}

	var resourceConfig *types.APIResource
	for _, resource := range resources {
		if resource.Kind == manifest.Kind {
			resourceConfig = &resource
			break
		}
	}

	if resourceConfig.Kind == "" {
		return fmt.Errorf("resource %s not found", manifest.Kind)
	}

	return sendRequest(apiEndpoint, resourceConfig.SingularName, manifest)
}

func sendRequest(apiEndpoint string, singularName string, manifest *manifest.Manifest) error {
	body, err := json.Marshal(manifest.Spec)
	if err != nil {
		return err
	}
	resp, err := util.SendRequest("POST", apiEndpoint+"/"+singularName, body)
	if err != nil {
		return err
	}
	fmt.Println(string(resp))
	fmt.Println("Response: ", string(resp))
	return nil
}
