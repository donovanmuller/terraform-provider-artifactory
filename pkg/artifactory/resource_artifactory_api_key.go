package artifactory

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	v1 "github.com/atlassian/go-artifactory/v2/artifactory/v1"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceArtifactoryApiKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiKeyCreate,
		Read:   resourceApiKeyRead,
		Delete: resourceApiKeyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"api_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func unpackApiKey(s *schema.ResourceData) *v1.ApiKey {
	d := &ResourceData{s}
	apiKey := new(v1.ApiKey)
	apiKey.ApiKey = d.getStringRef("api_key", false)

	return apiKey
}

func packApiKey(apiKey *v1.ApiKey, d *schema.ResourceData) error {
	hasErr := false
	logErr := cascadingErr(&hasErr)

	logErr(d.Set("api_key", apiKey.ApiKey))

	if hasErr {
		return fmt.Errorf("failed to pack api key")
	}

	return nil
}

func resourceApiKeyCreate(d *schema.ResourceData, m interface{}) error {
	c := m.(*ArtClient).ArtOld

	apiKey, _, err := c.V1.Security.CreateApiKey(context.Background())
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(hashcode.String(*apiKey.ApiKey)))
	return resourceApiKeyRead(d, m)
}

func resourceApiKeyRead(d *schema.ResourceData, m interface{}) error {
	c := m.(*ArtClient).ArtOld

	apiKey, resp, err := c.V1.Security.GetApiKey(context.Background())
	if resp.StatusCode == http.StatusNotFound {
		d.SetId("")
		return nil
	} else if err != nil {
		return err
	}

	return packApiKey(apiKey, d)
}

func resourceApiKeyDelete(d *schema.ResourceData, m interface{}) error {
	c := m.(*ArtClient).ArtOld
	_, resp, err := c.V1.Security.RevokeApiKey(context.Background())
	if resp.StatusCode == http.StatusNotFound {
		return nil
	}

	return err
}
