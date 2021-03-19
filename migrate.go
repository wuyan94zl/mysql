package mysql

func AutoMigrate(lists map[string]interface{}){
	for _, v := range lists {
		_ = DB.AutoMigrate(v)
	}
}
