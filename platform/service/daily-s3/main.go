package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("begin schedulehandler")
	defer fmt.Println("end schedulehandler")

	startTs := time.Now().Unix()

	for i := 0; i < 20; i++ {
		fmt.Println("Value of i:", i)
	}

	endTs := time.Now().Unix()
	fmt.Println("Total time: ", (endTs - startTs))

}
