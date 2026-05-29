package bootstrap

import "github.com/muhammadfarrasfajri/koperasi-gerai-be/config"

func InitDatabase() {
	config.ConnectMySQL()
}
