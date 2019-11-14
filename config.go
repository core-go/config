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
func LoadConfig(parentPath string, directory string, env string, c interface{}, fileNames ...string) interface{} {
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

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("config file not found")
		default:
			panic(err)
		}
	}
	if len(env) > 0 {
		env2 := strings.ToLower(env)
		for _, fileName2 := range fileNames {
			name := fileName2 + "-" + env2
			viper.SetConfigName(name)
			viper.MergeInConfig()
		}
	}
	bindEnvs(c)
	viper.Unmarshal(c)
	return c
}

// bindEnvs function will bind ymal file to struc model
func bindEnvs(conf interface{}, parts ...string) {
	ifv := reflect.Indirect(reflect.ValueOf(conf))
	ift := reflect.TypeOf(ifv)
	for i := 0; i < ift.NumField(); i++ {
		v := ifv.Field(i)
		t := ift.Field(i)
		tv, ok := t.Tag.Lookup("mapstructure")
		if !ok {
			continue
		}
		switch v.Kind() {
		case reflect.Struct:
			bindEnvs(v.Interface(), append(parts, tv)...)
		default:
			viper.BindEnv(strings.Join(append(parts, tv), "."))
		}
	}
}

func LoadMap(parentPath string, directory string, env string, fileNames ...string) map[string]string {
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

	if err := viper.ReadInConfig(); err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError:
			log.Println("config file not found")
		default:
			panic(err)
		}
	}
	if len(env) > 0 {
		env2 := strings.ToLower(env)
		for _, fileName2 := range fileNames {
			name := fileName2 + "-" + env2
			viper.SetConfigName(name)
			viper.MergeInConfig()
		}
	}
	viper.Unmarshal(&innerMap)
	return innerMap
}

func LoadCredentials(parentPath string, directory string, filename string) ([]byte, error) {
	file := "./" + directory + "/" + filename
	if !fileExists(file) {
		file = "./" + parentPath + "/" + directory + "/" + filename
	}
	return ioutil.ReadFile(file)
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
