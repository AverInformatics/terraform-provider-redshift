# Redshift Provider

Manage Redshift users, groups, privileges, databases and schemas. It runs the SQL queries necessary to manage these (CREATE USER, DELETE DATABASE etc) in transactions, and also reads the state from the tables that store this state, eg pg_user_info, pg_group etc. The underlying tables are more or less equivalent to the postgres tables, but some tables are not accessible in Redshift.

Currently, supports users, groups, schemas and databases. You can set privileges for groups on schemas. Per user schema privileges will be added at a later date.

Note that schemas are the lowest level of granularity here, tables should be created by some other tool, for instance flyway.

## Example Usage

```hcl
provider redshift {
  url = "localhost",
  user = "testroot",
  password = "Rootpass123",
  database = "dev"
}
```

## Argument Reference

* `url` - (Required) AWS Redshift Endpoint
* `user` - (Required) redshift database user
* `password` - (Required) redshift database user
* `database` - (Required) redshift database to target
