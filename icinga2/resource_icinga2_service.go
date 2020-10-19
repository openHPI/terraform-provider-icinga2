package icinga2

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/lrsmith/go-icinga2-api/iapi"
)

func resourceIcinga2Service() *schema.Resource {

	return &schema.Resource{
		Create: resourceIcinga2ServiceCreate,
		Exists: resourceIcinga2ServiceExists,
		Read:   resourceIcinga2ServiceRead,
		Delete: resourceIcinga2ServiceDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "ServiceName",
				ForceNew:    true,
			},
			"hostname": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Hostname",
				ForceNew:    true,
			},
			"check_command": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "CheckCommand",
				ForceNew:    true,
			},
			"vars": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"templates": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"zone": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Zone",
			},
		},
	}
}

func resourceIcinga2ServiceCreate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)

	hostname := d.Get("hostname").(string)
	name := d.Get("name").(string)

	var attrs iapi.ServiceAttrs
	attrs.CheckCommand = d.Get("check_command").(string)

	attrs.Vars = d.Get("vars").(map[string]interface{})

	attrs.Templates = make([]string, len(d.Get("templates").([]interface{})))
	for i, v := range d.Get("templates").([]interface{}) {
		attrs.Templates[i] = v.(string)
	}
	attrs.Zone = d.Get("zone").(string)

	services, err := client.CreateService(name, hostname, attrs)
	if err != nil {
		return err
	}

	found := false
	for _, service := range services {
		if service.Name == hostname+"!"+name {
			d.SetId(hostname + "!" + name)
			found = true
		}
	}

	if !found {
		return fmt.Errorf("Failed to Create Service %s!%s : %s", hostname, name, err)
	}

	return nil

}

func resourceIcinga2ServiceRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)

	hostname := d.Get("hostname").(string)
	name := d.Get("name").(string)

	services, err := client.GetService(name, hostname)
	if err != nil {
		return err
	}

	for _, service := range services {
		if service.Name == hostname+"!"+name {
			d.SetId(hostname + "!" + name)
			d.Set("hostname", hostname)
			d.Set("check_command", service.Attrs.CheckCommand)
			d.Set("vars", service.Attrs.Vars)
			d.Set("zone", service.Attrs.Zone)
		}
	}

	return nil
}

func resourceIcinga2ServiceExists(d *schema.ResourceData, meta interface{}) (bool, error) {

	client := meta.(*iapi.Server)

	hostname := d.Get("hostname").(string)
	name := d.Get("name").(string)

	services, err := client.GetService(name, hostname)
	if err != nil {
		return false, err
	}

	for _, service := range services {
		if service.Name == hostname+"!"+name {
			return true, nil
		}
	}

	return false, nil
}

func resourceIcinga2ServiceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*iapi.Server)

	hostname := d.Get("hostname").(string)
	name := d.Get("name").(string)

	err := client.DeleteService(name, hostname)
	if err != nil {
		return fmt.Errorf("Failed to Delete Service %s!%s : %s", hostname, name, err)
	}

	return nil
}
