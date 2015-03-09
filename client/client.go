package main

import(
    "fmt"
    "math"
    "net"
    "os"
    "time"
)
func main(){
    out := ""
    lastFlush := 0
    sinceLast := 0

    conn, _ := net.DialTCP("tcp", nil, &net.TCPAddr{net.ParseIP(os.Args[1]), 8080, ""})
    conn.Write([]byte("GET / HTTP 1.1\r\n\r\n"))

    //synchronize
    var buff = make([]byte, 3000)
    i, err := conn.Read(buff)
    lastTime := time.Now()
    startTime := lastTime
    os.Stderr.Write(buff[:i])
    for err = nil; err == nil; {
        i, err = conn.Read(buff)
        cTime := time.Now()
        os.Stderr.Write(buff[:i])

        //Determine the number of time cycles passed
        zeros := int(cTime.Sub(lastTime) / 5 / time.Millisecond)
        for i := 0; i < zeros; i++{
            out += "0"
            sinceLast++
            if sinceLast == 16{
                os.Stdout.Write(binconvert(demanchester(out[lastFlush:lastFlush + 16])))
                lastFlush += 16
                sinceLast = 0
            }
        }
        if err == nil{
            out += "1"
            sinceLast++
            if sinceLast == 16{
                os.Stdout.Write(binconvert(demanchester(out[lastFlush:lastFlush + 16])))
                lastFlush += 16
                sinceLast = 0
            }
        }
        lastTime = cTime
    }
    fmt.Fprintf(os.Stderr, "Rate: %fbps\n", float64(lastFlush) / float64(time.Now().Sub(startTime)) *  float64(time.Second))
    //fmt.Println("\nThe message was: " + out)
    //fmt.Println("The message really was: " + string(binconvert(out)))
}

func demanchester(in string) (string){
    ret := ""
    for b := 0; b < len(in); b += 2{
        ret += string(in[b])
    }
    return ret
}

func binconvert(in string) ([]byte){
    l := len(in) / 8
    var ret = make([]byte, l)

    for i := 0; i < l; i++{
        for j := 0; j < 8; j++{
            if in[i * 8 + (7 - j)] == 0x31{
                ret[i] += byte(math.Pow(2, float64(j)))
            }
        }
    }
    return ret
}

