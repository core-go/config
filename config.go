package config

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strconv"
	"strings"
)

func Load(c interface{}, fileNames ...string) error {
	return LoadConfig("", "", c, fileNames...)
}

func LoadConfig(parentPath string, directory string, c interface{}, fileNames ...string) error {
	env := os.Getenv("ENV")
	return LoadConfigWithEnv(parentPath, directory, env, c, fileNames...)
}
// Load function will read config from environment or config file.
func LoadConfigWithEnv(parentPath string, directory string, env string, c interface{}, fileNames ...string) error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.SetConfigType("yaml")

	fileCount := len(fileNames)
	if fileCount > 0  {
		if len(parentPath) == 0 && len(directory) == 0 {
			viper.AddConfigPath("./")
		} else {
			viper.AddConfigPath("./" + directory + "/")
			if len(parentPath) > 0 {
				viper.AddConfigPath("./" + parentPath + "/" + directory + "/")
			}
		}

		viper.SetConfigName(fileNames[0])
		if er1 := viper.ReadInConfig(); er1 != nil {
			switch er1.(type) {
			case viper.ConfigFileNotFoundError:
				log.Println("config file not found")
			default:
				return er1
			}
		}

		for i:=1; i< fileCount; i++ {
			viper.SetConfigName(fileNames[i])
			if er2b := viper.MergeInConfig(); er2b != nil {
				switch er2b.(type) {
				case viper.ConfigFileNotFoundError:
					break
				default:
					return er2b
				}
			}
		}
	} else {
		return fmt.Errorf("have no config file")
	}

	if len(env) > 0 {
		env2 := strings.ToLower(env)
		for _, fileName2 := range fileNames {
			name0 := fileName2 + "." + env2
			viper.SetConfigName(name0)
			er2a := viper.MergeInConfig()
			if er2a != nil {
				switch er2a.(type) {
				case viper.ConfigFileNotFoundError:
					break
				default:
					return er2a
				}
			}
			name1 := fileName2 + "-" + env2
			viper.SetConfigName(name1)
			er2b := viper.MergeInConfig()
			if er2b != nil {
				switch er2b.(type) {
				case viper.ConfigFileNotFoundError:
					break
				default:
					return er2b
				}
			}
		}
	}
	er3 := BindEnvs(c)
	if er3 != nil {
		return er3
	}
	er4 := viper.Unmarshal(c)
	return er4
}

// bindEnvs function will bind ymal file to struc model
func BindEnvs(conf interface{}, parts ...string) error {
	ifv := reflect.Indirect(reflect.ValueOf(conf))
	ift := reflect.TypeOf(ifv)
	num := min(ift.NumField(), ifv.NumField())
	for i := 0; i < num; i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			return BindEnvs(v.Interface(), append(parts, tv)...)
		default:
			return viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
	return nil
}

func min(n1 int, n2 int) int {
	if n1 < n2 {
		return n1
	}
	return n2
}
