package consul

import (
	"fmt"

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

		DeprecationMessage: `The /operator/license has been removed in Consul v1.10.0 and this resource will be removed in a future version of the Terraform provider.`,

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
	client, _, wOpts := getClient(d, meta)
	operator := client.Operator()

	license := d.Get("license").(string)

	_, err := operator.LicensePut(license, wOpts)
	if err != nil {
		return fmt.Errorf("failed to set license: %v", err)
	}

	return resourceConsulLicenseRead(d, meta)
}

func resourceConsulLicenseRead(d *schema.ResourceData, meta interface{}) error {
	client, qOpts, _ := getClient(d, meta)
	operator := client.Operator()

	licenseReply, err := operator.LicenseGet(qOpts)
	if err != nil {
		return fmt.Errorf("failed to read license: %v", err)
	}

	d.SetId(licenseReply.License.LicenseID)

	sw := newStateWriter(d)
	sw.set("valid", licenseReply.Valid)
	sw.set("license_id", licenseReply.License.LicenseID)
	sw.set("customer_id", licenseReply.License.CustomerID)
	sw.set("installation_id", licenseReply.License.InstallationID)
	sw.set("issue_time", licenseReply.License.IssueTime.String())
	sw.set("start_time", licenseReply.License.StartTime.String())
	sw.set("expiration_time", licenseReply.License.ExpirationTime.String())
	sw.set("product", licenseReply.License.Product)
	sw.set("features", licenseReply.License.Features)
	sw.set("warnings", licenseReply.Warnings)

	return sw.error()
}

func resourceConsulLicenseDelete(d *schema.ResourceData, meta interface{}) error {
	client, _, wOpts := getClient(d, meta)
	operator := client.Operator()

	_, err := operator.LicenseReset(wOpts)
	if err != nil {
		return fmt.Errorf("failed to remove license: %v", err)
	}

	d.SetId("")
	return nil
}
