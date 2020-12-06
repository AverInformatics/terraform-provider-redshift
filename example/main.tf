provider redshift {
  url = var.url
  user = var.username
  password = var.password
  database = var.database_primary
  ssl_mode = "disable"
}

resource redshift_user test_user {
  username = "testusernew"
  password = "Testpass123"
  connection_limit = "4"
  createdb = true
}

resource redshift_user test_user2 {
  username = "test_user8"
  password = "Testpass123"
  connection_limit = "1"
  createdb = true
}


resource redshift_group test_group {
  group_name = "test_group"
  users = [redshift_user.test_user.id, redshift_user.test_user2.id]
}

resource redshift_schema test_schema {
  schema_name = "test_schemax"
  cascade_on_delete = true
}

resource redshift_group_schema_privilege testgroup_testchema_privileges {
  schema_id = redshift_schema.test_schema.id
  group_id = redshift_group.test_group.id
  select = true
  insert = true
  update = true
}
