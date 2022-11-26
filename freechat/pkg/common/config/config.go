package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	// Root folder of this project
	Root = filepath.Join(filepath.Dir(b), "../../..")
)

var Config config

type config struct {
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
		LogLevel       int      `yaml:"logLevel"`
		SlowThreshold  int      `yaml:"slowThreshold"`
	}
	Mongo struct {
		DBUri                string   `yaml:"dbUri"`
		DBAddress            []string `yaml:"dbAddress"`
		DBDirect             bool     `yaml:"dbDirect"`
		DBTimeout            int      `yaml:"dbTimeout"`
		DBDatabase           string   `yaml:"dbDatabase"`
		DBSource             string   `yaml:"dbSource"`
		DBUserName           string   `yaml:"dbUserName"`
		DBPassword           string   `yaml:"dbPassword"`
		DBMaxPoolSize        int      `yaml:"dbMaxPoolSize"`
		DBRetainChatRecords  int      `yaml:"dbRetainChatRecords"`
		ChatRecordsClearTime string   `yaml:"chatRecordsClearTime"`
	}
	Redis struct {
		DBAddress     []string `yaml:"dbAddress"`
		DBMaxIdle     int      `yaml:"dbMaxIdle"`
		DBMaxActive   int      `yaml:"dbMaxActive"`
		DBIdleTimeout int      `yaml:"dbIdleTimeout"`
		DBUserName    string   `yaml:"dbUserName"`
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
	cfgName := os.Getenv("FREECHAT_CONFIG_NAME")
	if len(cfgName) != 0 {
		bytes, err := ioutil.ReadFile(filepath.Join(cfgName, "config", "config.yaml"))
		if err != nil {
			bytes, err = ioutil.ReadFile(filepath.Join(Root, "config", "config.yaml"))
			if err != nil {
				panic(err.Error() + " config: " + filepath.Join(cfgName, "config", "config.yaml"))
			}
		} else {
			Root = cfgName
		}
		if err = yaml.Unmarshal(bytes, &Config); err != nil {
			panic(err.Error())
		}
	} else {
		bytes, err := ioutil.ReadFile(filepath.Join(Root, "config", "config.yaml"))
		if err != nil {
			panic(err.Error())
		}
		if err = yaml.Unmarshal(bytes, &Config); err != nil {
			panic(err.Error())
		}
	}
}
