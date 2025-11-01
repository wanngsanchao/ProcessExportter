package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

//get the pid of the service_test process
func GetserviceTestPids(name string) ([]int,error) {
    if name == "" {
        return nil,errors.New("the process name is nil")
    }

    cmd := exec.Command("ps","aux")
    output,err := cmd.Output()

    if err != nil {
        return nil,errors.New("the command output error")
    }

    alllines := strings.Split(string(output),"\n")
    var pids []int

    for _,line := range alllines {
        if strings.Contains(line,name) && ! strings.Contains(line,"grep") {
            part := strings.Fields(line)

            if len(part) >= 2 {
                pid,err := strconv.Atoi(part[1])

                if err == nil {
                    pids = append(pids,pid)
                }
            }
        }
    }

    if len(pids) == 0 {
        return nil,errors.New(fmt.Sprintf("the %s process is not exits",name))
    }

    return pids,nil
}
