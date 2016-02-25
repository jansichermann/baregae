package cloudsqlutil

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type SqlObject interface {
	InsertParameters() map[string]interface{}
	TableName() string
}

func Database(user, password, projectId, cloudsqlInstance, database string) (*sql.DB, error) {
	return sql.Open("mysql", user+":"+password+"@cloudsql("+projectId+":"+cloudsqlInstance+")/"+database)
}

func InsertObject(db *sql.DB, obj SqlObject) error {
	insertParameters := obj.InsertParameters()
	parameterNames := make([]string, 0, len(insertParameters))
	parameterValues := make([]interface{}, 0, len(insertParameters))
	placeholders := make([]string, 0, len(insertParameters))
	for k, v := range insertParameters {
		parameterNames = append(parameterNames, k)
		parameterValues = append(parameterValues, v)
		placeholders = append(placeholders, "?")
	}
	keyString := strings.Join(parameterNames, ", ")
	placeholderString := strings.Join(placeholders, ", ")
	queryString := "INSERT INTO " + obj.TableName() + " (" + keyString + ") VALUES (" + placeholderString + ")"
	_, err := db.Exec(queryString, parameterValues...)
	return err
}

type Migration struct {
	Migration string
	Rollback  string
}

func (m Migration) TableName() string {
	return "migrations"
}

func (m Migration) InsertParameters() map[string]interface{} {
	return map[string]interface{}{"migration": m.Migration, "rollback": m.Rollback}
}

func EnsureMigrationTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS migrations (id int auto_increment primary key, migration varchar(255), rollback varchar(255));")
	return err;
}

func EnsureMigrations(db *sql.DB, migrations []Migration) error {
	rows, err := db.Query("SELECT id FROM migrations ORDER BY id desc limit 1")
	if err != nil {
		return err
	}
	defer rows.Close()

	lastMigration := 0
	if rows.Next() {
		rows.Scan(&lastMigration)
	}

	for i := lastMigration; i < len(migrations); i++ {
		m := migrations[i]

		_, err := db.Exec(m.Migration)
		if err != nil {
			return err
		}

		err = InsertObject(db, m)
		if err != nil {
			return err
		}
	}
	return nil
}
