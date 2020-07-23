package okta

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/okta/okta-sdk-golang/v2/okta/query"
)

func dataSourceAuthServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAuthServerRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"audiences": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"kid": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials_last_rotated": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials_next_rotation": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"credentials_rotation_mode": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAuthServerRead(d *schema.ResourceData, m interface{}) error {
	name := d.Get("name").(string)
	client := getOktaClientFromMetadata(m)
	context := getOktaContextFromMetadata(m)

	// authServer, err := getSupplementFromMetadata(m).FindAuthServer(name, &query.Params{})

	authServers, _, err := client.AuthorizationServer.ListAuthorizationServers(context, &query.Params{})
	if err != nil {
		return err
	}
	if authServers == nil {
		return fmt.Errorf("No authorization server found with provided name %s", name)
	}

	if err != nil {
		return err
	}

	for _, authServer := range authServers {
		if authServer.Name == name {
			d.SetId(authServer.Id)
			d.Set("name", authServer.Name)
			d.Set("description", authServer.Description)
			d.Set("audiences", convertStringSetToInterface(authServer.Audiences))
			d.Set("credentials_rotation_mode", authServer.Credentials.Signing.RotationMode)
			d.Set("kid", authServer.Credentials.Signing.Kid)
			d.Set("credentials_next_rotation", authServer.Credentials.Signing.NextRotation.String())
			d.Set("credentials_last_rotated", authServer.Credentials.Signing.LastRotated.String())
			d.Set("status", authServer.Status)
		}
	}

	if len(d.Id()) == 0 {
		return fmt.Errorf("No authorization server found with provided name %s", name)
	}

	return nil
}
