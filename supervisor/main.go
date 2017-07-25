// supervisor project main.go
package main

import (
	//	"fmt"

	"supervisor/work"
	//	"supervisor/conf"
)

func main() {

	work.Start()

	//		xxxMap, err := conf.ReadFile("conf.json")
	//		if err != nil {
	//			fmt.Println("readFile: ", err.Error())
	//			return
	//		}
	//			fmt.Println(xxxMap)
}
