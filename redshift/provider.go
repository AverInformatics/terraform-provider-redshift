package redshift

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"log"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Description: "Redshift url",
				Required:    true,
			},
			"user": {
				Type:        schema.TypeString,
				Description: "master user",
				Required:    true,
			},
			"password": {
				Type:        schema.TypeString,
				Description: "master password",
				Required:    true,
			},
			"port": {
				Type:        schema.TypeString,
				Description: "port",
				Optional:    true,
				Default:     "5439",
			},
			"sslmode": {
				Type:        schema.TypeString,
				Description: "SSL mode (require, disable, verify-ca, verify-full)",
				Optional:    true,
				Default:     "require",
			},
			"database": {
				Type:        schema.TypeString,
				Description: "default database",
				Optional:    true,
				Default:     "dev",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"redshift_user":     redshiftUser(),
			"redshift_group":    redshiftGroup(),
			"redshift_database": redshiftDatabase(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {

	config := Config{
		url:      d.Get("url").(string),
		user:     d.Get("user").(string),
		password: d.Get("password").(string),
		port:     d.Get("port").(string),
		sslmode:  d.Get("sslmode").(string),
		database: d.Get("database").(string),
	}

	log.Println("[INFO] Initializing Redshift client")
	client := config.Client()
	db := client.db

	//Test the connection
	err := db.Ping()

	if err != nil {
		log.Println("[ERROR] Could not establish Redshift connection with credentials")
		panic(err.Error()) // proper error handling instead of panic
	}

	return client, nil
}
