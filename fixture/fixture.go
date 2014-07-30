package fixture

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"net/smtp"
	"net/mail"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

type EmailUser struct {
	Username    string
	Password    string
	EmailServer string
	Port        int
}

func sendEmail(send_to string, subj string, content string) {

	emailUser := &EmailUser{"fixture.plumber", "plumber!", "smtp.gmail.com", 587}

	auth := smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer)

	from := mail.Address{"Fixture Plumber", "fixture.plumber.com"}
	to := mail.Address{"Test", send_to}
	title := subj
 
	body := content;
 
	header := make(map[string]string)
	header["From"] = from.String()
	header["To"] = to.String()
	header["Subject"] = title
	header["MIME-Version"] = "1.0"
	header["Content-Type"] = "text/plain; charset=\"utf-8\""
	header["Content-Transfer-Encoding"] = "base64"
 
	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + base64.StdEncoding.EncodeToString([]byte(body))

	err := smtp.SendMail(
		emailUser.EmailServer + ":" + strconv.Itoa(emailUser.Port),
		auth,
		from.Address,
		[]string{to.Address},
		[]byte(message),
	)

	if err != nil {
	 	log.Print("ERROR: attempting to send a mail ", err)
	}
}

type Notifier struct {
	fd_threshold int
	routine_threshold int
	email string
	last_message_fd *time.Time
	last_message_routine *time.Time
}

func (n *Notifier) checkRoutineAlert() {
	num_routines := runtime.NumGoroutine()
	
	fmt.Println("Number of go routines: ", num_routines)
	if num_routines >= n.routine_threshold {
		subject := "ALERT: Routine Threshold Exceeded"
		message := fmt.Sprintf("Number of current Go Routines is %d", num_routines)
		n.AlertUser(subject, message, 1)
	}

}

func (n *Notifier) checkFileAlert() {
	var fd_count bytes.Buffer

    process_id := os.Getpid()
    c1 := exec.Command("lsof", "-p", strconv.Itoa(process_id))
    c2 := exec.Command("wc", "-l")
    c2.Stdin, _ = c1.StdoutPipe()
    c2.Stdout = &fd_count
    _ = c2.Start()
    _ = c1.Run()
    _ = c2.Wait()

    // fmt.Println("Number of fds: ", fd_count.String())

    count,_ := strconv.Atoi(fd_count.String())

    if count >= n.fd_threshold {
		subject := "ALERT: File Descriptor Threshold Exceeded"
		message := fmt.Sprintf("Number of current file descriptors %d", count)
		n.AlertUser(subject, message, 0)
    }

}

func (n *Notifier) AlertUser(subject string, body string, alert_type int) {
	if alert_type == 0 {
		if n.last_message_fd == nil {
			sendEmail(n.email, subject, body)
			x := new(time.Time)
			*x = time.Now().UTC()
			n.last_message_fd = x

			return
		}

		diff := time.Now().UTC().Sub(*n.last_message_fd)
		
		if diff.Hours() >= 1 {
			sendEmail(n.email, subject, body)
		}

	} else {
		if n.last_message_routine == nil {
			sendEmail(n.email, subject, body)
			x := new(time.Time)
			*x = time.Now().UTC()
			n.last_message_routine = x
			fmt.Println("why do i keep getting nil")
			return
		}

		diff := time.Now().UTC().Sub(*n.last_message_routine)
		fmt.Println("The diff is: %s", diff)
		fmt.Println("Hours is: %d", diff.Hours())

		if diff.Hours() >= 1 {
			sendEmail(n.email, subject, body)
		}
	}

}

func RunAnalysis(fd_thresh int, r_thresh int, freq_sec int, email string) {
	if fd_thresh < -1 || r_thresh < -1 {
		fmt.Errorf("Argument less than -1")
		return
	}

	if freq_sec <= 0 {
		fmt.Errorf("Frequency too low: %d", freq_sec)
		return
	}

	n := &Notifier{fd_thresh, r_thresh, email, nil, nil}

	for {
		if n.fd_threshold != -1 {
			n.checkFileAlert()
		}

		if n.routine_threshold != -1 {
			n.checkRoutineAlert()
		}

		time.Sleep(time.Duration(freq_sec) * time.Second)
	}
}
