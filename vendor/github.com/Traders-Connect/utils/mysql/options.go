package mysql

import (
	"time"

	gLogger "gorm.io/gorm/logger"
)

type options struct {
	logLevel               gLogger.LogLevel
	maxOpenConnections     int
	maxIdleConnections     int
	connectionMaxLifeTime  time.Duration
	connectionMaxIdleTime  time.Duration
	skipDefaultTransaction bool
	createBatchSize        int
	//Replicas []string
}

func defaultOpts() *options {
	return &options{
		logLevel:               gLogger.Error,
		maxOpenConnections:     10,
		maxIdleConnections:     5,
		connectionMaxLifeTime:  time.Hour,
		connectionMaxIdleTime:  100 * time.Second,
		skipDefaultTransaction: false,
		createBatchSize:        1000,
	}
}

// SetMaxOpenConnections sets the number of maximum open connections. default=10
func SetMaxOpenConnections(num int) func(opt *options) {
	return func(opt *options) {
		opt.maxOpenConnections = num
	}
}

// SetMaxIdleConnections sets the number of maximum idle connections. default=5
func SetMaxIdleConnections(num int) func(opt *options) {
	return func(opt *options) {
		opt.maxIdleConnections = num
	}
}

// SetConnectionMaxLifeTime sets connection max lifetime. default=1 hour
func SetConnectionMaxLifeTime(t time.Duration) func(opt *options) {
	return func(opt *options) {
		opt.connectionMaxLifeTime = t
	}
}

// SetConnectionMaxIdleTime sets connection max idle time. default=100 seconds
func SetConnectionMaxIdleTime(t time.Duration) func(opt *options) {
	return func(opt *options) {
		opt.connectionMaxIdleTime = t
	}
}

// SetLogLevel sets log level for gorm logger. default=error
func SetLogLevel(level gLogger.LogLevel) func(opt *options) {
	return func(opt *options) {
		opt.logLevel = level
	}
}

// SetSkipDefaultTransaction sets the value for skip gorm's default transaction. default=false
func SetSkipDefaultTransaction(skip bool) func(opt *options) {
	return func(opt *options) {
		opt.skipDefaultTransaction = skip
	}
}

// SetCreateBatchSize sets default create batch size
func SetCreateBatchSize(size int) func(opt *options) {
	return func(opt *options) {
		opt.createBatchSize = size
	}
}
