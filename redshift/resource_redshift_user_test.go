package redshift

import (
	"database/sql"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"log"
	"strconv"
	"testing"
)

func TestAccRedshiftUser_basic(t *testing.T) {
	rInt := acctest.RandInt()
	var user testAccRedshiftUserExpectedAttributes

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create User
			{
				Config: testAccRedshiftUserConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckRedShiftUserExists("redshift_user.test_user", &user),
					testAccDataSourceRedshiftUserAttributes(&user, &testAccRedshiftUserExpectedAttributes{
						Username: fmt.Sprintf("test_user_%d", rInt),
						ValidUntil: "",
						PasswordDisabled: false,
						CreateDB: false,
						ConnectionLimit: "UNLIMITED",
						SyslogAccess: "RESTRICTED",
						SuperUser: false,
					}),
				),
			},
		},
	})
}

func testAccDataSourceRedshiftUserAttributes(src, n *testAccRedshiftUserExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		user := s.RootModule().Resources[src]
		userResource := user.Primary.Attributes

		search := s.RootModule().Resources[n]
		searchResource := search.Primary.Attributes

		testAttributes := []string{
			"username",
		}

		for _, attribute := range testAttributes {
			if searchResource[attribute] != userResource[attribute] {
				return fmt.Errorf("Expected user's parameter `%s` to be: %s, but got: `%s`", attribute, userResource[attribute], searchResource[attribute])
			}
		}

		return nil
	}
}

func testAccRedshiftUserConfig(rInt int) string {
	return fmt.Sprintf(`
resource redshift_user test_user {
  username = "test_user_%d"
  password = "Testpass123"
}
  `, rInt)
}

type testAccRedshiftUserExpectedAttributes struct {
	Username         string
	ValidUntil       string
	PasswordDisabled bool
	CreateDB         bool
	ConnectionLimit  string
	SyslogAccess     string
	SuperUser        bool
	Usesysid         string
}

func testAccCheckRedShiftUserExists(n string, user *testAccRedshiftUserExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}
		userID := rs.Primary.Attributes["usesysid"]
		if userID == "" {
			return fmt.Errorf("No usesysid is set")
		}
		redshiftClient := testAccProvider.Meta().(*Client).db

		tx, txErr := redshiftClient.Begin()

		if txErr != nil {
			panic(txErr)
		}

		id, _ := strconv.Atoi(userID)

		var readUserQuery = "select usename, usecreatedb, usesuper, valuntil, useconnlimit " +
			"from pg_user_info where usesysid = $1"

		log.Print("Reading redshift user with query: " + readUserQuery)

		err := tx.QueryRow(readUserQuery, d.Id()).Scan(&user)

		if err != nil {
			log.Print("Reading user does not exist")
			log.Print(err)
			return err
		}

		if err != nil {
			tx.Rollback()
			return err
		}

		tx.Commit()

		gotUser, _, err := conn.Users.GetUser(id)
		if err != nil {
			return err
		}
		*user = *gotUser
		return nil
	}
}

func testAccCheckGitlabUserAttributes(user *gitlab.User, want *testAccGitlabUserExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		return nil
	}
}