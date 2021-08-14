package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func flagUsage() {
	usageText := `create and source virtualenv environment

Usage:
venv [arguments]
The commands are:
create     create virtualenv environment and setting direnv
del        delete virtualenv environment and setting direnv
path       show venv path`
	fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func venv_path() string {
	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	current_dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	hash := md5.Sum([]byte(current_dir))
	basename := filepath.Base(current_dir)
	hash_path := fmt.Sprintf("%x", hash)
	return homedir + "/" + ".venv/" + hash_path + "/" + basename
}

func create_venv() {
	venv_path := venv_path()

	_, err := os.Stat(venv_path)
	if !os.IsNotExist(err) {
		fmt.Println("already created virtualenv for this directry")
		fmt.Println("Please 'venv del'")
		return
	}
	_, err = exec.Command("python", "-m", "virtualenv", venv_path).Output()
	if err != nil {
		log.Fatal(err)
	}
	output, err := exec.Command("python", "--version").Output()
	if err != nil {
		log.Fatal(err)
	}
	version := strings.TrimSpace(string(output))
	fmt.Println("create virtualenv [", version, "]")
}

func create_env() {
	venv_path := venv_path()
	fp, err := os.Create(".envrc")
	if err != nil {
		log.Fatal(err)
	}
	fp.WriteString("source " + venv_path + "/bin/activate\n")
	fp.WriteString("unset PS1")
	defer fp.Close()
	output, err := exec.Command("direnv", "allow").Output()
	if err != nil {
		fmt.Println("koko")
		log.Fatal(output, err)
	}
}

func show_path() {
	venv_path := venv_path()
	_, err := os.Stat(venv_path)
	if os.IsNotExist(err) {
		fmt.Println("not created virtualenv for this directry")
		return
	}
	fmt.Println(venv_path)
}

func delete_venv() {
	venv_path := venv_path()
	dir_path := filepath.Dir(venv_path)
	_, err := os.Stat(dir_path)
	if !os.IsNotExist(err) {
		err = os.RemoveAll(dir_path)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err = os.Stat(".envrc")
	if !os.IsNotExist(err) {
		err = os.Remove(".envrc")
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("remove virtualenv for this directry")
}

func main() {
	flag.Usage = flagUsage
	_ = flag.NewFlagSet("create", flag.ExitOnError)
	_ = flag.NewFlagSet("path", flag.ExitOnError)
	_ = flag.NewFlagSet("del", flag.ExitOnError)
	if len(os.Args) == 1 {
		flag.Usage()
		return
	}
	switch os.Args[1] {
	case "create":
		create_venv()
		create_env()
	case "path":
		show_path()
	case "del":
		delete_venv()
	default:
		flag.Usage()
	}
	//flag.Parse()
	//args := flag.Args()
}
