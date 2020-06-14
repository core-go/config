package config

import (
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
)

// Load function will read config from environment or config file.
func LoadConfig(parentPath string, directory string, env string, c interface{}, fileNames ...string) error {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.SetConfigType("yaml")

	for _, fileName := range fileNames {
		viper.SetConfigName(fileName)
	}

	viper.AddConfigPath("./" + directory + "/")
	if len(parentPath) > 0 {
		viper.AddConfigPath("./" + parentPath + "/" + directory + "/")
	}

	if er1 := viper.ReadInConfig(); er1 != nil {
		switch er1.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("config file not found")
		default:
			return er1
		}
	}
	if len(env) > 0 {
		env2 := strings.ToLower(env)
		for _, fileName2 := range fileNames {
			name0 := fileName2 + "." + env2
			viper.SetConfigName(name0)
			er2a := viper.MergeInConfig()
			if er2a != nil {
				return er2a
			}
			name1 := fileName2 + "-" + env2
			viper.SetConfigName(name1)
			er2b := viper.MergeInConfig()
			if er2b != nil {
				return er2b
			}
		}
	}
	er3 := bindEnvs(c)
	if er3 != nil {
		return er3
	}
	er4 := viper.Unmarshal(c)
	return er4
}

// bindEnvs function will bind ymal file to struc model
func bindEnvs(conf interface{}, parts ...string) error {
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
			return bindEnvs(v.Interface(), append(parts, tv)...)
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

func LoadMap(parentPath string, directory string, env string, fileNames ...string) (map[string]string, error) {
	innerMap := make(map[string]string)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_", ".", "_"))
	viper.SetConfigType("yaml")

	for _, fileName := range fileNames {
		viper.SetConfigName(fileName)
	}

	viper.AddConfigPath("./" + directory + "/")
	if len(parentPath) > 0 {
		viper.AddConfigPath("./" + parentPath + "/" + directory + "/")
	}

	if er1 := viper.ReadInConfig(); er1 != nil {
		switch er1.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("config file not found")
		default:
			return nil, er1
		}
	}
	if len(env) > 0 {
		env2 := strings.ToLower(env)
		for _, fileName2 := range fileNames {
			name0 := fileName2 + "." + env2
			viper.SetConfigName(name0)
			er2a := viper.MergeInConfig()
			if er2a != nil {
				return nil, er2a
			}
			name1 := fileName2 + "-" + env2
			viper.SetConfigName(name1)
			er2b := viper.MergeInConfig()
			if er2b != nil {
				return nil, er2b
			}
		}
	}
	er3 := viper.Unmarshal(&innerMap)
	return innerMap, er3
}

func LoadFile(parentPath string, directory string, env string, filename string) ([]byte, error) {
	if len(env) > 0 {
		indexDot := strings.LastIndex(filename, ".")
		if indexDot >= 0 {
			file := "./" + directory + "/" + filename[0:indexDot] + "-" + env + filename[indexDot:]
			if !fileExists(file) {
				file = "./" + parentPath + "/" + directory + "/" + filename[0:indexDot] + "-" + env + filename[indexDot:]
			}
			if fileExists(file) {
				return ioutil.ReadFile(file)
			}
		}
	}

	file := "./" + directory + "/" + filename
	if !fileExists(file) {
		file = "./" + parentPath + "/" + directory + "/" + filename
	}
	return ioutil.ReadFile(file)
}

func LoadCredentials(parentPath string, directory string, env string, filename string) ([]byte, error) {
	return LoadFile(parentPath, directory, env, filename)
}

func LoadText(parentPath string, directory string, env string, filename string) (string, error) {
	rs, err := LoadFile(parentPath, directory, env, filename)
	if err != nil {
		return "", err
	}
	return string(rs), nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
