package main

import (
    "fmt"
    "io/ioutil"
    "net"
    "os"
    "time"
)

func handle(conn *net.TCPConn, data [][]byte, encoded string) {
    //Read the "GET"
    var buff []byte
    conn.Read(buff)
    conn.SetNoDelay(true)

    //Send synchoronization message
    s := []byte("HTTP/1.1 200 OK\r\n\r\n")
    conn.Write(s)

    //Send the data
    d := 0
    for b := range(encoded){
        if encoded[b] == 0x30{
            time.Sleep(5 * time.Millisecond)
            continue
        }
        time.Sleep(1 * time.Millisecond)
        conn.Write(data[d])
        d++
    }
    conn.CloseWrite()
    conn.CloseRead()
}
func encode(m1 []byte) (string/*, int*/){
    ret := ""
    //ret, popcount := "", 0
    for b := range(m1){
        s := fmt.Sprintf("%08b", m1[b])
        /*for bit := 0; bit < 8; bit++{
            if s[bit] == 0x31{
                popcount++
            }
        }*/
        ret += s
    }
    return ret/*, popcount*/
}
func manchester(in string) (string){
    ret := ""
    for b := range(in){
        if in[b] == 0x30{
            ret += "01"
        } else {
            ret += "10"
        }
    }
    return ret
}
func main() {
    //Steg'd data
    msgbytes, _ := ioutil.ReadAll(os.Stdin)
    //themessaget, popcount := encode(msgbytes)
    themessage := manchester(encode(msgbytes))
    popcount := len(themessage) / 2
    fmt.Println(themessage)

    //The file to host
    f, _ := os.Open(os.Args[1])
    fdata, _ := ioutil.ReadAll(f)
    f.Close()
    chunksize := len(fdata) / popcount
    var hosted [][]byte
    for chunk := 0; chunk < popcount - 1; chunk++{
        hosted = append(hosted, fdata[chunk * chunksize : chunksize * (chunk + 1)])
    }
    hosted = append(hosted, fdata[(popcount - 1) * chunksize:])

    //Serve
    listener, _ := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP("0.0.0.0"), 8080, ""})
    for {
        conn, _ := listener.AcceptTCP()
        go handle(conn, hosted, themessage)
    }
}
