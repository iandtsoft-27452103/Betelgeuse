package shogi

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
)

type Communicate struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser
	reader *bufio.Reader
}

func InitCommunicate() Communicate {
	var c Communicate
	var err error
	c.cmd = exec.Command("python", "-u", "python0/Gift.py")
	c.stdin, err = c.cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	c.stdout, err = c.cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	c.stderr, err = c.cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	c.reader = bufio.NewReader(c.stdout)

	return c
}

func Boot(cm *Communicate) {

	/*var wg sync.WaitGroup
	wg.Add(2)

	// stdout ログ読み取り (非同期)
	go func() {
		defer wg.Done()
		scanner := bufio.NewScanner(cm.stdout)
		for scanner.Scan() {
			fmt.Println("[PY-OUT]", scanner.Text())
		}
	}()

	wg.Wait()

	// コマンド終了待機
	if err := cm.cmd.Wait(); err != nil {
		fmt.Println("python exited with error:", err)
		return
	}*/

	if err := cm.cmd.Start(); err != nil {
		log.Fatal(err)
	}

	// Python から受信
	resp, _ := cm.reader.ReadString('\n')
	fmt.Println("Response:", resp)

	//reader.Discard(len(resp))

	resp, _ = cm.reader.ReadString('\n')
	fmt.Println("Response:", resp)

	resp, _ = cm.reader.ReadString('\n')
	fmt.Println("Response:", resp)

	//fmt.Println("aaa.")
}

func ThrowRequest(cm *Communicate, msg string) {
	fmt.Fprintln(cm.stdin, msg)
}

func ReceiveResponse(cm *Communicate) string {
	resp, _ := cm.reader.ReadString('\n')
	//resp, _ := cm.reader.ReadString('\n')
	//fmt.Println("Response:", resp)
	return resp
}
