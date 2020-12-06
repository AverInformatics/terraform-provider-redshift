# <resource name> Resource/Data Source

Description of what this resource does, with links to official
app/service documentation.

## Example Usage

```hcl
resource redshift_user test_user {
  username = "testusernew" # User name are not immutable.
  password = "Testpass123" # You can pass an md5 encrypted password here by prefixing the hash with md5
  valid_until = "2018-10-30" # See below for an example with 'password_disabled'
  connection_limit = "4"
  create_db = true
  syslog_access = "UNRESTRICTED"
  superuser = true
}
```

## Argument Reference

* `username` - (Required) redshift username to create in SQL
* `password` - (Required) redshift password to create in SQL, You can pass an md5 encrypted password here by prefixing the hash with md5
* `valid_until` - (Optional) date the user credentials expire
* `connection_limit` - (Optional) Number of connections provider can establish to redshift, defaults to UNLIMITED
* `password_disabled` - (Optional) defaults to false
* `create_db` - (Optional) Allows user to create new database defaults to false
* `syslog_access` - (Optional)  allow user access to redshift system logs, default to RESTRICTED
