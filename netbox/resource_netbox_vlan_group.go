package netbox

import (
	"strconv"

	"github.com/fbreckle/go-netbox/netbox/client"
	"github.com/fbreckle/go-netbox/netbox/client/ipam"
	"github.com/fbreckle/go-netbox/netbox/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetboxVlanGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetboxVlanGroupCreate,
		Read:   resourceNetboxVlanGroupRead,
		Update: resourceNetboxVlanGroupUpdate,
		Delete: resourceNetboxVlanGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"slug": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.StringLenBetween(0, 30),
			},
			"scope_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"scope_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceNetboxVlanGroupCreate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	
	name := d.Get("name").(string)

	scopeid := int64(d.Get("scope_id").(int))
	scopetype := d.Get("scope_type").(string)

	slugValue, slugOk := d.GetOk("slug")
	var slug string
	// Default slug to name attribute if not given
	if !slugOk {
		slug = name
	} else {
		slug = slugValue.(string)
	}

	description := d.Get("description").(string)

	data := models.VLANGroup{}

	data.Name = &name
	data.ScopeID = &scopeid
	data.ScopeType = scopetype
	data.Description = description
	data.Slug = &slug




	params := ipam.NewIpamVlanGroupsCreateParams().WithData(&data)
	res, err := api.Ipam.IpamVlanGroupsCreate(params, nil)
	if err != nil {
		return err
	}
	d.SetId(strconv.FormatInt(res.GetPayload().ID, 10))

	return resourceNetboxVlanGroupUpdate(d, m)
}

func resourceNetboxVlanGroupRead(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := ipam.NewIpamVlanGroupsReadParams().WithID(id)

	res, err := api.Ipam.IpamVlanGroupsRead(params, nil)
	if err != nil {
		return err
		errorcode := err.(*ipam.IpamVlanGroupsReadDefault).Code()
		if errorcode == 404 {
			// If the ID is updated to blank, this tells Terraform the resource no longer exists (maybe it was destroyed out of band). Just like the destroy callback, the Read function should gracefully handle this case. https://www.terraform.io/docs/extend/writing-custom-providers.html
			d.SetId("")
			return nil
		}
		return err
	}

	if res.GetPayload().Name != nil {
		d.Set("name", res.GetPayload().Name)
	}

	if res.GetPayload().ScopeID != nil {
		d.Set("scope_id", res.GetPayload().ScopeID)
	}

	if res.GetPayload().ScopeType != "" {
		d.Set("scope_type", res.GetPayload().ScopeType)
	}

	if res.GetPayload().Description != "" {
		d.Set("description", res.GetPayload().Description)
	}

	if res.GetPayload().Slug != nil {
		d.Set("slug", res.GetPayload().Slug)
	}


	// if res.GetPayload().Site != nil {
	// 	d.Set("site_id", res.GetPayload().Site.ID)
	// }

	return nil
}

func resourceNetboxVlanGroupUpdate(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	data := models.VLANGroup{}
	name := d.Get("name").(string)
	scopeid := int64(d.Get("scope_id").(int))
	description := d.Get("description").(string)

	data.Name = &name

	slugValue, slugOk := d.GetOk("slug")
	var slug string
	// Default slug to name attribute if not given
	if !slugOk {
		slug = name
	} else {
		slug = slugValue.(string)
	}
	data.Slug = &slug

	data.ScopeID = &scopeid
	data.ScopeType = d.Get("scope_type").(string)

	data.Description = description

	// if siteID, ok := d.GetOk("site_id"); ok {
	// 	data.Site = int64ToPtr(int64(siteID.(int)))
	// }

	
	params := ipam.NewIpamVlanGroupsUpdateParams().WithID(id).WithData(&data)
	_, err := api.Ipam.IpamVlanGroupsUpdate(params, nil)
	if err != nil {
		return err
	}
	return resourceNetboxVlanGroupRead(d, m)
}

func resourceNetboxVlanGroupDelete(d *schema.ResourceData, m interface{}) error {
	api := m.(*client.NetBoxAPI)
	id, _ := strconv.ParseInt(d.Id(), 10, 64)
	params := ipam.NewIpamVlansDeleteParams().WithID(id)
	_, err := api.Ipam.IpamVlansDelete(params, nil)
	if err != nil {
		return err
	}

	return nil
}
