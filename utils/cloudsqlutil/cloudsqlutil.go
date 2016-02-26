package cloudsqlutil

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type AppDatabase struct {
  db *sql.DB
}

func (db AppDatabase)InsertObject(obj SqlObject) error {
	return insertObject(db.db, obj)
}

func (db AppDatabase)EnsureMigrationTable() error {
	return ensureMigrationTable(db.db);
}

func (db AppDatabase)EnsureMigrations(migrations []Migration) error {
	return ensureMigrations(db.db, migrations);
}

type SqlObject interface {
	InsertParameters() map[string]interface{}
	TableName() string
}

func Database(user, password, projectId, cloudsqlInstance, database string) (AppDatabase, error) {
	db, err := sql.Open("mysql", user+":"+password+"@cloudsql("+projectId+":"+cloudsqlInstance+")/"+database)
	adb := AppDatabase{db: db}
	return adb, err;
}

func insertObject(db *sql.DB, obj SqlObject) error {
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

func ensureMigrationTable(db *sql.DB) error {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS migrations (id int auto_increment primary key, migration varchar(255), rollback varchar(255));")
	return err;
}

func ensureMigrations(db *sql.DB, migrations []Migration) error {
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

		err = insertObject(db, m)
		if err != nil {
			return err
		}
	}
	return nil
}
