package secrethub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/secrethub/secrethub-go/internals/api"
)

func resourceAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccessRuleSet,
		Read:   resourceAccessRuleRead,
		Update: resourceAccessRuleSet,
		Delete: resourceAccessRuleDelete,
		Schema: map[string]*schema.Schema{
			"dir_path": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The path of the directory on which the permission holds.",
			},
			"account_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the account (username or service ID) for which the permission holds.",
			},
			"permission": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The permission that the account has on the given directory: read, write or admin",
			},
		},
	}
}

func resourceAccessRuleSet(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Get("dir_path").(string)
	permission := d.Get("permission").(string)
	account := d.Get("account_name").(string)

	_, err := client.AccessRules().Set(path, permission, account)

	d.SetId(path + ":" + account)

	return err
}

func resourceAccessRuleRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 2 {
		return fmt.Errorf("malformed ID: %s is not a valid access rule ID, expected <path>:<account_name>", d.Id())
	}
	path := parts[0]
	account := parts[1]

	accessRule, err := client.AccessRules().Get(path, account)
	if err == api.ErrAccessRuleNotFound {
		// The secret was deleted outside of the current Terraform workspace, so invalidate this resource
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	err = d.Set("permission", accessRule.Permission.String())
	if err != nil {
		return err
	}
	return nil
}

func resourceAccessRuleDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Get("dir_path").(string)
	account := d.Get("account_name").(string)

	return client.AccessRules().Delete(path, account)
}
