/**
 * Auth :   liubo
 * Date :   2020/1/31 22:32
 * Comment: Cmder
 */

package main

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type Cmder struct {
	cmd *exec.Cmd

	wg sync.WaitGroup
	pErr io.ReadCloser
	pOut io.ReadCloser
	pIn io.WriteCloser

	errBuffer bytes.Buffer
	errNextLine int
	outBuffer bytes.Buffer
	outNextLine int
}
func (self *Cmder) Write(data []byte) {
	self.pIn.Write(data)
}
func (self *Cmder) Run(pCallback func(words string)) {

	self.wg.Add(2)

	var err = self.cmd.Start()
	if err != nil {
		fmt.Println(err.Error())
	}

	self.errBuffer.Reset()
	self.outBuffer.Reset()

	go func() {
		self.copyAndCapture(self.pErr, &self.errBuffer, &self.errNextLine, pCallback)
		self.wg.Done()
	}()
	go func() {
		self.copyAndCapture(self.pOut, &self.outBuffer, &self.outNextLine, pCallback)
		self.wg.Done()
	}()

	err = self.cmd.Wait()
	if err != nil {
		fmt.Println(err.Error())
	}

	self.wg.Wait()
}

func (self *Cmder) copyAndCapture(r io.ReadCloser, input *bytes.Buffer, nextLine *int, pCallback func(words string)) error {

	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			input.Write(d)

			if pCallback != nil {
				pCallback(string(d))
			}
			//var data = input.Bytes()
			//for pCallback != nil {
			//	var idx = bytes.IndexByte( data[*nextLine:], '\n')
			//	if idx == -1 {
			//		break
			//	}
			//
			//	var line = data[*nextLine:idx]
			//	pCallback(string(line))
			//
			//	*nextLine = idx + 1
			//}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return err
		}
	}
}

func NewCmder(exe string, workDir string, arg ...string) *Cmder {
	var ret = &Cmder{}
	var cmd = exec.Command(exe, arg...)
	ret.cmd	= cmd
	if len(workDir) > 0 {
		ret.cmd.Dir = workDir
	}
	ret.pErr, _ =  cmd.StderrPipe()
	ret.pOut, _ = cmd.StdoutPipe()
	ret.pIn, _ = cmd.StdinPipe()

	return ret
}