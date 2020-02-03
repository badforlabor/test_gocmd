/**
 * Auth :   liubo
 * Date :   2020/1/31 20:55
 * Comment: 模拟控制台输入输出
 */

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Println("output 1")
	os.Stderr.WriteString("output err\n")
	//os.Exit(0)

	fmt.Print("enter a word...")
	var data, err =  bufio.NewReader(os.Stdin).ReadBytes('\n')
	if err == nil {
		fmt.Print(string(data))
	}

	fmt.Println("done.")
}
