package main

import (
    "bytes"
    "io/ioutil"
    "strconv"
    "fmt"
)

var info []string

func proc() {
    fullchunk, err := ioutil.ReadFile("sample.torrent")
    maxl := len(fullchunk)
    if err != nil {
        return
    }
    for i, byt := range fullchunk {

        // 105 -> i (integer value)
        if byt == 105 {
            nxt := ""
            for j := i + 1;; j++ {
                seg := fullchunk[j:j+1]
                if bytes.ContainsAny(seg, "1234567890") {
                    fmt.Printf("i>%s\n", seg)
                    nxt += string(seg)
                } else {
                    fmt.Printf("i breaking at %q\n", seg)
                    break
                }
            }
            info = append(info, nxt)
        }

        // 58 -> :
        if byt == 58 {
            prev := ""
            for j := 1; j <= 4; j++ {
                seg := fullchunk[(i-j):((i-j)+1)]
                fmt.Printf(">%+q\n", seg)
                if bytes.ContainsAny(seg, "1234567890") {
                    println("curr", prev)
                    prev = string(seg) + prev
                } else {
                    fmt.Printf("breaking at>%+q\n", seg)
                    break
                }
            }
            if prev != "" {
                stride, err := strconv.Atoi(prev)
                if err != nil {
                    fmt.Printf("%v", err)
                    return
                }
                println("final s:", prev)
                println("final i:", stride)
                if (i+1)+stride > maxl {
                    return
                }

                sc := fullchunk[i+1:(i+1)+stride]
                info = append(info, string(sc))
                // LEAP
                i = (i+1) + stride
                continue
            }
        }

    }
    //fmt.Printf("%s", fullchunk)
}

func main() {
    proc()
    for _, p := range info {
        println(p)
    }
}
