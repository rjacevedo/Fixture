package fixture

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"bytes"
	"time"
	"runtime"
)

func checkRoutineAlert(routine_thresh int) bool {
	num_routines := runtime.NumGoroutine()
	
	fmt.Println("Number of go routines: ", num_routines)
	if num_routines >= routine_thresh {
		return true
	}

	return false
}

func checkFileAlert(fd_thresh int) bool {
	var fd_count bytes.Buffer

    process_id := os.Getpid()
    c1 := exec.Command("lsof", "-p", strconv.Itoa(process_id))
    c2 := exec.Command("wc", "-l")
    c2.Stdin, _ = c1.StdoutPipe()
    c2.Stdout = &fd_count
    _ = c2.Start()
    _ = c1.Run()
    _ = c2.Wait()

    fmt.Println("Number of fds: ", fd_count.String())

    count,_ := strconv.Atoi(fd_count.String())

    if count >= fd_thresh {
    	return true
    }

    return false
}

func RunAnalysis(fd_thresh int, r_thresh int, freq_sec int) {
	for {
		fmt.Println("Analysis running!")

		if checkFileAlert(fd_thresh) {
			fmt.Println("ALERT: number of open files is ", fd_thresh)
		}

		if checkRoutineAlert(r_thresh) {
			fmt.Println("ALERT: number of routines is ", r_thresh)
		}

		time.Sleep(time.Duration(freq_sec) * time.Second)
	}
}
