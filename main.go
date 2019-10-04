/**
MY OWN VERSION OF MAIN...........
*/
package main

import (
	"log"
	"time"
)

....
"time"
)

const version = "0.0.1 Gasoline"

var configuration *configs.ViperConfiguration

func init() {

	configuration = configs.NewConfiguration()
	configuration.Init()

	debug := configuration.GetBool("debug")
	log.Init(debug)

	log.Info("----------------------------------------------")
	log.Info("XXXXXXx Starting.... version: " + version)
	log.Info("----------------------------------------------")

}

func main() {

	if err := database.InitMySQLDatabase(configuration.Get("mysqlURI")); err != nil {
		log.Fatal(err)
	}

	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for {
			service.UpdateAccountBalance(configuration)
			<-ticker.C
		}
	}()

	select {}
}
