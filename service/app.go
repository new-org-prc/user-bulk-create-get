package service

import (
	"log"
	"sika/config"
	"sika/internal/address"
	"sika/internal/user"
	"sika/pkg/storage"

	"gorm.io/gorm"
)

type AppContainer struct {
	cfg    config.Config
	dbConn *gorm.DB
	userService    *UserService
}

func NewAppContainer(cfg config.Config) (*AppContainer, error) {
	app := &AppContainer{
		cfg: cfg,
	}

	app.mustInitDB()

	app.SetUserService()
	return app, nil
}

func (a *AppContainer) mustInitDB() {
	if a.dbConn != nil {
		return
	}

	db, err := storage.NewPostgresGormConnection(a.cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	a.dbConn = db

	err = storage.Migrate(a.dbConn)
	if err != nil {
		log.Fatal("Migration failed: ", err)
	}
}

func(a *AppContainer)SetUserService(){
	if a.userService != nil{
		return
	}
	a.userService = NewUserService(user.NewOps(storage.NewUserRepo(a.dbConn)), address.NewOps(storage.NewAddressRepo(a.dbConn)))
}

func(a *AppContainer)UserService()*UserService{
	return a.userService
}
