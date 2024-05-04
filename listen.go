/*
Largely experimental. I'm not even sure if this is possible.
The elementary idea currently is having a never ending goroutine being fed strings through a channel which get executed.

In reality this probably needs something a bit more involved.

*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

type ExecChannel chan string

func main() {

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(filepath.Glob("*"))
	fmt.Print(">> ")
	text, _ := reader.ReadString('\n')
	fmt.Println(text)

}

// this is a subroutine
func listener(ec ExecChannel) {
	for {
		// exec here
		var x chan string
		select {

		case x <- ec:
			// exec(x)

		}
	}
}
