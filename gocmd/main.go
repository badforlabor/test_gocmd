/**
 * Auth :   liubo
 * Date :   2020/1/31 20:55
 * Comment:
 */

package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)
func copyAndCapture(w io.Writer, r io.Reader, input io.Writer) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
			fmt.Println("read--->", string(d))
			fmt.Println("<---read")
			if strings.Contains( string(d), "enter a word...") {
				input.Write([]byte("123\n"))
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}


func main() {

	var cmd = NewCmder("simple.exe", "")
	cmd.Run(func(self* Cmder, words string) {
		if words == "enter a word..." {
			self.Write([]byte("123"))
		}
	})
	fmt.Println(string(cmd.outBuffer.Bytes()))
	fmt.Println(string(cmd.errBuffer.Bytes()))

	fmt.Println("gocmd done.")
}

func main1() {

	var cmd = exec.Command("simple.exe")
	cmd.Dir = "./"

	var pErr, _ =  cmd.StderrPipe()
	var pOut, _ = cmd.StdoutPipe()
	var pIn, _ = cmd.StdinPipe()

	var e = cmd.Start()
	if e != nil {
		fmt.Println(e.Error())
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		copyAndCapture(os.Stdout, pOut, pIn)
		wg.Done()
	}()
	go func() {
		copyAndCapture(os.Stderr, pErr, pIn)
		wg.Done()
	}()

	e = cmd.Wait()
	pIn.Write([]byte("1234444\n"))
	fmt.Println("cmd wait done")
	if e != nil {
		fmt.Println("cmd:", e.Error())
	}
	wg.Wait()

	fmt.Println("done.")
}