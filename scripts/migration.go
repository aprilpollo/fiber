package main

import (
	config "aprilpollo/internal/adapter/config"

	"fmt"
	"log"
	"strings"
	gormOrm "aprilpollo/internal/adapter/storage/gorm"
	"aprilpollo/internal/adapter/storage/gorm/models"
	"aprilpollo/internal/adapter/storage/gorm/views"

	"github.com/fatih/color"
)


var (
	green = color.New(color.FgGreen).SprintFunc()
	blue  = color.New(color.FgBlue).SprintFunc()
	cyan  = color.New(color.FgCyan).SprintFunc()
	red   = color.New(color.FgRed).SprintFunc()
)

func printSuccess(msg string) {
	fmt.Printf("%s %s\n", green("[v]"), msg)
}

func printInfo(msg string) {
	fmt.Printf("%s %s\n", blue("[INFO]"), msg)
}

func printHeader(msg string) {
	fmt.Printf("\n%s %s\n", cyan("[INFO]"), msg)
	fmt.Printf("%s\n", cyan(strings.Repeat("=", 50)))
}

func main() {
	printHeader("Starting Database Migration")

	printInfo("Loading configuration...")
	config, err := config.GetConfig()
	if err != nil {
		log.Fatalf("%s failed to load configuration: %v", red("[x]"), err)
	}

	printInfo("Connecting to database...")
	db, err := gormOrm.NewGormDB(config.Database, nil)
	if err != nil {
		log.Fatalf("%s failed to connect DB: %v", red("[x]"), err)
	}
	defer db.Close()
	printSuccess("Database connected successfully")

	printInfo("Running auto migration...")
	if err := db.Migrate(models.All()...); err != nil {
		log.Fatalf("%s failed to migrate: %v", red("[x]"), err)
	}

	printInfo("Dropping existing views...")
	for name := range views.Views {
		dropSQL := fmt.Sprintf("DROP VIEW IF EXISTS %s CASCADE;", name)
		if err := db.GetDB().Exec(dropSQL).Error; err != nil {
			log.Fatalf("%s failed to drop view %s: %v", red("[x]"), name, err)
		}
		printSuccess(fmt.Sprintf("Dropped view: %s", name))
	}

	printInfo("Creating new views...")
	for name, query := range views.Views {
		if err := db.GetDB().Exec(query).Error; err != nil {
			log.Fatalf("%s failed to create view %s: %v", red("[x]"), name, err)
		}
		printSuccess(fmt.Sprintf("Created view: %s", name))
	}

	printHeader("Migration Completed Successfully!")
}