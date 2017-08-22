package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

type File struct {
	Path   string
	Length uint64
}

type Tracker struct {
	Name         string
	Announce     string
	AnnList      []string
	Comment      string
	CreatedBy    string
	PieceLength  uint32
	CreationDate uint64
	Files        []*File
	info_hash    string
}

var pieces []string
var core Tracker

func proc() {
	fullchunk, err := ioutil.ReadFile("sample.torrent")
	maxl := len(fullchunk)
	if err != nil {
		return
	}

	currKey := ""
	for i, byt := range fullchunk {

		// 105 -> i (integer value)
		if byt == 105 {
			nxt := ""
			for j := i + 1; ; j++ {
				seg := fullchunk[j : j+1]
				if bytes.ContainsAny(seg, "1234567890") {
					nxt += string(seg)
				} else {
					break
				}
			}

			if nxt == "" {
				continue
			}

			val, err := strconv.Atoi(nxt)
			if err != nil {
				println("ATOI err w: ", nxt)
				continue
			}

			switch currKey {
			case "creation date":
				core.CreationDate = uint64(val)
				continue
			case "piece length":
				core.PieceLength = uint32(val)
				continue
			case "pieces":
				return
			case "length":
				core.Files[len(core.Files)-1].Length = uint64(val)
				continue
			}
		}

		// 58 -> :
		if byt == 58 {
			prev := ""
			for j := 1; j <= 8; j++ {
				seg := fullchunk[(i - j):((i - j) + 1)]
				if bytes.ContainsAny(seg, "1234567890") {
					prev = string(seg) + prev
				} else {
					break
				}
			}
			if prev != "" {
				stride, err := strconv.Atoi(prev)
				if err != nil {
					fmt.Printf("%v", err)
					return
				}
				if (i+1)+stride > maxl {
					return
				}

				sc := strings.ToLower(string(fullchunk[i+1 : (i+1)+stride]))
				// LEAP
				i = (i + 1) + stride

				switch sc {
				case "name":
					currKey = "name"
					continue
				case "announce":
					currKey = "announce"
					continue
				case "announce-list":
					currKey = "announce-list"
					continue
				case "comment":
					currKey = "comment"
					continue
				case "created by":
					currKey = "created by"
					continue
				case "creation date":
					currKey = "creation date"
					continue
				case "piece length":
					currKey = "piece length"
					continue
				case "pieces":
					currKey = "pieces"
					continue
				case "files":
					currKey = "files"
					continue
				case "length":
					currKey = "length"
					core.Files = append(core.Files, &File{})
					continue
				case "path":
					currKey = "path"
					continue
				case "info":
					remainder := fullchunk[i+1 : maxl-2]
					hs := sha1.New()
					hs.Write(remainder)
					csum := hs.Sum(nil)
					core.info_hash = fmt.Sprintf("% x", csum)
					continue
				}

				switch currKey {
				case "name":
					core.Name = sc
					continue
				case "announce":
					core.Announce = sc
					continue
				case "announce-list":
					core.AnnList = append(core.AnnList, sc)
					continue
				case "comment":
					core.Comment = sc
					continue
				case "created by":
					core.CreatedBy = sc
					continue
				case "path":
					core.Files[len(core.Files)-1].Path = sc
					continue
				}
				continue
			}
		}

	}
}

func main() {
	t1 := time.Now()
	core = Tracker{}
	proc()
	fmt.Printf("\nunder %s\n", time.Since(t1))
	fmt.Printf("%+v", core)
	for _, fl := range core.Files {
		fmt.Printf("%+v\n", fl)
	}
}
