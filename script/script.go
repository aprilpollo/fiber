package main

import (
	"aprilpollo/internal/adapters/config"
	"aprilpollo/internal/adapters/storage/orm"
	"aprilpollo/internal/adapters/storage/orm/models"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"strconv"

	"regexp"
	"os"
	"log"
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

// isDuplicateError checks if the error is a duplicate key violation (SQLSTATE 23505)
func isDuplicateError(err error) bool {
	return parseErrorMessage(err) == "SQLSTATE 23505"
}

// appendResultRow appends a row to the table based on the operation result
func appendResultRow(mTable table.Writer, name string, err error, failCount *int, successCount *int) {
	if err != nil {
		msg := parseErrorMessage(err)
		var status string
		if isDuplicateError(err) {
			status = text.Colors{text.FgYellow}.Sprint("✓ Skipped")
			*successCount++
		} else {
			status = text.Colors{text.FgRed}.Sprint("✗ Failed")
			*failCount++
		}
		mTable.AppendRow(table.Row{name, status, msg})
	} else {
		mTable.AppendRow(table.Row{name, text.Colors{text.FgGreen}.Sprint("✓ Created"), "SUCCESS"})
		*successCount++
	}
}

func main(){
	cfg, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := orm.NewGormDB(cfg.Database, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}, true)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	mTable := table.NewWriter()
	mTable.SetOutputMirror(os.Stdout)
	mTable.SetStyle(table.StyleRounded)
	mTable.Style().Title.Align = text.AlignCenter
	mTable.Style().Options.DoNotColorBordersAndSeparators = true
	mTable.Style().Options.DrawBorder = false
	mTable.Style().Options.SeparateColumns = true
	mTable.Style().Options.SeparateFooter = true
	mTable.Style().Options.SeparateHeader = true
	mTable.Style().Options.SeparateRows = false

	mTable.AppendHeader(table.Row{"NAME", "STATUS", "MESSAGE"})
	mTable.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, WidthMin: 20, AlignHeader: text.AlignLeft},
		{Number: 2, WidthMin: 20, AlignHeader: text.AlignLeft},
		{Number: 3, WidthMin: 20, AlignHeader: text.AlignLeft},
	})

	successCount := 0
	failCount := 0

	// ── Migration ────────────────────────────────────────────────────────────
	mTable.AppendRow(table.Row{"TABLES", "", "MIGRATION TABLES"})
	mTable.AppendRow(table.Row{"-", "-", "-"})

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

	// ── Seed: Users ──────────────────────────────────────────────────────────
	mTable.AppendRow(table.Row{"", "", ""})
	mTable.AppendRow(table.Row{"USERS", "", "SEED DATA"})
	mTable.AppendRow(table.Row{"-", "-", "-"})

	seedUsers := []models.UserModel{
		{Name: "Alice Johnson", Email: "alice@example.com"},
		{Name: "Bob Smith",    Email: "bob@example.com"},
		{Name: "Carol White",  Email: "carol@example.com"},
	}

	for _, u := range seedUsers {
		err := db.GetDB().Create(&u).Error
		appendResultRow(mTable, u.Email, err, &failCount, &successCount)
	}

	// ── Seed: Books ──────────────────────────────────────────────────────────
	mTable.AppendRow(table.Row{"", "", ""})
	mTable.AppendRow(table.Row{"BOOKS", "", "SEED DATA"})
	mTable.AppendRow(table.Row{"-", "-", "-"})

	seedBooks := []models.BookModel{
		{Title: "The Go Programming Language", Author: "Alan Donovan"},
		{Title: "Clean Architecture",          Author: "Robert C. Martin"},
		{Title: "Domain-Driven Design",        Author: "Eric Evans"},
	}

	for _, b := range seedBooks {
		err := db.GetDB().Create(&b).Error
		appendResultRow(mTable, b.Title, err, &failCount, &successCount)
	}

	mTable.AppendFooter(table.Row{"Summary", "", "Success: " + strconv.Itoa(successCount) + " Failed: " + strconv.Itoa(failCount)})
	mTable.Render()
}
