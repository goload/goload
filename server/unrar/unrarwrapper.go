package unrar

import (
	"os/exec"
	"log"
	"os"
)

func Unrar(file string,dest string,password string) error{
	if password =="" {
		password = "Empty"
	}
	cmd := exec.Command("unrar", "x", "-o-","-p"+password, file, dest )
	cmd.Stdout = os.Stdout
    err := cmd.Start()
	if  err != nil {
		log.Fatal(err)
	}
	log.Printf("Extracting "+file+" to "+ dest)
	err = cmd.Wait()
	if(err == nil) {
		log.Printf("Done extracting "+file)
	} else {
		log.Printf("Error extracting "+file+" password wrong/missing or part missing?")
	}
	return err
}