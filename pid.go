package wzrpc

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/TwiN/go-color"
)

type PID int

func (pid PID) BumpWindows() bool {
	pr := "Realtime" // "Realtime" "High priority" "Above normal"
	whereClause := fmt.Sprintf(`processid='%v'`, pid)
	a := []string{"process", "where", whereClause, "CALL", "setpriority", pr}
	c := exec.Command("wmic", a...)
	return c.Run() == nil
}

func (pid PID) BumpLinux() bool {
	c := exec.Command(fmt.Sprintf(`renice -15 -p %v`, pid))
	return c.Run() == nil
}

func (pid PID) log(success bool) {
	var colour, sign = color.Yellow, " :) "
	if !success {
		colour, sign = color.Red, " :| "
	}
	log.Println(sign + "Reniced: " + color.Ize(colour, pid))
}

func (pid PID) Bump() bool {
	ok := pid.BumpLinux() || pid.BumpWindows()
	pid.log(ok)
	return ok
}
