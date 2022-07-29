package gdb

import (
	"fmt"
	"github.com/Ravior/gserver/core/os/glog"
	"github.com/Ravior/gserver/core/util/gconfig"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"time"
)

const (
	// Default conn name for instance usage.
	defaultDbConn      = "default"
	defaultMaxIdleConn = 10  // 最大空闲连接数，默认10
	defaultMaxOpenConn = 100 // 最大打开连接数，默认100
)

var (
	// Instances map containing db instances.
	dbInstances = make(map[string]*gorm.DB, 0)
)

func InitMysql() {
	for pool, config := range gconfig.Global.DataBase {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", config.User, config.Password, config.Host, config.Port, config.DbName)
		logLevel := logger.Error
		if config.ShowLog {
			logLevel = logger.Info
		}
		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.New(
				log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
				logger.Config{
					SlowThreshold:             time.Second, // Slow SQL threshold
					LogLevel:                  logLevel,    // Log level
					IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound errors for logger
					Colorful:                  true,        // Disable color
				},
			),
			NamingStrategy: schema.NamingStrategy{
				TablePrefix:   config.Prefix,
				SingularTable: true, // use singular table name, table for `User` would be `user` with this option enabled
			},
		})

		if err != nil {
			glog.Errorf("链接数据库出现错误, 链接信息：%s, 错误信息：%v", dsn, err)
			return
		}

		// 链接池设置
		if sqlDB, err := db.DB(); err == nil {
			// 设置链接池最大空闲连接数
			maxIdleConn := defaultMaxIdleConn
			if config.MaxIdleConn > 0 {
				maxIdleConn = config.MaxIdleConn
			}
			sqlDB.SetMaxIdleConns(maxIdleConn)
			// 设置链接池最大打开连接数
			maxOpenConn := defaultMaxOpenConn
			if config.MaxOpenConn > 0 {
				maxOpenConn = config.MaxOpenConn
			}
			sqlDB.SetMaxOpenConns(maxOpenConn)
			sqlDB.SetConnMaxLifetime(time.Hour)
		}

		dbInstances[pool] = db
	}
}

func GetDb(pools ...string) *gorm.DB {
	pool := defaultDbConn
	if len(pools) > 0 && pools[0] != "" {
		pool = pools[0]
	}

	if db, ok := dbInstances[pool]; ok {
		return db
	}
	return nil
}
