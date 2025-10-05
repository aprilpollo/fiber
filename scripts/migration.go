package main

import (
	"log"
	"os"
	"regexp"

	config "aprilpollo/internal/adapter/config"
	gormOrm "aprilpollo/internal/adapter/storage/gorm"
	"aprilpollo/internal/adapter/storage/gorm/models"
	"aprilpollo/internal/adapter/storage/gorm/views"

	"strconv"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// parseErrorMessage extracts the essential error information
func parseErrorMessage(err error) string {
	errorStr := err.Error()

	// Pattern to extract SQLSTATE code from any PostgreSQL error
	re := regexp.MustCompile(`\(SQLSTATE ([A-Z0-9]+)\)`)
	matches := re.FindStringSubmatch(errorStr)

	if len(matches) == 2 {
		return "SQLSTATE " + matches[1]
	}

	// If no SQLSTATE found, return the original error but truncated
	if len(errorStr) > 100 {
		return errorStr[:100] + "..."
	}

	return errorStr
}

func main() {
	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}
	db, err := gormOrm.NewGormDB(config.Database, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // Disable logging for migration
	})
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mTable := table.NewWriter()
	mTable.SetOutputMirror(os.Stdout)
	mTable.SetStyle(table.StyleRounded)
	//mTable.SetTitle("MIGRATION DATABASE")
	mTable.Style().Title.Align = text.AlignCenter
	mTable.Style().Options.DoNotColorBordersAndSeparators = true
	mTable.Style().Options.DrawBorder = false
	mTable.Style().Options.SeparateColumns = true
	mTable.Style().Options.SeparateFooter = true
	mTable.Style().Options.SeparateHeader = true
	mTable.Style().Options.SeparateRows = false

	// Set column widths and alignment

	mTable.AppendHeader(table.Row{"TABLES & VIEWS", "STATUS", "MESSAGE"})
	mTable.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 20, AlignHeader: text.AlignLeft},
		{Number: 2, WidthMin: 20, AlignHeader: text.AlignLeft},
		{Number: 3, WidthMin: 20, AlignHeader: text.AlignLeft},
	})

	successCount := 0
	failCount := 0

	for _, model := range models.All() {
		if err := db.Migrate(model); err != nil {
			mTable.AppendRow(table.Row{
				model.TableName(),
				text.Colors{text.FgRed}.Sprint("✗ Failed"),
				parseErrorMessage(err),
			})
			failCount++
		} else {
			mTable.AppendRow(table.Row{
				model.TableName(),
				text.Colors{text.FgGreen}.Sprint("✓ Migrated"),
				"SUCCESS",
			})
			successCount++
		}
	}

	mTable.AppendRow(table.Row{"-", "-", "-"})

	for viewName, viewSQL := range views.Views {
		if err := db.GetDB().Exec(viewSQL).Error; err != nil {
			mTable.AppendRow(table.Row{viewName, text.Colors{text.FgRed}.Sprint("✗ Failed"), parseErrorMessage(err)})
			failCount++
		} else {
			mTable.AppendRow(table.Row{viewName, text.Colors{text.FgGreen}.Sprint("✓ Migrated"), "SUCCESS"})
			successCount++
		}
	}

	mTable.AppendFooter(table.Row{"Summary", "", "Success: " + strconv.Itoa(successCount) + " Failed: " + strconv.Itoa(failCount)})
	mTable.Render()
}
