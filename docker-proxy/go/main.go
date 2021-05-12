package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var isUdp *bool = flag.Bool("u", false, "use UDP")
var localAddr *string = flag.String("l", ":8080", "local address")
var remoteAddr *string = flag.String("r", "localhost:80", "remote address")

func main() {
	flag.Parse()
	msg := LogInit{
		Local:  *localAddr,
		Remote: *remoteAddr,
	}

	if *isUdp {
		PrintStderr(msg.SetProto("udp"))
		udp(*localAddr, *remoteAddr)
	} else {
		PrintStderr(msg.SetProto("tcp"))
		tcp(*localAddr, *remoteAddr)
	}
}

func tcp(localAddr, remoteAddr string) {
	listener, err := net.Listen("tcp", localAddr)
	if err != nil {
		PrintStderr(LogError{
			Message: fmt.Sprintf("error while listenting to local addr %q", localAddr),
			Error:   err,
		})
		return
	}

	for {
		local, err := listener.Accept()
		if err != nil {
			PrintStderr(LogError{
				Message: fmt.Sprintf("error while accepting connection"),
				Error:   err,
			})
			continue
		}

		defer local.Close()

		go func() {
			remote, err := net.Dial("tcp", remoteAddr)
			if err != nil {
				PrintStderr(LogError{
					Message: fmt.Sprintf("error while dialing remote addr %q", remoteAddr),
					Error:   err,
				})
				return
			}

			PrintStdout(LogNewConnection{
				Protocol: "tcp",
				Client:   local.RemoteAddr().String(),
				Local:    localAddr,
				Source:   remote.LocalAddr().String(),
				Remote:   remoteAddr,
			})

			go io.Copy(remote, local)
			io.Copy(local, remote)

			remote.Close()
			local.Close()
		}()
	}
}

func udp(localAddr, remoteAddr string) {
	connections := map[string]time.Time{}
	ticker := time.NewTicker(5 * time.Second)

	// cleanup connections
	go func() {
		for {
			select {
			case <-ticker.C:
				expires := time.Now().Add(-1 * time.Minute)
				for key, time := range connections {
					if time.Before(expires) {
						delete(connections, key)
					}
				}
			}
		}
	}()

	lAddr, err := net.ResolveUDPAddr("udp", localAddr)
	if err != nil {
		PrintStderr(LogError{
			Message: fmt.Sprintf("error while resolving local addr %q", localAddr),
			Error:   err,
		})
		return
	}

	rAddr, err := net.ResolveUDPAddr("udp", remoteAddr)
	if err != nil {
		PrintStderr(LogError{
			Message: fmt.Sprintf("error while resolving remote addr %q", remoteAddr),
			Error:   err,
		})
		return
	}

	listener, err := net.ListenUDP("udp", lAddr)
	if err != nil {
		PrintStderr(LogError{
			Message: fmt.Sprintf("error while listenting to local addr %q", localAddr),
			Error:   err,
		})
		return
	}

	remote, err := net.DialUDP("udp", nil, rAddr)
	if err != nil {
		PrintStderr(LogError{
			Message: fmt.Sprintf("error while dialing remote addr %q", remoteAddr),
			Error:   err,
		})
		return
	}

	buffer := make([]byte, 1500)
	for {
		n, addr, err := listener.ReadFromUDP(buffer)
		if err != nil {
			PrintStderr(LogError{
				Message: "error while receiving a packet",
				Error:   err,
			})
			continue
		}

		if _, ok := connections[addr.String()]; !ok {
			PrintStdout(LogNewConnection{
				Protocol: "udp",
				Client:   addr.String(),
				Local:    localAddr,
				Source:   remote.LocalAddr().String(),
				Remote:   remoteAddr,
			})

			connections[addr.String()] = time.Now()
		}

		remote.Write(buffer[0:n])
	}
}

type LogNewConnection struct {
	Protocol string `json:"protocol"`
	Client   string `json:"client"`
	Local    string `json:"local"`
	Source   string `json:"source"`
	Remote   string `json:"remote"`
}

type LogInit struct {
	Protocol string `json:"protocol"`
	Local    string `json:"local"`
	Remote   string `json:"remote"`
}

func (l *LogInit) SetProto(proto string) LogInit {
	l.Protocol = proto
	return *l
}

type LogError struct {
	Message string `json:"message"`
	Error   error  `json:"error"`
}

func PrintStdout(data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "{\"time\":\"%s\",\"message\":\"Failed to generate json.\"}", time.Now())
	} else {
		fmt.Fprintln(os.Stdout, string(json))
	}
}

func PrintStderr(data interface{}) {
	json, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "{\"time\":\"%s\",\"message\":\"Failed to generate json.\"}", time.Now())
	} else {
		fmt.Fprintln(os.Stderr, string(json))
	}
}
