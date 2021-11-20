package main

type Config struct {
	Influxdb struct {
		Bucket string `mapstructure:"bucket" validate:"required"`
		Org    string `mapstructure:"org" validate:"required"`
		Url    string `mapstructure:"url" validate:"required"`
		Token  string `mapstructure:"token" validate:"required"`
	} `mapstructure:"influxdb"`
	Tauron struct {
		Login     string `mapstructure:"login" validate:"required"`
		Password  string `mapstructure:"password" validate:"required"`
		SmartNr   string `mapstructure:"smart-nr" validate:"required"`
		StartDate string `mapstructure:"start-date" validate:"required"`
	} `mapstructure:"tauron"`
}
