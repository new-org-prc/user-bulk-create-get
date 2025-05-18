package main

import (
	"flag"
	"log"
	"os"
	http_server "sika/api/http"
	"sika/config"
	"sika/pkg/load"
	"sika/service"
	"time"
)

var configPath = flag.String("config", "config.yaml", "configuration path")
var inputFilePath = flag.String("file", "data/users_data.json", "json file path")

const importFlagFile = ".data_imported"

func isDataImported() bool {
	_, err := os.Stat(importFlagFile)
	return err == nil
}

func markDataAsImported() error {
	file, err := os.Create(importFlagFile)
	if err != nil {
		return err
	}
	return file.Close()
}

func main() {
	flag.Parse()
	cfg := readConfig()
	app, err := service.NewAppContainer(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if !isDataImported() {

		if err := app.UserService().ClearUserAndAddressDataFromDB(); err != nil {
			log.Fatalf("Error clearing existing data: %v", err)
		}
		start := time.Now()
		userData, err := load.LoadData(*inputFilePath)
		if err != nil {
			log.Printf("Error loading data: %v", err)
		} else {
			err = app.UserService().ImportUsers(userData)
			if err != nil {
				log.Printf("Error importing users: %v", err)
			} else {
				if err := markDataAsImported(); err != nil {
					log.Printf("Warning: Could not mark data as imported: %v", err)
				}
			}
			log.Printf("it took %s to import data successfully", time.Since(start))
		}
	} else {
		log.Println("data already imported, if you need to import again, please delete the file .data_imported from the root directory and run the program again")
	}

	http_server.Run(cfg, app)
}

func readConfig() config.Config {

	if cfgPathEnv := os.Getenv("APP_CONFIG_PATH"); len(cfgPathEnv) > 0 {
		*configPath = cfgPathEnv
	}

	if len(*configPath) == 0 {
		log.Fatal("configuration file not found")
	}

	cfg, err := config.ReadStandard(*configPath)

	if err != nil {
		log.Fatal(err)
	}

	return cfg
}