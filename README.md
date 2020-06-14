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
We have 2 config files, which are put into "resource" directory
- the default config file: application.yaml
- for SIT environment: application-sit.yaml

#### application.yaml
```yaml
ldap:
  server: localhost:389
  binding_format: uid=%s,ou=users,dc=sample,dc=com,dc=vn

redis_url: redis://localhost:6379
```

In the SIT environment, we just override the configuration of ldap server and redis server:
#### application-sit.yaml
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
	// "resource" is the directory, which contains application.yaml and application-sit.yaml 
	config.LoadConfig("authentication", "resource", env, &conf, "application")
	log.Println("config ", conf)
}
```
