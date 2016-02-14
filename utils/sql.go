package utils

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

type SqlObject interface {
	InsertParameters() map[string]interface{}
	TableName() string
}

func CloudSqlDatabase(user, password, projectId, cloudsqlInstance, database string) (*sql.DB, error) {
	return sql.Open("mysql", user+":"+password+"@cloudsql("+projectId+":"+cloudsqlInstance+")/"+database)
}

func InsertObject(db *sql.DB, obj SqlObject) error {
	insertParameters := obj.InsertParameters()
	parameterNames := make([]string, len(insertParameters))
	parameterValues := make([]interface{}, len(insertParameters))
	placeholders := make([]string, len(insertParameters))
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
