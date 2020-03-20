package service

import "fmt"

func ExampleRequest() {
	result, _ := new(Result).Request("0600556", "20200324", "20200324")
	fmt.Println(result)
	// Output:
	// []
}
