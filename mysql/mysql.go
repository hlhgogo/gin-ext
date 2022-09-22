package mysql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/hlhgogo/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var ClientMap = make(map[string]*gorm.DB)
var clientMapLock sync.Mutex

const (
	DefaultCli = "default"
)

func Load() error {
	configMap, err := config.GetMysqlOptions()
	if err != nil {
		return err
	}

	for connName, v := range configMap {
		dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", v.Username, v.Password, v.Addr, v.Name)
		conf := mysql.Config{
			DSN:                       dsn,
			DefaultStringSize:         0,     // string 类型字段的默认长度
			DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
			DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
			DontSupportRenameColumn:   true,  //  用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
			SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
		}
		if err := initClient(connName, conf, v.MaxConnNum, v.MaxIdleConn); err != nil {
			return err
		}
	}
	return nil
}

func initClient(client string, conf mysql.Config, maxConnNum int, maxIdleConn int) error {
	clientMapLock.Lock()
	defer clientMapLock.Unlock()

	Conn, err := gorm.Open(mysql.New(conf), &gorm.Config{SkipDefaultTransaction: false})
	if err != nil {
		return err
	}

	sqlDB, err := Conn.DB()
	if err != nil {
		return err
	}

	// 设置数据库连接池参数
	sqlDB.SetMaxOpenConns(maxConnNum)  // 设置数据库连接池最大连接数
	sqlDB.SetMaxIdleConns(maxIdleConn) // 最多空闲数量
	ClientMap[client] = Conn

	return nil
}

// DefaultClient 获取default客户端
func DefaultClient() *gorm.DB {
	return ClientMap[DefaultCli]
}

func Client(name string) *gorm.DB {
	return ClientMap[name]
}
