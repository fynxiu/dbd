package executor

type mysqlBlackHoleLogger struct{}

func (l *mysqlBlackHoleLogger) Print(v ...interface{}) {}
