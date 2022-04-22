package executor

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/fynxiu/dbd/internal/constant"

	"github.com/go-sql-driver/mysql"
	"github.com/golang/glog"
)

const (
	databaseName       = constant.DatabaseName
	statementDelimiter = ";\n"
)

func init() {
	registerExecutor(constant.EngineMysql, newMysqlExecutor)
}

var _ Executor = (*mysqlExecutor)(nil)

func newMysqlExecutor(DataSourceName string) Executor {
	mysql.SetLogger(&mysqlBlackHoleLogger{})
	return &mysqlExecutor{
		dsn: DataSourceName,
	}
}

type mysqlExecutor struct {
	dsn  string
	db   *sql.DB
	once sync.Once
}

// ExecuteScript implements Executor
func (e *mysqlExecutor) ExecuteScript(script string) error {
	glog.V(3).Infof("ExecuteScript: %s", script)
	var err error
	e.once.Do(func() {
		if err = e.ensureCleanEnviroment(); err != nil {
			glog.Errorf("ExecuteScript failed, %v", err)
			err = fmt.Errorf("ExecuteScript failed, %v", err)
		}
	})
	if err != nil {
		return err
	}

	db, err := e.getConn()
	if err != nil {
		return err
	}

	for _, x := range strings.Split(script, statementDelimiter) {
		if x = strings.TrimSpace(x); x == "" {
			continue
		}
		if _, err := db.Exec(x); err != nil {
			glog.Infof("ExecuteScript failed, %v\nscript:%v", err, script)
			return err
		}
	}

	return nil
}

// Schema implements Executor
func (e *mysqlExecutor) Schema() (string, error) {
	db, err := e.getConn()
	if err != nil {
		return "", err
	}

	tables, err := getTables(db)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, x := range tables {
		var tn, ts string
		if err := db.QueryRow(fmt.Sprintf("SHOW CREATE TABLE %s", x)).Scan(&tn, &ts); err != nil {
			return "", err
		}
		sb.WriteString("\n")
		sb.WriteString(ts)
		sb.WriteString(";\n")
	}

	return sb.String(), nil
}

// Dispose implements Executor
func (e *mysqlExecutor) Dispose() error {
	if e.db != nil {
		return e.db.Close()
	}
	return nil
}

func (e *mysqlExecutor) getConn() (*sql.DB, error) {
	var err error
	if e.db != nil {
		return e.db, nil
	}
	e.db, err = sql.Open("mysql", e.dsn)
	if err != nil {
		return nil, err
	}
	e.db.SetConnMaxLifetime(time.Minute)
	e.db.SetMaxOpenConns(10)
	e.db.SetMaxIdleConns(10)

	timeout := time.After(time.Second * 30)
	ticker := time.NewTicker(time.Millisecond * 500)
	defer ticker.Stop()
	tick := ticker.C
	for {
		select {
		case <-timeout:
			return nil, fmt.Errorf("connect mysql timeout")
		case <-tick:
			if err := e.db.Ping(); err != nil {
				glog.V(2).Infof("ping mysql failed, %v", err)
				continue
			}
			return e.db, nil
		}
	}
}

func getTables(db *sql.DB) ([]string, error) {
	var tables []string
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tables, nil
}

func (e *mysqlExecutor) ensureCleanEnviroment() error {
	db, err := e.getConn()
	if err != nil {
		return fmt.Errorf("ensureCleanEnviroment failed, %v", err)
	}

	// tables, err := getTables(db)
	// if err != nil {
	// 	return fmt.Errorf("ensureCleanEnviroment failed, %v", err)
	// }

	// for _, x := range tables {
	// 	if _, err := db.Exec(fmt.Sprintf("DROP TABLE %s", x)); err != nil {
	// 		return fmt.Errorf("ensureCleanEnviroment failed, %v", err)
	// 	}
	// }

	if _, err := db.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s;", databaseName)); err != nil {
		glog.Error(err)
		return err
	}
	if _, err := db.Exec(fmt.Sprintf("CREATE DATABASE %s;", databaseName)); err != nil {
		glog.Error(err)
		return err
	}
	if _, err := db.Exec(fmt.Sprintf("USE %s;", databaseName)); err != nil {
		glog.Error(err)
		return err
	}

	return nil
}
