package config

import "os"

func Set() {
	os.Setenv("SERVER", "dbs1.dc.pizzasoft.ru")
	//os.Setenv("SERVER", "localhost")
	os.Setenv("PORT", "9000")
	os.Setenv("METHOD", "POST")
	os.Setenv("DISTRICTS", "5")
	return
}
