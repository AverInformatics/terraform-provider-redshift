package redshift

import (
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccRedshiftUser_basic(t *testing.T) {
	rInt := acctest.RandInt()
	var user testAccRedshiftUserExpectedAttributes
	vu := sql.NullString{String: "false", Valid: true}
	pd := sql.NullBool{Bool: false, Valid: false}
	su := sql.NullBool{Bool:false, Valid: true}
	cd := sql.NullBool{Bool: false, Valid: true}
	cl := sql.NullString{String: "UNLIMITED", Valid: true}
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
						Username:         fmt.Sprintf("test_user_%d", rInt),
						ValidUntil:       vu,
						PasswordDisabled: pd,
						CreateDB:         cd,
						ConnectionLimit:  cl,
						SuperUser:        su,
					}),
				),
			},
		},
	})
}

func testAccDataSourceRedshiftUserAttributes(src, n *testAccRedshiftUserExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if src.ConnectionLimit != n.ConnectionLimit {
			return errors.New("connection limit not equal")
		}

		if src.PasswordDisabled != n.PasswordDisabled {
			return errors.New("password disabled not equal")
		}

		if src.ValidUntil  != n.ValidUntil {
			return errors.New("valid until not equal")
		}

		if src.SuperUser != n.SuperUser {
			return errors.New("superuser not equal")
		}

		if src.CreateDB != n.CreateDB {
			return errors.New("createdb not equal")
		}

		if src.Username != n.Username {
			return errors.New("username not equal")
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
	Username         string         `sql:"usename"`
	ValidUntil       sql.NullString `sql:"valuntil"`
	PasswordDisabled sql.NullBool   `sql:"usecatupd"`
	CreateDB         sql.NullBool   `sql:"usecreatedb"`
	ConnectionLimit  sql.NullString `sql:"useconnlimit"`
	SuperUser        sql.NullBool   `sql:"usesuper"`
	Usesysid         string         `sql:"usesysid"`
}

func testAccCheckRedShiftUserExists(n string, user *testAccRedshiftUserExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}
		userID := rs.Primary.Attributes["id"]
		if userID == "" {
			return fmt.Errorf("No usesysid is set")
		}
		redshiftClient := testAccProvider.Meta().(*Client).db

		tx, txErr := redshiftClient.Begin()

		if txErr != nil {
			panic(txErr)
		}

		id, _ := strconv.Atoi(userID)

		var readUserQuery = "select usename, usecreatedb, usesuper, usecatupd, valuntil, useconnlimit " +
			"from pg_user_info where usesysid = $1"

		err := tx.QueryRow(readUserQuery, id).Scan(&user.Username, &user.CreateDB, &user.SuperUser, &user.ValidUntil, &user.PasswordDisabled, &user.ConnectionLimit)

		if err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}
}
