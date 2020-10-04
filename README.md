# Config Loader for Golang 

- Support to load the configurations
- Support to merge the environment variables to the configurations
- Support to override the configurations based on the environment (For example: SIT environment, UAT environment)

## Installation

Please make sure to initialize a Go module before installing common-go/config:

```shell
go get -u github.com/common-go/config
```

Import:

```go
import "github.com/common-go/config"
```

## Example
We have 2 config files, which are put into "configs" directory
- the default config file: config.yaml
- for SIT environment: config-sit.yaml

#### config.yaml
```yaml
ldap:
  server: localhost:389
  binding_format: uid=%s,ou=users,dc=sample,dc=com,dc=vn

redis_url: redis://localhost:6379
```

In the SIT environment, we just override the configuration of ldap server and redis server:
#### config-sit.yaml
```yaml
ldap:
  server: sit-server:389

redis_url: redis://redis-sit:6379
```

```go
type LDAPConfig struct {
	Server        string `mapstructure:"server"`
	BindingFormat string `mapstructure:"binding_format"`
}

type RootConfig struct {
	Server      app.ServerConfig `mapstructure:"server"`
	Ldap        LDAPConfig       `mapstructure:"ldap"`
	RedisClient string           `mapstructure:"redis_client"`
}

func main() {
	env := os.Getenv("ENV")
	var conf RootConfig
	// "authentication" is the directory, which contains source code
	// "configs" is the directory, which contains config.yaml and config-sit.yaml 
	config.LoadConfigWithEnv("authentication", "configs", env, &conf, "config")
	log.Println("config ", conf)

	var conf2 RootConfig
	config.Load(&conf2, "configs/config")
	log.Println("config2 ", conf2)
}
```
