package consul

import (
	"fmt"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceConsulLicense() *schema.Resource {
	return &schema.Resource{
		Create: resourceConsulLicenseCreate,
		Read:   resourceConsulLicenseRead,
		Update: resourceConsulLicenseCreate,
		Delete: resourceConsulLicenseDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			// Input
			"datacenter": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"license": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			// Output
			"valid": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"license_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"installation_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issue_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiration_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flags": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"features": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"warnings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceConsulLicenseCreate(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	operator := client.Operator()

	license := d.Get("license").(string)
	datacenter, err := getDC(d, client, meta)
	if err != nil {
		return fmt.Errorf("failed to read datacenter: %v", err)
	}

	wOpts := &consulapi.WriteOptions{
		Datacenter: datacenter,
	}

	_, err = operator.LicensePut(license, wOpts)
	if err != nil {
		return fmt.Errorf("failed to set license: %v", err)
	}

	return resourceConsulLicenseRead(d, meta)
}

func resourceConsulLicenseRead(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	operator := client.Operator()

	datacenter, err := getDC(d, client, meta)
	if err != nil {
		return fmt.Errorf("failed to read datacenter: %v", err)
	}

	qOpts := &consulapi.QueryOptions{
		Datacenter: datacenter,
	}

	licenseReply, err := operator.LicenseGet(qOpts)
	if err != nil {
		return fmt.Errorf("failed to read license: %v", err)
	}

	d.SetId(licenseReply.License.LicenseID)

	if err = d.Set("valid", licenseReply.Valid); err != nil {
		return fmt.Errorf("failed to set 'valid': %v", err)
	}
	if err = d.Set("license_id", licenseReply.License.LicenseID); err != nil {
		return fmt.Errorf("failed to set 'license_id': %v", err)
	}
	if err = d.Set("customer_id", licenseReply.License.CustomerID); err != nil {
		return fmt.Errorf("failed to set 'customer_id': %v", err)
	}
	if err = d.Set("installation_id", licenseReply.License.InstallationID); err != nil {
		return fmt.Errorf("failed to set 'installation_id': %v", err)
	}
	if err = d.Set("issue_time", licenseReply.License.IssueTime.String()); err != nil {
		return fmt.Errorf("failed to set 'issue_time': %v", err)
	}
	if err = d.Set("start_time", licenseReply.License.StartTime.String()); err != nil {
		return fmt.Errorf("failed to set 'start_time': %v", err)
	}
	if err = d.Set("expiration_time", licenseReply.License.ExpirationTime.String()); err != nil {
		return fmt.Errorf("failed to set 'expiration_time': %v", err)
	}
	if err = d.Set("product", licenseReply.License.Product); err != nil {
		return fmt.Errorf("failed to set 'product': %v", err)
	}
	if err = d.Set("flags", licenseReply.License.Flags); err != nil {
		return fmt.Errorf("failed to set 'flags': %v", err)
	}
	if err = d.Set("features", licenseReply.License.Features); err != nil {
		return fmt.Errorf("failed to set 'features': %v", err)
	}
	if err = d.Set("warnings", licenseReply.Warnings); err != nil {
		return fmt.Errorf("failed to set 'warnings': %v", err)
	}

	return nil
}

func resourceConsulLicenseDelete(d *schema.ResourceData, meta interface{}) error {
	client := getClient(meta)
	operator := client.Operator()

	datacenter, err := getDC(d, client, meta)
	if err != nil {
		return fmt.Errorf("failed to read datacenter: %v", err)
	}

	wOpts := &consulapi.WriteOptions{
		Datacenter: datacenter,
	}

	_, err = operator.LicenseReset(wOpts)
	if err != nil {
		return fmt.Errorf("failed to remove license: %v", err)
	}

	d.SetId("")
	return nil
}
