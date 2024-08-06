package mysql

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
)

const (
	DbOpSuccess = "success"
	DbOpError   = "error"
)

// InstrumentedMysql is an instance of mysql client instrumented using prometheus
type InstrumentedMysql struct {
	*gorm.DB
	MetricCount    *prometheus.CounterVec
	MetricDuration *prometheus.HistogramVec
	Log            *zap.SugaredLogger
}

// NewInstrumentedMysql returns a new instance of an instrumented mysql
func NewInstrumentedMysql(dsn, namespace, subsystem string, logger *zap.SugaredLogger, opts ...func(opt *options)) (*InstrumentedMysql, error) {
	if logger == nil {
		return nil, errors.New("invalid logger")
	}

	opt := defaultOpts()

	for _, option := range opts {
		option(opt)
	}

	conf := &gorm.Config{
		Logger: gLogger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), gLogger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  opt.logLevel,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}),
		SkipDefaultTransaction: opt.skipDefaultTransaction,
		CreateBatchSize:        opt.createBatchSize,
	}

	db, err := gorm.Open(mysql.Open(dsn), conf)
	if err != nil {
		return nil, err
	}

	d, err := db.DB()
	if err != nil {
		return nil, err
	}

	d.SetMaxOpenConns(opt.maxOpenConnections)
	d.SetMaxIdleConns(opt.maxIdleConnections)
	d.SetConnMaxLifetime(opt.connectionMaxLifeTime)
	d.SetConnMaxIdleTime(opt.connectionMaxIdleTime)

	return &InstrumentedMysql{
		DB: db,
		MetricCount: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "db_ops_total",
				Help:      "The total number of database requests per function",
			},
			[]string{"op", "status"},
		),
		MetricDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "db_ops_duration_seconds",
				Help:      "The amount of time database request functions take",
			},
			[]string{"op", "status"},
		),
		Log: logger,
	}, nil
}
