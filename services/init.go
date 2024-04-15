package services

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/email"
	"gitlab.tgl-cloud.com/dx-ecosystem/crm/pkg/jwt"
)

var (
	TokenMaker 	*jwt.JWTMaker
	Config     	*ServerConfig
	Mailer 		*email.Mailer
)

func InitServices() {
	viper.SetDefault("jwt_key.crm", "golanginuse-secret-key")
	viper.SetDefault("jwt_key.tms", "oq0gLsgK0jX9CDh5s0kfrYXvpRO0nHDG")
	TokenMaker = &jwt.JWTMaker{
		SecretKey: viper.GetString("jwt_key.crm"),
		TMSSecretKey: viper.GetString("jwt_key.tms"),
		Lifetime:  4320,
		Issuer:    "TGL Solutions",
	}

	viper.SetDefault("email.display_name", "skyACE株式会社")

	Config = &ServerConfig{
		Format: &FormatConfig{
			DateFormat: "2006-01-02",
			TimeFormat: "2006-01-02 15:04:05",
		},
		Token: &TokenConfig{
			SecretKey: "Z1iZkwEIT02A0H8O4TJn0YLwhUxHzG07",
		},
		Salt: &SaltConfig{
			Key: "xILncfMVgOS0V6nz",
		},
		Domain: &DomainConfig{
			Domain: viper.GetString("service.domain"),
		},
		FileStorage: &FileStorageConfig{
			Folder:         "./public/",
			TempFolder: 	"./temp/",
			ThumbnailWidth: 200,
		},
		Mail: &MailConfig{
			Host: viper.GetString("email.host"),
			Port: viper.GetInt("email.port"),
			DisplayName: viper.GetString("email.display_name"),
			MailAddress: viper.GetString("email.mail_address"),
			Password: viper.GetString("email.password"),
		},
	}

	if Config.Domain.Domain[len(Config.Domain.Domain)-1:] != "/" {
		Config.Domain.Domain = Config.Domain.Domain + "/"
	}

	if _, err := os.Stat(Config.FileStorage.Folder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(Config.FileStorage.Folder, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	if _, err := os.Stat(Config.FileStorage.TempFolder); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(Config.FileStorage.TempFolder, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	Mailer = email.NewMailer(Config.Mail.Host, Config.Mail.Port, Config.Mail.MailAddress, Config.Mail.Password, Config.Mail.DisplayName)

	initHashIDs()

	initUploaders()

	StartMailQueue()

	fmt.Printf("INFO: Config.Domain %s \n", Config.Domain)
}
