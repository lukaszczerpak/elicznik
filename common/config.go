package common

import "time"

type AppConfig struct {
	General struct {
		Timezone          string         `mapstructure:"timezone" validate:"required"`
		Location          *time.Location `mapstructure:"-"`
		DeleteBeforeWrite bool           `mapstructure:"-"`
		DumpToFile        string         `mapstructure:"-"`
		Debug             bool           `mapstructure:"-"`
	} `mapstructure:"general"`
	Influxdb struct {
		Bucket      string `mapstructure:"bucket" validate:"required"`
		Org         string `mapstructure:"org" validate:"required"`
		Measurement string `mapstructure:"measurement" validate:"required"`
		Url         string `mapstructure:"url" validate:"required"`
		Token       string `mapstructure:"token" validate:"required"`
	} `mapstructure:"influxdb"`
	Tauron struct {
		Login         string  `mapstructure:"login" validate:"required"`
		Password      string  `mapstructure:"password" validate:"required"`
		SmartNr       string  `mapstructure:"smart-nr" validate:"required"`
		StartDate     string  `mapstructure:"start-date" validate:"required"`
		StorageFactor float64 `mapstructure:"storage-factor" validate:"required"`
	} `mapstructure:"tauron"`
}
