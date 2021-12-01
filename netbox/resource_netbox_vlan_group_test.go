package netbox

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/fbreckle/go-netbox/netbox/client"
	"github.com/fbreckle/go-netbox/netbox/client/ipam"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccNetboxVlanGroupFullDependencies(testName string) string {
	return fmt.Sprintf(`
resource "netbox_tag" "test" {
  name = "%[1]s"
}

resource "netbox_tenant" "test" {
  name = "%[1]s"
}

resource "netbox_site" "test" {
  name = "%[1]s"
  status = "active"
}
`, testName)
}
func TestAccNetboxVlanGroup_basic(t *testing.T) {

	testSlug := "vlan_basic"
	testName := testAccGetTestName(testSlug)
	testDescription := "Test Description"
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccNetboxVlanFullDependencies(testName) + fmt.Sprintf(`
resource "netbox_vlan_group" "test_basic" {
  name = "%s"
  slug = "%s"
  description = "%s"
}`, testName, testSlug, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan_group.test_basic", "name", testName),
					resource.TestCheckResourceAttr("netbox_vlan_group.test_basic", "slug", testSlug),
					resource.TestCheckResourceAttr("netbox_vlan_group.test_basic", "description", testDescription),
				),
			},
			{
				ResourceName:      "netbox_vlan_group.test_basic",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetboxVlanGroup_with_dependencies(t *testing.T) {

	testSlug := "vlan_group_with_dependencies"
	testName := testAccGetTestName(testSlug)
	testDescription := "Test Description"
	resource.ParallelTest(t, resource.TestCase{
		Providers: testAccProviders,
		PreCheck:  func() { testAccPreCheck(t) },
		Steps: []resource.TestStep{
			{
				Config: testAccNetboxVlanFullDependencies(testName) + fmt.Sprintf(`
resource "netbox_vlan_group" "test_with_dependencies" {
  name = "%s"
  slug = "%s"
  description = "%s"
}`, testName, testSlug, testDescription),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("netbox_vlan_group.test_with_dependencies", "name", testName),
					resource.TestCheckResourceAttr("netbox_vlan_group.test_with_dependencies", "slug", testSlug),
					resource.TestCheckResourceAttr("netbox_vlan_group.test_with_dependencies", "description", testDescription),
					// resource.TestCheckResourceAttr("netbox_vlan.test_with_dependencies", "status", "active"),
					// resource.TestCheckResourceAttrPair("netbox_vlan.test_with_dependencies", "tenant_id", "netbox_tenant.test", "id"),
					// resource.TestCheckResourceAttrPair("netbox_vlan.test_with_dependencies", "site_id", "netbox_site.test", "id"),
					// resource.TestCheckResourceAttr("netbox_vlan.test_with_dependencies", "tags.#", "1"),
					// resource.TestCheckResourceAttr("netbox_vlan.test_with_dependencies", "tags.0", testName),
				),
			},
			{
				ResourceName:      "netbox_vlan_group.test_with_dependencies",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func init() {
	resource.AddTestSweepers("netbox_vlan_group", &resource.Sweeper{
		Name:         "netbox_vlan_group",
		Dependencies: []string{},
		F: func(region string) error {
			m, err := sharedClientForRegion(region)
			if err != nil {
				return fmt.Errorf("Error getting client: %s", err)
			}
			api := m.(*client.NetBoxAPI)
			params := ipam.NewIpamVlansListParams()
			res, err := api.Ipam.IpamVlansList(params, nil)
			if err != nil {
				return err
			}
			for _, vlan := range res.GetPayload().Results {
				if strings.HasPrefix(*vlan.Name, testPrefix) {
					deleteParams := ipam.NewIpamVlansDeleteParams().WithID(vlan.ID)
					_, err := api.Ipam.IpamVlansDelete(deleteParams, nil)
					if err != nil {
						return err
					}
					log.Print("[DEBUG] Deleted a vlan")
				}
			}
			return nil
		},
	})
}
