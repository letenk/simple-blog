package main

import (
	"github.com/jabutech/simple-blog/helper"
	"github.com/jabutech/simple-blog/router"
	"github.com/jabutech/simple-blog/util"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load config
	config, err := util.LoadConfig(".") // "." as location file app.env in root folder
	helper.FatalError("error load config: ", err)

	// Open connection to db
	db, err := gorm.Open(mysql.Open(config.DBSource), &gorm.Config{})
	helper.FatalError("Connection to database failed!, err: ", err)

	// Router
	server := router.NewRouter(db)

	// Run router
	server.Run()
}
