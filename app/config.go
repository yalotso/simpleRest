package app

import (
	"fmt"
	"github.com/spf13/viper"
)

var Config appConfig

type appConfig struct {
	ServerPort                string `mapstructure:"server_port"`
	ServerAddress             string `mapstructure:"server_address"`
	DSN                       string `mapstructure:"dsn"`
	SmtpHost                  string `mapstructure:"smtp_host"`
	SmtpPort                  int    `mapstructure:"smtp_port"`
	EmailUser                 string `mapstructure:"email_user"`
	EmailPassword             string `mapstructure:"email_password"`
	EmailVerificationTemplate string `mapstructure:"email_verification_template"`
	ResetPasswordTemplate     string `mapstructure:"reset_password_template"`
	JWTSigningMethod          string `mapstructure:"jwt_signing_method"`
	JWTSigningKey             string `mapstructure:"jwt_signing_key"`
	JWTVerificationKey        string `mapstructure:"jwt_verification_key"`
}

// LoadConfig loads configuration from the given list of paths and populates it into the Config variable.
// The configuration file(s) should be named as app.yaml.
// Environment variables with the prefix "RESTFUL_" in their names are also read automatically.
func LoadConfig(configPaths ...string) error {
	v := viper.New()
	v.SetConfigName("app")
	v.SetConfigType("yaml")
	v.SetEnvPrefix("restful")
	v.AutomaticEnv()
	for _, path := range configPaths {
		v.AddConfigPath(path)
	}
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read the configuration file: %s", err)
	}
	if err := v.Unmarshal(&Config); err != nil {
		return err
	}
	return nil
}
