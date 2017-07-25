package work

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func ReadFile(filename string, info *Info) error {
	bytes, err := ioutil.ReadFile(filename)

	if err != nil {
		fmt.Println("ReadFile: ", err.Error())
		return err
	}

	if err := json.Unmarshal(bytes, info); err != nil {
		fmt.Println("Unmarshal: ", err.Error())
		return err
	}

	return nil
}
