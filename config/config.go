package config

import (
	"flag"
	"io/ioutil"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Instance *Config

type Config struct {
	Env             string `yaml:"Env"`        // 环境：prod、dev
	BaseUrl         string `yaml:"BaseUrl"`    // base url
	Port            string `yaml:"Port"`       // 端口
	LogFile         string `yaml:"LogFile"`    // 日志文件
	ShowSql         bool   `yaml:"ShowSql"`    // 是否显示日志
	StaticPath      string `yaml:"StaticPath"` // 静态文件目录
	PuzzleUrl       string `yaml:"PuzzleUrl"`
	NovelPath       string `yaml:"NovelPath"`
	NovelPathTxt    string `yaml:"NovelPathTxt"`
	NovelPathOutput string `yaml:"NovelPathOutput"`

	// 数据库配置
	DB struct {
		Url          string `yaml:"Url"`
		MaxIdleConns int    `yaml:"MaxIdleConns"`
		MaxOpenConns int    `yaml:"MaxOpenConns"`
	} `yaml:"DB"`

	Cron struct {
		Image   string `yaml:"Image"`
		Gold    string `yaml:"Gold"`
		Gold2   string `yaml:"Gold2"`
		Article string `yaml:"Article"`
		Weather string `yaml:"Weather"`
		Poetry  string `yaml:"Poetry"`
		Clock   string `yaml:"Clock"`
	} `yaml:"Cron"`

	Redis struct {
		Addr string `yaml:"Addr"`
		Pwd  string `yaml:"Pwd"`
	} `yaml:"Redis"`

	DAYU struct {
		APP_KEY           string `yaml:"APP_KEY"`
		APP_SECRET        string `yaml:"APP_SECRET"`
		SMS_TEMPLATE_CODE string `yaml:"SMS_TEMPLATE_CODE"`
	} `yaml:"DAYU"`

	// Github
	Github struct {
		ClientID     string `yaml:"ClientID"`
		ClientSecret string `yaml:"ClientSecret"`
	} `yaml:"Github"`

	// OSChina
	OSChina struct {
		ClientID     string `yaml:"ClientID"`
		ClientSecret string `yaml:"ClientSecret"`
	} `yaml:"OSChina"`

	// QQ登录
	QQConnect struct {
		AppId  string `yaml:"AppId"`
		AppKey string `yaml:"AppKey"`
	} `yaml:"QQConnect"`

	// 阿里云oss配置
	Uploader struct {
		Enable    string `yaml:"Enable"`
		AliyunOss struct {
			Host            string `yaml:"Host"`
			Bucket          string `yaml:"Bucket"`
			RealBucket      string `yaml:"RealBucket"`
			Endpoint        string `yaml:"Endpoint"`
			RegionId        string `yaml:"RegionId"`
			AccessId        string `yaml:"AccessId"`
			AccessSecret    string `yaml:"AccessSecret"`
			RoleArn         string `yaml:"RoleArn"`
			RoleSessionName string `yaml:"RoleSessionName"`
			RemotePath      string `yaml:"RemotePath"`
			RealPath        string `yaml:"RealPath"`
			StyleSplitter   string `yaml:"StyleSplitter"`
			StyleAvatar     string `yaml:"StyleAvatar"`
			StylePreview    string `yaml:"StylePreview"`
			StyleSmall      string `yaml:"StyleSmall"`
			StyleDetail     string `yaml:"StyleDetail"`
		} `yaml:"AliyunOss"`
		Local struct {
			Host     string `yaml:"Host"`
			Path     string `yaml:"Path"`
			BookPath string `yaml:"BookPath"`
			LogoPath string `yaml:"LogoPath"`
		} `yaml:"Local"`
	} `yaml:"Uploader"`

	// 百度ai
	BaiduAi struct {
		ApiKey    string `yaml:"ApiKey"`
		SecretKey string `yaml:"SecretKey"`
	} `yaml:"BaiduAi"`

	// 百度SEO相关配置
	// 文档：https://ziyuan.baidu.com/college/courseinfo?id=267&page=2#h2_article_title14
	BaiduSEO struct {
		Site  string `yaml:"Site"`
		Token string `yaml:"Token"`
	} `yaml:"BaiduSEO"`

	// 神马搜索SEO相关
	// 文档：https://zhanzhang.sm.cn/open/mip
	SmSEO struct {
		Site     string `yaml:"Site"`
		UserName string `yaml:"UserName"`
		Token    string `yaml:"Token"`
	} `yaml:"SmSEO"`

	// smtp
	Smtp struct {
		Host     string `yaml:"Host"`
		Port     string `yaml:"Port"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
		SSL      bool   `yaml:"SSL"`
	} `yaml:"Smtp"`

	// es
	Es struct {
		Url   string `yaml:"Url"`
		Index string `yaml:"Index"`
	} `yaml:"Es"`
}

func init() {
	var filename = flag.String("config", "./music.yaml", "配置文件路径")
	Instance = &Config{}
	if yamlFile, err := ioutil.ReadFile(*filename); err != nil {
		logrus.Error(err)
	} else if err = yaml.Unmarshal(yamlFile, Instance); err != nil {
		logrus.Error(err)
	}
}
