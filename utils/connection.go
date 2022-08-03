package utils

import (
	"bufio"
	// "fmt"
	"io"
	// "log"
	"net"
	"strings"
	"github.com/beego/beego/v2/core/logs"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"os"
)

type Connection struct {
	*ssh.Client
	password string
}

func Connect(addr, user, password string) (*Connection, error) {
	sshConfig := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.HostKeyCallback(func(hostname string, remote net.Addr, key ssh.PublicKey) error { return nil }),
	}

	conn, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, err
	}

	return &Connection{conn, password}, nil

}
func (conn *Connection) SendCommands(cmds ...string) ([]byte, error) {
	session, err := conn.NewSession()
	if err != nil {
		logs.Error(err)
		return nil,err
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	err = session.RequestPty("xterm", 80, 40, modes)
	if err != nil {
		return []byte{}, err
	}

	in, err := session.StdinPipe()
	if err != nil {
		// log.Fatal(err)
		logs.Error(err)
		return nil,err
	}

	out, err := session.StdoutPipe()
	if err != nil {
		// log.Fatal(err)
		logs.Error(err)
		return nil,err
	}

	var output []byte

	go func(in io.WriteCloser, out io.Reader, output *[]byte) {
		var (
			line string
			r    = bufio.NewReader(out)
		)
		for {
			b, err := r.ReadByte()
			if err != nil {
				break
			}

			*output = append(*output, b)

			if b == byte('\n') {
				line = ""
				continue
			}

			line += string(b)

			if strings.HasPrefix(line, "[sudo] password for ") && strings.HasSuffix(line, ": ") {
				_, err = in.Write([]byte(conn.password + "\n"))
				if err != nil {
					break
				}
			}
		}
	}(in, out, &output)

	cmd := strings.Join(cmds, "; ")
	_, err = session.Output(cmd)
	if err != nil {
		return []byte{}, err
	}

	return output, nil

}
///create sfpt client 
func(conn *Connection)Createsfptclient()(*sftp.Client,error){
	sftpClient,serr:=sftp.NewClient(conn.Client)
	if(serr!=nil){
		return nil,serr
	}
	return sftpClient,nil
}
///download file from sftp
func(conn *Connection)Downloadfile(sc *sftp.Client, remoteFile, localFile string)(err error){
	srcFile, err := sc.OpenFile(remoteFile, (os.O_RDONLY))

	if err != nil {
		logs.Error("Unable to open remote file: %v\n", err)
        // fmt.Fprintf(os.Stderr, "Unable to open remote file: %v\n", err)
        logs.Error(err)
		return err
    }
	defer srcFile.Close()
	dstFile, err := os.Create(localFile)
   
	if err != nil {
		logs.Error("Unable to open local file: %v\n", err)
        // fmt.Fprintf(os.Stderr, "Unable to open local file: %v\n", err)
        logs.Error(err)
		return err
    }
	defer dstFile.Close()
    _, err = io.Copy(dstFile, srcFile)
    if err != nil {
		logs.Error(err)
        // fmt.Fprintf(os.Stderr, "Unable to download remote file: %v\n", err)
        return err
    }
    // fmt.Fprintf(os.Stdout, "%d bytes copied\n", bytes)
    // logs.Info("%d bytes copied\n", bytes)
	return nil
}