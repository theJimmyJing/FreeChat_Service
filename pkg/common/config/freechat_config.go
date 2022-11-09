package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
)

var FreechatConfig freechatConfig

type freechatConfig struct {
	ServerVersion string `yaml:"serverversion"`

	Api struct {
		GinPort  []int  `yaml:"openImApiPort"`
		ListenIP string `yaml:"listenIP"`
	}
	CmsApi struct {
		GinPort  []int  `yaml:"openImCmsApiPort"`
		ListenIP string `yaml:"listenIP"`
	}

	Mysql struct {
		DBAddress      []string `yaml:"dbMysqlAddress"`
		DBUserName     string   `yaml:"dbMysqlUserName"`
		DBPassword     string   `yaml:"dbMysqlPassword"`
		DBDatabaseName string   `yaml:"dbMysqlDatabaseName"`
		DBTableName    string   `yaml:"DBTableName"`
		DBMsgTableNum  int      `yaml:"dbMsgTableNum"`
		DBMaxOpenConns int      `yaml:"dbMaxOpenConns"`
		DBMaxIdleConns int      `yaml:"dbMaxIdleConns"`
		DBMaxLifeTime  int      `yaml:"dbMaxLifeTime"`
	}
	Mongo struct {
		DBUri               string `yaml:"dbUri"`
		DBAddress           string `yaml:"dbAddress"`
		DBDirect            bool   `yaml:"dbDirect"`
		DBTimeout           int    `yaml:"dbTimeout"`
		DBDatabase          string `yaml:"dbDatabase"`
		DBSource            string `yaml:"dbSource"`
		DBUserName          string `yaml:"dbUserName"`
		DBPassword          string `yaml:"dbPassword"`
		DBMaxPoolSize       int    `yaml:"dbMaxPoolSize"`
		DBRetainChatRecords int    `yaml:"dbRetainChatRecords"`
	}
	Redis struct {
		DBAddress     []string `yaml:"dbAddress"`
		DBMaxIdle     int      `yaml:"dbMaxIdle"`
		DBMaxActive   int      `yaml:"dbMaxActive"`
		DBIdleTimeout int      `yaml:"dbIdleTimeout"`
		DBPassWord    string   `yaml:"dbPassWord"`
		EnableCluster bool     `yaml:"enableCluster"`
	}

	Freechat struct {
		ListenIP string `yaml:"listenIP"`
		Port     []int  `yaml:"freechatPort"`

		TestDepartMentID string `yaml:"testDepartMentID"`
		ImAPIURL         string `yaml:"imAPIURL"`
	}
	Rtc struct {
		SignalTimeout string `yaml:"signalTimeout"`
	} `yaml:"rtc"`
}

func init() {
	cfgName := os.Getenv("CONFIG_NAME")
	fmt.Println(Root, cfgName)

	if len(cfgName) == 0 {
		cfgName = Root + "/config/freechat.yaml"
	}

	bytes, err := ioutil.ReadFile(cfgName)
	if err != nil {
		panic(err.Error())
	}
	if err = yaml.Unmarshal(bytes, &FreechatConfig); err != nil {
		panic(err.Error())
	}
}
