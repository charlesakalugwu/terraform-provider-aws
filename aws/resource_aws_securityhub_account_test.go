package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/securityhub"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/terraform-providers/terraform-provider-aws/atest"
	awsprovider "github.com/terraform-providers/terraform-provider-aws/provider"
)

func testAccAWSSecurityHubAccount_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { atest.PreCheck(t) },
		ErrorCheck:   atest.ErrorCheck(t, securityhub.EndpointsID),
		Providers:    atest.Providers,
		CheckDestroy: testAccCheckAWSSecurityHubAccountDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAWSSecurityHubAccountConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSSecurityHubAccountExists("aws_securityhub_account.example"),
				),
			},
			{
				ResourceName:      "aws_securityhub_account.example",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckAWSSecurityHubAccountExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := atest.Provider.Meta().(*awsprovider.AWSClient).SecurityHubConn

		_, err := conn.GetEnabledStandards(&securityhub.GetEnabledStandardsInput{})

		if err != nil {
			// Can only read enabled standards if Security Hub is enabled
			if tfawserr.ErrMessageContains(err, "InvalidAccessException", "not subscribed to AWS Security Hub") {
				return fmt.Errorf("Security Hub account not found")
			}
			return err
		}

		return nil
	}
}

func testAccCheckAWSSecurityHubAccountDestroy(s *terraform.State) error {
	conn := atest.Provider.Meta().(*awsprovider.AWSClient).SecurityHubConn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_securityhub_account" {
			continue
		}

		_, err := conn.GetEnabledStandards(&securityhub.GetEnabledStandardsInput{})

		if err != nil {
			// Can only read enabled standards if Security Hub is enabled
			if tfawserr.ErrMessageContains(err, "InvalidAccessException", "not subscribed to AWS Security Hub") {
				return nil
			}
			return err
		}

		return fmt.Errorf("Security Hub account still exists")
	}

	return nil
}

func testAccAWSSecurityHubAccountConfig() string {
	return `
resource "aws_securityhub_account" "example" {}
`
}
