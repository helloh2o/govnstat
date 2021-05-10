package main

import (
	"flag"
	"fmt"
	"log"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	args = flag.String("p", "-m", "vnstat args")
	max  = flag.Float64("max", 999.00, "max gb traffic todo")
	loop = flag.Int64("loop", 5, "how many minutes to check loop")
	ver  = flag.Int64("ver", 2, "the version of vnstat")
)

const (
	GiB = "GiB"
	MiB = "MiB"
)

func main() {
	flag.Parse()
	// check on startup
	check()
	// check loop
	tk := time.NewTicker(time.Minute * time.Duration(*loop))
	for {
		<-tk.C
		check()
	}
}

func check() {
	if *ver == 1 {
		RunCommand("vnstat", "-u")
	}
	result := RunCommandWith("vnstat", *args)
	lines := strings.Split(result, "\n")
	switch *args {
	case "-m":
		// date
		fmtStr := "%d-%d"
		y := time.Now().Year()
		m := time.Now().Month()
		ys := strconv.Itoa(y) + "-"
		month := ""
		for _, line := range lines {
			vvs := strings.Split(line, "|")
			for i, v := range vvs {
				if i == 0 {
					// not month data break
					if !strings.Contains(v, ys) {
						break
					}
					if m < 10 {
						fmtStr = "%d-0%d"
					}
					month = fmt.Sprintf(fmtStr, y, m)
					if !strings.Contains(v, month) {
						log.Printf("Skip passed month data %s", v)
						break
					}
				}
				if i == 2 {
					v = strings.Trim(v, " ")
					total := strings.Split(v, " ")
					if len(total) == 2 {
						amount, _ := strconv.ParseFloat(total[0], 32)
						dw := total[1]
						switch dw {
						case GiB:
							if amount < *max {
								log.Printf("%s traffic ok, used:: %.2f, left:: %.2f", month, amount, *max-amount)
							} else {
								log.Printf("%s -_-, oops traffic used up, god!!! exec poweroff cmd now.", month)
								RunCommandWith("poweroff")
							}
						default:
							log.Printf("ignore the wd -> %s", dw)
							return
						}
					}
				}
			}
		}
	case "--json":
	default:
		log.Printf("unhandle result:: %s", result)
	}
}


func RunCommand(cmdName string, arg ...string) {
	cmd := exec.Command(cmdName, arg...)
	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
	if err != nil {
		panic(err)
		return
	}
	// run
	if err := cmd.Start(); err != nil {
		panic(err)
		return
	}
	// result
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Println(string(opBytes))
}

func RunCommandWith(cmdName string, arg ...string) string {
	cmd := exec.Command(cmdName, arg...)
	stdout, err := cmd.StdoutPipe()
	defer stdout.Close()
	if err != nil {
		panic(err)
		return err.Error()
	}
	// run
	if err := cmd.Start(); err != nil {
		panic(err)
		return err.Error()
	}
	// result
	opBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	return string(opBytes)
}
