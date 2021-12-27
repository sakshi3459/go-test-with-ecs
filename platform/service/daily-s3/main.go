package main

import (
	"fmt"
	"time"
	//"errors"
	
	//"./lib"
	//"github.com/sakshi3459/go-test-with-ecs/platform/service/daily-s3/lib"
)

func main() {
	fmt.Println("begin schedulehandler")
	defer fmt.Println("end schedulehandler")

	startTs := time.Now().Unix()

	for i := 0; i < 20; i++ {
		fmt.Println("Value of i:", i)
	}
	
	//lib.Test()
	
	//err := errors.New("error in sample appln")
	//fmt.Println("Failure: ", err)
	//return

	endTs := time.Now().Unix()
	fmt.Println("Total time: ", (endTs - startTs))

}
