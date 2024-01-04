package itypes

import "fmt"

type DatabaseInfo struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
	SSLMode  string
	TimeZone string
}

func (dbInfo *DatabaseInfo) ToString() string {
	return fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=%v TimeZone=%v", dbInfo.Host, dbInfo.User, dbInfo.Password, dbInfo.DBName, dbInfo.Port, dbInfo.SSLMode, dbInfo.TimeZone)
}
