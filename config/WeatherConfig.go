package config

import (
	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

type WeatherConfig struct {
	FileName string
}

func NewConfigByName(fileName string) *WeatherConfig {
	c := WeatherConfig{FileName: fileName}
	if err := c.initConfig(); err != nil {
		log.Info(err)
		os.Exit(-1)
	}
	return &c
}
func NewConfig() *WeatherConfig {
	c := WeatherConfig{}
	if err := c.initConfig(); err != nil {
		log.Info(err)
		os.Exit(-1)
	}
	return &c
}

func (c *WeatherConfig) initConfig() error {
	if c.FileName != "" {
		// 如果指定了配置文件，则解析指定的配置文件
		viper.SetConfigFile(c.FileName)
	} else {
		// 如果没有指定配置文件，则解析默认的配置文件
		viper.AddConfigPath("WeatherForGo/")
		viper.SetConfigName("config")
	}
	// 设置配置文件格式为YAML
	viper.SetConfigType("yaml")
	// viper解析配置文件
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// 监听配置文件是否改变,用于热更新
func (c *WeatherConfig) WatchConfig(f func()) {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Info("WeatherConfig file changed: ", e.Name)
		f()
	})
}

func (c *WeatherConfig) GetValue(key string) string {
	return viper.GetString(key)
}

func (c *WeatherConfig) GetValues(key string) []string {
	return viper.GetStringSlice(key)
}

func (c *WeatherConfig) GetViperUnmarshal(rawVal interface{}) error {
	return viper.Unmarshal(rawVal)
}
