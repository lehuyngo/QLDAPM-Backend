package services

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/util"
)

type SaltConfig struct {
	Key string
}

type FormatConfig struct {
	DateFormat string
	TimeFormat string
}

type FileStorageConfig struct {
	Folder         string
	ThumbnailWidth int
	TempFolder     string
}

type TokenConfig struct {
	SecretKey string
}

type DomainConfig struct {
	Domain string
}

type MailConfig struct {
	Host string
	Port int
	DisplayName string
	MailAddress string
	Password string
}

type ServerConfig struct {
	Format  *FormatConfig
	Token 	*TokenConfig
	Salt	*SaltConfig
	Domain	*DomainConfig
	FileStorage *FileStorageConfig
	Mail *MailConfig
}

func InitViper() {

	cfgFile := "./config.toml"
	fmt.Printf("Read config file from env: [%s] \n", cfgFile)

	folder, fileName, ext, err := util.ExtractFilePath(cfgFile)
	if err != nil {
		fmt.Printf("Extract config file failed %s err: %s \n", viper.ConfigFileUsed(), err.Error())
		os.Exit(-1)
	}
	fmt.Printf("Extract config file success folder[%s] fileName[%s] ext[%s] \n", folder, fileName, ext)

	// Setting
	viper.AddConfigPath(folder)
	viper.SetConfigName(fileName)
	viper.AutomaticEnv()
	viper.SetConfigType(ext)

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("FATAL: Viper using config file failed %s err: %s \n", viper.ConfigFileUsed(), err.Error())
		os.Exit(-1)
	}

	fmt.Printf("Service using config file: %s \n", viper.ConfigFileUsed())
	//watch on config change
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s", e.Name)
	})

	fmt.Println("Start initialize config success.")
}
