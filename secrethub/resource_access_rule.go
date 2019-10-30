package secrethub

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/secrethub/secrethub-go/internals/api"
)

func resourceAccessRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceAccessRuleCreate,
		Read:   resourceAccessRuleRead,
		Update: resourceAccessRuleUpdate,
		Delete: resourceAccessRuleDelete,
		Schema: map[string]*schema.Schema{
			"dir": {
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
		Importer: &schema.ResourceImporter{
			State: resourceAccessRuleImport,
		},
	}
}

func resourceAccessRuleCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Get("dir").(string)
	permission := d.Get("permission").(string)
	account := d.Get("account_name").(string)

	_, err := client.AccessRules().Get(path, account)
	if err == nil {
		return fmt.Errorf("access rule already exists: %s:%s", path, account)
	} else if err != api.ErrAccessRuleNotFound {
		return err
	}

	_, err = client.AccessRules().Set(path, permission, account)
	if err != nil {
		return err
	}

	d.SetId(path + ":" + account)

	return resourceAccessRuleRead(d, m)
}

func resourceAccessRuleUpdate(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Get("dir").(string)
	permission := d.Get("permission").(string)
	account := d.Get("account_name").(string)

	_, err := client.AccessRules().Set(path, permission, account)
	if err != nil {
		return err
	}

	return resourceAccessRuleRead(d, m)
}

func resourceAccessRuleRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path, account, err := resourceAccessRuleParseID(d.Id())
	if err != nil {
		return err
	}

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

	path := d.Get("dir").(string)
	account := d.Get("account_name").(string)

	return client.AccessRules().Delete(path, account)
}

func resourceAccessRuleParseID(id string) (string, string, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed ID: %s is not a valid access rule ID, expected <path>:<account_name>", id)
	}
	path := parts[0]
	account := parts[1]

	return path, account, nil
}

func resourceAccessRuleImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	path, account, err := resourceAccessRuleParseID(d.Id())
	if err != nil {
		return nil, err
	}

	err = d.Set("dir", path)
	if err != nil {
		return nil, err
	}

	err = d.Set("account_name", account)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
