package apollo

import (
	"errors"
	"fmt"
	"github.com/apolloconfig/agollo/v4"
	agolloconfig "github.com/apolloconfig/agollo/v4/env/config"
	"github.com/apolloconfig/agollo/v4/storage"
	"github.com/spf13/viper"
	"gitlab.chuangzhen-sh.net/golang/xgo/logging/applogger"
	"go.uber.org/zap"
	"log"
	"strings"
	"sync"
)

var (
	once     sync.Once
	instance *client

	// ErrUninitialized returned when client not initialized.
	ErrClientUninitialized = errors.New("apollo client not initialized")
)

type client struct {
	agollo.Client
	VMap *viper.Viper
}

type Changer struct {
	sync.RWMutex
}

func (c *Changer) OnChange(event *storage.ChangeEvent) {
	c.Lock()
	defer c.Unlock()
	for key, value := range event.Changes {
		if value.ChangeType == storage.ADDED {
			instance.VMap.Set(key, value.NewValue)
		} else if value.ChangeType == storage.MODIFIED {
			instance.VMap.Set(key, value.NewValue)
		} else if value.ChangeType == storage.DELETED {
			instance.VMap.Set(key, nil)
		}
		applogger.Info(fmt.Sprintf("key:%s---ChangeType:%v---OldValue:%v---NewValue:%v", key, value.ChangeType, value.OldValue, value.NewValue))
	}
}

func (c *Changer) OnNewestChange(event *storage.FullChangeEvent) {

}

func Init(cfg agolloconfig.AppConfig) {
	once.Do(func() {
		loggerInterface := zap.S()
		defer loggerInterface.Sync()
		agollo.SetLogger(loggerInterface)
		ac, err := agollo.StartWithConfig(func() (*agolloconfig.AppConfig, error) {
			return &cfg, nil
		})

		if err != nil {
			log.Fatal(err)
		}

		instance = &client{
			Client: ac,
			VMap:   viper.New(),
		}

		for _, n := range strings.Split(cfg.NamespaceName, ",") {
			cache := instance.GetConfigCache(n)

			cache.Range(func(key, value interface{}) bool {
				instance.VMap.Set(key.(string), value)
				log.Println(fmt.Sprintf("key:%v, value:%v, type:%T", key, value, value))
				return true
			})
		}

		instance.AddChangeListener(new(Changer))
	})
}

func RenderConfig(v interface{}) {
	if err := instance.VMap.Unmarshal(&v); err != nil {
		log.Fatal(err)
	}
}

func Client() *client {
	if instance == nil {
		log.Fatal(ErrClientUninitialized)
	}
	return instance
}
