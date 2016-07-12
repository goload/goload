package unrar

import (
	"os/exec"
	"log"
	"errors"
	"regexp"
	"os"
	"strconv"
)

const PERCENT_REGEX = `(?P<PC>[1]{0,1}[0-9]{1,2})%`

type UnrarProgess struct {
	Progess float64
	Error error
}

type PercentWriter struct {
	c chan UnrarProgess
}

func (p PercentWriter) Write(b []byte) (n int, err error) {
	re := regexp.MustCompile(PERCENT_REGEX)
	n1 := re.SubexpNames()
	matches := re.FindAllStringSubmatch(string(b), -1)
	if (matches == nil) {
		return len(b),nil
	}
	r2 := matches[0]
	md := map[string]string{}
	for i, n := range r2 {
		md[n1[i]] = n
	}
	percent,_:=strconv.Atoi(md["PC"])
	p.c <- UnrarProgess{Progess:float64(percent)}
	return len(b), nil
}

func Unrar(filePath string, dest string, password string) (chan UnrarProgess) {
	c := make(chan UnrarProgess)
	go func() {
		if password == "" {
			password = "Empty"
		}
		cmd := exec.Command("unrar", "x", "-o-", "-p" + password, filePath, dest)
		cmd.Stdout = &PercentWriter{c:c}
		cmd.Stderr = os.Stderr
		if err := cmd.Start(); err != nil {
			c <- UnrarProgess{Progess:0.0,Error:errors.New("Unrar could not be executed")}
			close(c)
			return
		}
		log.Printf("Extracting " + filePath + " to " + dest)
		err := cmd.Wait()
		if (err == nil) {
			c <- UnrarProgess{Progess:100.0}
			close(c)
		} else {
			c <- UnrarProgess{Progess:0.0,Error:errors.New("Error extracting " + filePath + " password wrong/missing or part missing?")}
			close(c)
		}
	} ()
	return c
}