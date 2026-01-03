package databases

import (
	"fmt"
	"log"

	"github.com/miceremwirigi/journey-family-sermons/m/cmd/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func StartDatabase(env string) (db *gorm.DB) {
	log.Println("Starting Database ...")
	conf := config.LoadConfig()
	if conf.DatabaseUrl != "" {
		db, err := gorm.Open(postgres.Open(conf.DatabaseUrl), &gorm.Config{})
		if err == nil {
		log.Println("Successfully Started Database")
		return db
		}
	}

	var (
		db_host, db_user, db_pass, db_name, db_ssl, db_port string
		err                                                 error
	)
	if env == "test" {
		_, db_host, db_user, db_pass, db_name, db_ssl, db_port, err = LoadTestDatabaseConfig()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		_, db_host, db_user, db_pass, db_name, db_ssl, db_port, err = LoadDatabaseConfig()
		if err != nil {
			log.Fatal(err)
		}
	}

	dsn := fmt.Sprintf(`host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Africa/Nairobi`,
		db_host, db_user, db_pass, db_name, db_port, db_ssl)
	log.Println("dsn: "+ dsn)
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err == nil {
		log.Println("Successfully Started Database")
	} else {
		log.Fatal(err)
	}
	return db
}
