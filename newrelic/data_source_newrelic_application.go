package newrelic

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	newrelic "github.com/paultyng/go-newrelic/v4/api"
)

func dataSourceNewRelicApplication() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNewRelicApplicationRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instance_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
			},
			"host_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
			},
		},
	}
}

func dataSourceNewRelicApplicationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ProviderConfig).Client

	log.Printf("[INFO] Reading New Relic applications")

	name := d.Get("name").(string)
	filter := newrelic.ApplicationsFilters{Name: &name}
	applications, err := client.QueryApplications(filter)
	if err != nil {
		return err
	}

	var application *newrelic.Application

	for _, a := range applications {
		if a.Name == name {
			application = &a
			break
		}
	}

	if application == nil {
		return fmt.Errorf("the name '%s' does not match any New Relic applications", name)
	}

	d.SetId(strconv.Itoa(application.ID))
	d.Set("name", application.Name)
	d.Set("instance_ids", application.Links.InstanceIDs)
	d.Set("host_ids", application.Links.HostIDs)

	return nil
}
