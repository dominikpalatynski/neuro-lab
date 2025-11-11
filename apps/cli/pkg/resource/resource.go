package resource

import (
	"cli/pkg/manifest"
	"cli/pkg/util"
	"encoding/json"
	"fmt"

	"cli/pkg/config"
	"types"

	apierrors "github.com/neuro-lab/errors"
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

	fmt.Println("Applying ", manifest.Kind)
	resp, err := util.SendRequest("POST", apiEndpoint+"/"+singularName, body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		var errorResponse apierrors.ErrorResponse
		if err := json.Unmarshal(resp.Body, &errorResponse); err != nil {
			return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, resp.Body)
		}

		if errorResponse.Type == apierrors.TypeValidationFailed {
			return fmt.Errorf("%s: %s\n%v", errorResponse.Title, errorResponse.Detail, errorResponse.Errors)
		}

		return fmt.Errorf("%s: %s", errorResponse.Title, errorResponse.Detail)
	}

	// Success: status code is 2xx
	fmt.Println("Resource applied successfully")
	return nil
}
