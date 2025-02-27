package connect_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/connect"
	"github.com/hashicorp/aws-sdk-go-base/v2/awsv1shim/v2/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfconnect "github.com/hashicorp/terraform-provider-aws/internal/service/connect"
)

func testAccPhoneNumber_basic(t *testing.T) {
	ctx := acctest.Context(t)
	var v connect.DescribePhoneNumberOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_connect_phone_number.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, connect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPhoneNumberDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneNumberConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPhoneNumberExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "arn"),
					resource.TestCheckResourceAttr(resourceName, "country_code", "US"),
					resource.TestCheckResourceAttrSet(resourceName, "phone_number"),
					resource.TestCheckResourceAttr(resourceName, "status.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "status.0.status", connect.PhoneNumberWorkflowStatusClaimed),
					resource.TestCheckResourceAttrPair(resourceName, "target_arn", "aws_connect_instance.test", "arn"),
					resource.TestCheckResourceAttr(resourceName, "type", connect.PhoneNumberTypeDid),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPhoneNumber_description(t *testing.T) {
	ctx := acctest.Context(t)
	var v connect.DescribePhoneNumberOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	description := "example description"
	resourceName := "aws_connect_phone_number.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, connect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPhoneNumberDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneNumberConfig_description(rName, description),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPhoneNumberExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "description", description),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccPhoneNumber_prefix(t *testing.T) {
	ctx := acctest.Context(t)
	var v connect.DescribePhoneNumberOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	prefix := "+1"
	resourceName := "aws_connect_phone_number.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, connect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPhoneNumberDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneNumberConfig_prefix(rName, prefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPhoneNumberExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrSet(resourceName, "phone_number"),
					resource.TestMatchResourceAttr(resourceName, "phone_number", regexp.MustCompile(fmt.Sprintf("\\%s[0-9]{0,10}", prefix))),
					resource.TestCheckResourceAttr(resourceName, "prefix", prefix),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"prefix"},
			},
		},
	})
}

func testAccPhoneNumber_targetARN(t *testing.T) {
	ctx := acctest.Context(t)
	var v connect.DescribePhoneNumberOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	rName2 := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_connect_phone_number.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, connect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPhoneNumberDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneNumberConfig_targetARN(rName, rName2, "first"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPhoneNumberExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "target_arn", "aws_connect_instance.test", "arn"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPhoneNumberConfig_targetARN(rName, rName2, "second"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPhoneNumberExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttrPair(resourceName, "target_arn", "aws_connect_instance.test2", "arn"),
				),
			},
		},
	})
}

func testAccPhoneNumber_tags(t *testing.T) {
	ctx := acctest.Context(t)
	var v connect.DescribePhoneNumberOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_connect_phone_number.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, connect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPhoneNumberDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneNumberConfig_tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPhoneNumberExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccPhoneNumberConfig_tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPhoneNumberExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccPhoneNumberConfig_tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPhoneNumberExists(ctx, resourceName, &v),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccPhoneNumber_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	var v connect.DescribePhoneNumberOutput
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_connect_phone_number.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { acctest.PreCheck(ctx, t) },
		ErrorCheck:               acctest.ErrorCheck(t, connect.EndpointsID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckPhoneNumberDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccPhoneNumberConfig_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckPhoneNumberExists(ctx, resourceName, &v),
					acctest.CheckResourceDisappears(ctx, acctest.Provider, tfconnect.ResourcePhoneNumber(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckPhoneNumberExists(ctx context.Context, resourceName string, function *connect.DescribePhoneNumberOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Connect Phone Number not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Connect Phone Number ID not set")
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).ConnectConn()

		params := &connect.DescribePhoneNumberInput{
			PhoneNumberId: aws.String(rs.Primary.ID),
		}

		getFunction, err := conn.DescribePhoneNumberWithContext(ctx, params)
		if err != nil {
			return err
		}

		*function = *getFunction

		return nil
	}
}

func testAccCheckPhoneNumberDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_connect_phone_number" {
				continue
			}

			conn := acctest.Provider.Meta().(*conns.AWSClient).ConnectConn()

			params := &connect.DescribePhoneNumberInput{
				PhoneNumberId: aws.String(rs.Primary.ID),
			}

			_, err := conn.DescribePhoneNumberWithContext(ctx, params)

			if tfawserr.ErrCodeEquals(err, connect.ErrCodeResourceNotFoundException) {
				continue
			}

			if err != nil {
				return err
			}
		}

		return nil
	}
}

func testAccPhoneNumberConfig_base(rName string) string {
	return fmt.Sprintf(`
resource "aws_connect_instance" "test" {
  identity_management_type = "CONNECT_MANAGED"
  inbound_calls_enabled    = true
  instance_alias           = %[1]q
  outbound_calls_enabled   = true
}
`, rName)
}

func testAccPhoneNumberConfig_basic(rName string) string {
	return acctest.ConfigCompose(
		testAccPhoneNumberConfig_base(rName),
		`
resource "aws_connect_phone_number" "test" {
  target_arn   = aws_connect_instance.test.arn
  country_code = "US"
  type         = "DID"
}
`)
}

func testAccPhoneNumberConfig_description(rName, description string) string {
	return acctest.ConfigCompose(
		testAccPhoneNumberConfig_base(rName),
		fmt.Sprintf(`
resource "aws_connect_phone_number" "test" {
  target_arn   = aws_connect_instance.test.arn
  country_code = "US"
  type         = "DID"
  description  = %[1]q
}
`, description))
}

func testAccPhoneNumberConfig_prefix(rName, prefix string) string {
	return acctest.ConfigCompose(
		testAccPhoneNumberConfig_base(rName),
		fmt.Sprintf(`
resource "aws_connect_phone_number" "test" {
  target_arn   = aws_connect_instance.test.arn
  country_code = "US"
  type         = "DID"
  prefix       = %[1]q
}
`, prefix))
}

func testAccPhoneNumberConfig_targetARN(rName, rName2, selectTargetArn string) string {
	return acctest.ConfigCompose(
		testAccPhoneNumberConfig_base(rName),
		fmt.Sprintf(`
locals {
  select_target_arn = %[2]q
}

resource "aws_connect_instance" "test2" {
  identity_management_type = "CONNECT_MANAGED"
  inbound_calls_enabled    = true
  instance_alias           = %[1]q
  outbound_calls_enabled   = true
}

resource "aws_connect_phone_number" "test" {
  target_arn   = local.select_target_arn == "first" ? aws_connect_instance.test.arn : aws_connect_instance.test2.arn
  country_code = "US"
  type         = "DID"
}
`, rName2, selectTargetArn))
}

func testAccPhoneNumberConfig_tags1(rName, tag, value string) string {
	return acctest.ConfigCompose(
		testAccPhoneNumberConfig_base(rName),
		fmt.Sprintf(`
resource "aws_connect_phone_number" "test" {
  target_arn   = aws_connect_instance.test.arn
  country_code = "US"
  type         = "DID"

  tags = {
    %[1]q = %[2]q
  }
}
`, tag, value))
}

func testAccPhoneNumberConfig_tags2(rName, tag1, value1, tag2, value2 string) string {
	return acctest.ConfigCompose(
		testAccPhoneNumberConfig_base(rName),
		fmt.Sprintf(`
resource "aws_connect_phone_number" "test" {
  target_arn   = aws_connect_instance.test.arn
  country_code = "US"
  type         = "DID"

  tags = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, tag1, value1, tag2, value2))
}
