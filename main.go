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
	log.Info("•••••••••••••••••••••••")
	log.Info("MY APP Starting.... version: " + version)
	log.Info("••••••••••••••••••••••••")

	configuration = configs.NewConfiguration()
	configuration.Init()

	debug := configuration.GetBool("debug")
	if debug {
		log.Println("Running in debug mode")
	}

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
