package clients

import (
	"fmt"

	"gitlab.tgl-cloud.com/dx-ecosystem/crm/entities"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	// "github.com/spf13/viper"
)

var MySQLClient *gorm.DB

func NewMySQLClient() (*gorm.DB, error) {
	// username := os.Getenv("CRM_DATABASE_USERNAME")
	// password := os.Getenv("CRM_DATABASE_PASSWORD")
	// addr := os.Getenv("CRM_DATABASE_ADDRESS")
	// dbname := os.Getenv("CRM_DATABASE_DBNAME")
	// if username == "" {
	// 	username = viper.GetString("database.username")
	// }
	// if password == "" {
	// 	password = viper.GetString("database.password")
	// }
	// if addr == "" {
	// 	addr = viper.GetString("database.address")
	// }
	// if dbname == "" {
	// 	dbname = viper.GetString("database.dbname")
	// }
	addr:= "127.0.0.1:33061"
	dbname:= "crmdb"
	username := "root"
	password:= "123456"

	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&collation=utf8mb4_unicode_ci&parseTime=True",
		username, password, addr, dbname)

	var err error
	client, err := gorm.Open(mysql.Open(connStr), &gorm.Config{
		PrepareStmt:                              false,
		Logger:                                   logger.Default.LogMode(logger.Info),
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}

	client = client.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4 auto_increment=1")

	// set data models many to many relationship
	/*
		err = client.SetupJoinTable(&Product{}, "Factories", &ProductFactory{})
		if err != nil {
			return nil, err
		}
	*/

	return client, nil
}

func AutoMigrate() error {

	/*
		err := MySQLClient.AutoMigrate(&entities.MigrateUser{})
		if err != nil {
			return err
		}
	*/

	err := MySQLClient.AutoMigrate(&entities.User{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.Organization{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.MeetingHighlight{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.Meeting{}, &entities.MeetingNote{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.Attendee{}, &entities.Contributor{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.MeetingNoteEditor{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.File{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.Client{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.ClientActivity{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.ClientProjectActivity{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.ClientNote{}, &entities.ClientTag{}, &entities.ClientTagActivity{}, &entities.ClientNoteActivity{}, &entities.ClientClientTag{}, &entities.ClientAttachFile{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.Contact{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.ContactActivity{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.ContactNote{}, &entities.ContactTag{}, &entities.ContactTagActivity{}, &entities.ContactNoteActivity{}, &entities.ContactMailActivity{}, &entities.ContactContactTag{}, &entities.ContactAttachFile{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.ContactMailShortClick{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.ContactClientActivity{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.ClientContact{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.Project{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.ContactProject{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.ContactProjectActivity{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.DraftContact{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.Task{}, &entities.TaskAssignee{}, &entities.TaskAttachFile{}, &entities.TaskComment{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.Mail{}, &entities.MailAttachFile{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.BatchMail{}, &entities.BatchMailReceiver{}, &entities.BatchMailCarbonCopy{}, &entities.BatchMailAttachFile{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.TrackedURL{}, &entities.EventClickURL{})
	if err != nil {
		return err
	}

	err = MySQLClient.AutoMigrate(&entities.StaticFile{})
	if err != nil {
		return err
	}

	return nil
}
