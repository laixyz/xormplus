package xormhelper

import (
	"errors"
	_ "github.com/go-sql-driver/mysql"
	xorm "github.com/laixyz/xormplus"
	"github.com/laixyz/xormplus/log"
	"github.com/laixyz/xormplus/names"
	"sync"
	"time"
)

const (
	SnakeMapper = "snake"
	SameMapper  = "same"
	GonicMapper = "gonic"
)

type XormSession struct {
	session         *xorm.Session
	Name            string
	DSN             string
	Prefix          string
	Mapper          string
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
	Debug           bool
}

var ConnMaxLifetimeDefault time.Duration = 30 * time.Second
var MaxIdleConnsDefault int = 16
var MaxOpenConnsDefault int = 32

func (xc *XormSession) Connect() (err error) {
	engine, err := xorm.NewEngine("mysql", xc.DSN)
	if err != nil {
		return err
	}
	var tbMapper names.PrefixMapper
	switch xc.Mapper {
	case GonicMapper:
		tbMapper = names.NewPrefixMapper(names.GonicMapper{}, xc.Prefix)
	case SameMapper:
		tbMapper = names.NewPrefixMapper(names.SameMapper{}, xc.Prefix)
	case SnakeMapper:
		tbMapper = names.NewPrefixMapper(names.SnakeMapper{}, xc.Prefix)
	default:
		tbMapper = names.NewPrefixMapper(names.SnakeMapper{}, xc.Prefix)
	}
	engine.SetTableMapper(tbMapper)
	if xc.ConnMaxLifetime > 0 {
		engine.SetConnMaxLifetime(xc.ConnMaxLifetime)
	} else {
		engine.SetConnMaxLifetime(ConnMaxLifetimeDefault)
	}
	if xc.MaxIdleConns > 0 {
		engine.SetMaxIdleConns(xc.MaxIdleConns)
	} else {
		engine.SetMaxIdleConns(MaxIdleConnsDefault)
	}
	if xc.MaxOpenConns > 0 {
		engine.SetMaxOpenConns(xc.MaxOpenConns)
	} else {
		engine.SetMaxOpenConns(MaxOpenConnsDefault)
	}

	if xc.Debug {
		engine.ShowSQL(true)
		//xorm 日志设置，只显示错误日志
		engine.Logger().SetLevel(log.LOG_DEBUG)
	} else {
		engine.Logger().SetLevel(log.LOG_ERR)
	}
	xc.session = engine.NewSession()
	return nil
}
func (xc *XormSession) Close() error {
	return xc.session.Engine().Close()
}

type Xorms struct {
	Engine map[string]*XormSession
	sync.RWMutex
}

var PublicXorms Xorms = Xorms{Engine: make(map[string]*XormSession)}

// Register
func (xc *XormSession) Register() error {
	err := xc.Connect()
	if err != nil {
		return err
	}
	PublicXorms.Lock()
	defer PublicXorms.Unlock()
	if _, ok := PublicXorms.Engine[xc.Name]; ok {
		delete(PublicXorms.Engine, xc.Name)
	}
	PublicXorms.Engine[xc.Name] = xc
	return nil
}

/*
Register 注册一个mysql连接池
范例:
	给连接池取个名称
	Register("default","test:test@(127.0.0.1:3306)/db?charset=utf8mb4&parseTime=True&loc=Local")
	指定表前缀
	Register("default","test:test@(127.0.0.1:3306)/db?charset=utf8mb4&parseTime=True&loc=Local","prefix_")
	设置调试模式，会打印sql语句
	Register("default","test:test@(127.0.0.1:3306)/db?charset=utf8mb4&parseTime=True&loc=Local","prefix_","true")
	指定xorm的mapper规则
	Register("default","test:test@(127.0.0.1:3306)/db?charset=utf8mb4&parseTime=True&loc=Local","prefix_","true","gonic)
*/
func Register(poolName, dsn string, params ...string) error {
	var newXormConfig XormSession

	newXormConfig.Name = poolName
	newXormConfig.DSN = dsn
	l := len(params)
	if l > 0 && params[0] != "" {
		newXormConfig.Prefix = params[0]
	}
	if l > 1 && params[1] == "true" {
		newXormConfig.Debug = true
	}
	if l > 2 && params[2] != "" {
		newXormConfig.Mapper = params[2]
	} else {
		newXormConfig.Mapper = SnakeMapper
	}
	err := newXormConfig.Connect()
	if err != nil {
		return err
	}
	PublicXorms.Lock()
	defer PublicXorms.Unlock()
	if _, ok := PublicXorms.Engine[newXormConfig.Name]; ok {
		delete(PublicXorms.Engine, newXormConfig.Name)
	}
	PublicXorms.Engine[newXormConfig.Name] = &newXormConfig
	return nil
}

// Using 使用指定的名字的连接池， 不指定默认为default
func Using(params ...string) (*xorm.Session, error) {
	PublicXorms.Lock()
	defer PublicXorms.Unlock()
	var name string
	if len(params) == 0 {
		name = "default"
	} else {
		name = params[0]
	}
	if _, ok := PublicXorms.Engine[name]; ok {
		if err := PublicXorms.Engine[name].session.Ping(); err != nil {
			err = PublicXorms.Engine[name].Connect()
			if err != nil {
				return nil, err
			}
			return PublicXorms.Engine[name].session, err
		}
		return PublicXorms.Engine[name].session, nil
	}
	return nil, errors.New("the mysql connect not register")
}

/* Close 关闭连接
 范例：
	Close() 关闭名字为default的连接
	Close("my") 关闭名字为my的连接
	Close("my","your") 关闭多个连接
	Close("__all__") 关闭所有的连接。
*/
func Close(params ...string) (err error) {
	PublicXorms.Lock()
	defer PublicXorms.Unlock()
	var names []string
	if len(params) == 0 {
		names = []string{"default"}
	} else if len(params) == 1 && params[0] == "__all__" {
		for _, xc := range PublicXorms.Engine {
			err = xc.Close()
			if err != nil {
				return
			}
		}
		return nil
	} else {
		names = params
	}
	for _, name := range names {
		if _, ok := PublicXorms.Engine[name]; ok {
			err = PublicXorms.Engine[name].Close()
			if err != nil {
				return
			}
		}
	}
	return nil
}
