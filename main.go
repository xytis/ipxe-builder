package main

import (
	"github.com/gin-gonic/gin"

	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func cmdExec(command string, args []string, timeout int32) (rv string, rc bool) {
	cmd := exec.Cmd{Path: command, Args: args, Dir: "/ipxe"}

	var stdoutput bytes.Buffer
	var stderror bytes.Buffer
	cmd.Env = os.Environ()
	cmd.Stdout = &stdoutput
	cmd.Stderr = &stderror
	if err := cmd.Run(); err != nil {
		outputstr := stdoutput.String()
		errstr := stderror.String()
		fmt.Println("Failed to run", outputstr, "error:", errstr, err)
		return
	}

	rv = stdoutput.String()
	rc = true
	return
}

type buildPostData struct {
	Config string `json:"config"`
	Target string `json:"target"`
}

func abort(c *gin.Context, code int) {
	c.Status(code)
	c.Abort()
}

func main() {
	r := gin.Default()
	r.POST("/build", func(c *gin.Context) {
		data := &buildPostData{}
		err := c.Bind(data)
		if err != nil {
			abort(c, 500)
			return
		}
		if data.Config == "" {
			abort(c, 400)
			return
		}
		if data.Target == "" {
			data.Target = "bin/ipxe.usb"
		}
		if dest, err := ioutil.TempFile("/tmp", "embed-"); err == nil {
			defer os.Remove(dest.Name())
			dest.WriteString(data.Config)
			binary, _ := exec.LookPath("make")
			if _, rc := cmdExec(binary, []string{"make", data.Target, fmt.Sprintf("EMBED=%s", dest.Name())}, 600); !rc {
				abort(c, 500)
				return
			}
			c.Status(200)
			c.File(data.Target)
		}
	})

	r.Run(":8080")
}
