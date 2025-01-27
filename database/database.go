package database

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type FtpSource struct {
	User     string `json:"user"`
	Host     string `json:"host"`
	Password string `json:"password"`
}

type SourceConfig struct {
	ID     string          `gorm:"primaryKey"`
	Type   string          `json:"type"`
	Config json.RawMessage `json:"config"`
}

type SourceType int

const (
	SourceType_Ftp SourceType = iota + 1
	SourceType_S3
	SourceType_GCS
)

func Id() string {
	bytes := make([]byte, 24)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return hex.EncodeToString(bytes)
}

func getConfigData() SourceConfig {
	// generate random id of 24 characters for ID

	return SourceConfig{
		ID:   Id(),
		Type: "ftp",
		Config: json.RawMessage(`{
			"user": "nitesh",
			"host": "localhost",
			"password": "password"
		}`),
	}
}

func autoMigrate(db *gorm.DB) error {
	err := db.AutoMigrate(&SourceConfig{})
	return err
}

func fetchConfigData(db *gorm.DB) SourceConfig {
	source := SourceConfig{}
	if err := db.Model(&SourceConfig{}).First(&source).Error; err != nil {
		panic(err)
	}

	return source
}

func Database() (*gorm.DB, error) {

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", "localhost", "postgres", "password", "postgres", "5432")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Main() {
	db, err := Database()
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}

	err = autoMigrate(db)
	if err != nil {
		panic(err)
	}

	config := fetchConfigData(db)

	fmt.Println("Config: ", config.ID, config.Type, string(config.Config))

	//Unmarshall this config.Config to FtpSource struct
	ftpSource := FtpSource{}
	err = json.Unmarshal(config.Config, &ftpSource)
	if err != nil {
		panic(err)
	}

	fmt.Println("FTP Source: ", ftpSource.User, ftpSource.Host, ftpSource.Password)

	fmt.Print("DATA INSERTED SUCCESSFULLY")

	fmt.Println("ENUM VALUES : ", SourceType_Ftp)
	fmt.Println("ENUM VALUES : ", SourceType_S3)
	fmt.Println("ENUM VALUES : ", SourceType_GCS)

}

func CreateConfig(db *gorm.DB, config SourceConfig) error {
	if err := db.Model(&SourceConfig{}).Create(&config).Error; err != nil {
		return err
	}
	return nil
}
