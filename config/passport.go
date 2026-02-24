package config

type Passport struct {
	EndPoint     string `mapstructure:"end-point" json:"end-point" yaml:"end-point"`
	Application  string `mapstructure:"application" json:"application" yaml:"application"`
	Certificate  string `mapstructure:"certificate" json:"certificate" yaml:"certificate"`
	ClientUser   string `mapstructure:"client-user" json:"client-user" yaml:"client-user"`
	ClientID     string `mapstructure:"client-id" json:"client-id" yaml:"client-id"`
	ClientSecret string `mapstructure:"client-secret" json:"client-secret" yaml:"client-secret"`
}
