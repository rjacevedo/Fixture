package fixture

import (
	"./fixture"
	"fmt"
	"time"
)


func main() {
	go fixture.RunAnalysis(10, 20, 5)

	for {
		fmt.Println("Main Program Run")
		time.Sleep(time.Second)
	}

}