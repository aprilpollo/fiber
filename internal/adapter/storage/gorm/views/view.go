package views

var Views = map[string]string{
	"vw_users": `
	CREATE OR REPLACE VIEW vw_users AS
	SELECT * FROM users
	`,
}
