package secrethub

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/secrethub/secrethub-go/internals/api"
	"github.com/secrethub/secrethub-go/pkg/randchar"
)

func resourceSecret() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecretCreate,
		Read:   resourceSecretRead,
		Update: resourceSecretUpdate,
		Delete: resourceSecretDelete,
		Importer: &schema.ResourceImporter{
			State: resourceSecretImport,
		},
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The path where the secret will be stored.",
			},
			"version": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The version of the secret.",
			},
			"value": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"generate"},
				Description:   "The secret contents. Either `value` or `generate` must be defined.",
			},
			"generate": {
				Type:          schema.TypeList,
				Optional:      true,
				MaxItems:      1,
				Description:   "Settings for autogenerating a secret. Either `value` or `generate` must be defined.",
				ConflictsWith: []string{"value"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"length": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The length of the secret to generate.",
						},
						"use_symbols": {
							Type:        schema.TypeBool,
							Deprecated:  "use the charsets attribute instead",
							Optional:    true,
							Description: "Whether the secret should contain symbols.",
						},
						"charsets": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Define the set of characters to randomly generate a password from. Options are all, alphanumeric, numeric, lowercase, uppercase, letters, symbols and human-readable.",
						},
						"min": {
							Type:        schema.TypeMap,
							Optional:    true,
							Elem:        &schema.Schema{Type: schema.TypeInt},
							Description: "Ensure that the generated secret contains at least n characters from the given character set. Note that adding constraints reduces the strength of the secret.",
						},
					},
				},
			},
		},
	}
}

func resourceSecretCreate(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	valueStr := d.Get("value").(string)
	generateList := d.Get("generate").([]interface{})
	if valueStr == "" && len(generateList) == 0 {
		return fmt.Errorf("either 'value' or 'generate' must be specified")
	}

	var value []byte
	if valueStr != "" {
		value = []byte(valueStr)
	}

	if len(generateList) > 0 {
		settings := generateList[0].(map[string]interface{})
		useSymbols := settings["use_symbols"].(bool)
		length := settings["length"].(int)
		charsetSet := settings["charsets"].(*schema.Set)
		charsets := charsetSet.List()
		charset := randchar.Charset{}
		if len(charsets) == 0 {
			charset = randchar.Alphanumeric
		}
		if useSymbols {
			charset = charset.Add(randchar.Symbols)
		}
		for _, charsetName := range charsets {
			set, found := randchar.CharsetByName(charsetName.(string))
			if !found {
				return fmt.Errorf("unknown charset: %s", charsetName)
			}
			charset = charset.Add(set)
		}

		minRuleMap := settings["min"].(map[string]interface{})
		var minRules []randchar.Option
		for charset, min := range minRuleMap {
			n := min.(int)
			set, found := randchar.CharsetByName(charset)
			if !found {
				return fmt.Errorf("unknown charset: %s", charset)
			}
			minRules = append(minRules, randchar.Min(n, set))
		}

		var err error
		rand, err := randchar.NewRand(charset, minRules...)
		if err != nil {
			return err
		}
		value, err = rand.Generate(length)
		if err != nil {
			return err
		}
	}

	path := d.Get("path").(string)

	res, err := client.Secrets().Write(path, value)
	if err != nil {
		return err
	}

	d.SetId(string(path))
	err = d.Set("value", string(value))
	if err != nil {
		return err
	}
	err = d.Set("version", res.Version)
	if err != nil {
		return err
	}

	return resourceSecretRead(d, m)
}

func resourceSecretRead(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Id()

	remote, err := client.Secrets().Get(path)
	if err == api.ErrSecretNotFound {
		// The secret was deleted outside of the current Terraform workspace, so invalidate this resource
		d.SetId("")
		return nil
	}
	if err != nil {
		return err
	}

	prev := d.Get("version")
	if prev != remote.LatestVersion {
		// The secret has been updated outside of the current Terraform workspace, so the new secret version has to be fetched
		updated, err := client.Secrets().Versions().GetWithData(path)
		if err != nil {
			return err
		}
		err = d.Set("value", string(updated.Data))
		if err != nil {
			return err
		}
		err = d.Set("version", updated.Version)
		if err != nil {
			return err
		}
	}

	return nil
}

func resourceSecretUpdate(d *schema.ResourceData, m interface{}) error {
	return resourceSecretCreate(d, m)
}

func resourceSecretDelete(d *schema.ResourceData, m interface{}) error {
	provider := m.(providerMeta)
	client := *provider.client

	path := d.Id()

	return client.Secrets().Delete(path)
}

func resourceSecretImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	path := d.Id()

	err := api.ValidateSecretPath(path)
	if err != nil {
		return nil, err
	}

	err = d.Set("path", path)
	if err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}
