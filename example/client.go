package main

import (
    "fmt"
    "net"
    "log"
    "os"
    "bufio"
    "time"
    "flag"
    "os/user"
    "os/exec"
    "strings"
    "encoding/json"
)

var running bool;  // global variable if client is running
var debug = flag.Bool("d", false, "enable debug mode ( display debugging information )")
var server = flag.String("s", "127.0.0.1", "server name to connect to")
var port = flag.String("p", "9988" , "port number to connect to")
var current_user,_ = user.Current()

// func Log(v ...): loging. give log information if debug is true
type Action struct {
  Name string
  Target string
}

func Log(v ...string) {
    if *debug == true {
        ret := fmt.Sprint(v);
        log.Printf("CLIENT: %s", ret);
    }
}

// func test(): testing for error
func test(err error, mesg string) {
    if err!=nil {
         log.Printf("CLIENT: ERROR: %s : %s", mesg, err);
         os.Exit(-1);
    } else {
        Log("Ok: ", mesg);
    }
}

// read from connection and return true if ok
func Read(con net.Conn) string{
    var buf [4048]byte;
    _, err := con.Read(buf[0:4048]);
    if err!=nil {
        con.Close();
        running=false;
        return "Error in reading!";
    }
    str := string(buf[0:4048]);
    fmt.Println();
    return string(str);
}

// clientsender(): read from stdin and send it via network
func clientsender(cn net.Conn) {
    reader := bufio.NewReader(os.Stdin);
    for {
        input, err := reader.ReadBytes('\n')
        if err == nil  && len(input) > 1 {
            tokens := strings.Fields(string(input[0:len(input)-1]))

            if tokens[0] == "/quit" {
                cn.Write([]byte("is leaving..."))
                running = false
                break
            } else if tokens[0] == "/command" {
                if len(tokens) > 1 {
                    out, err := exec.Command(tokens[1], tokens[2:]...).Output()
                    if err != nil {
                        fmt.Printf("Error: %s\n", err)
                    } else {
                        cn.Write(out) // send output to server
                    }
                } else {
                    fmt.Printf("Usage:\n\t/command <exec> <arguments>\n\tEx: /command ls -l -a\n\n")
                }
                continue
            }
            Log("clientsender(): send: ", string(input[0:len(input)-1]))
            cn.Write(input[0:len(input)-1])
        }
    }
}

// clientreceiver(): wait for input from network and print it out
func clientreceiver(cn net.Conn) {
    for running {
        buf := make([]byte,2048)
        _,err := cn.Read(buf);
        if err != nil {
          panic(err)
        }
        var rec Action
        json.Unmarshal([]byte(strings.Trim(string(buf), "\x00")), &rec)
        //fmt.Printf(rec.Name)
        fmt.Printf("(%s)\n", string(buf))
        fmt.Printf("\n[%v]", rec.Name)
        fmt.Printf(" %v\n---------\n", rec.Target)
    }
}

func usage() {
    //fmt.Fprintf(os.Stderr, "usage: client 192.168.1.1:9999\n")
    flag.PrintDefaults()
    os.Exit(2)
}

func main() {
    flag.Usage = usage;
    flag.Parse();
    fmt.Print("Hello ")
    fmt.Print(current_user.Name)
    fmt.Print("\nWho lives in ")
    fmt.Print(current_user.HomeDir)
    fmt.Print(" \n")

    running = true;
    Log("main(): start ");

    destination := fmt.Sprintf("%s:%s", *server,*port);
    fmt.Println("Connecting to: ", destination);

    Log("main(): connecto to ", destination);
    cn, err := net.Dial("tcp", destination);
    test(err, "dialing");
    defer cn.Close();
    Log("main(): connected ");

    // get the user name
    fmt.Print("Please give your name: ");
    reader := bufio.NewReader(os.Stdin);
    name, _ := reader.ReadBytes('\n');

    fmt.Println("Enter /quit to quit");

    //cn.Write(strings.Bytes("User: "));
    cn.Write(name[0:len(name)-1]);

    // start receiver and sender
    Log("main(): start receiver");
    go clientreceiver(cn);
    Log("main(): start sender");
    go clientsender(cn);

    // wait for quiting (/quit). run until running is true
    for ;running; {
        time.Sleep(1*1e9);
    }
    Log("main(): stopped");
}
