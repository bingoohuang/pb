package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bingoohuang/pb"
)

func main() {
	num := flag.Int("demo", 1, "select demo number")
	flag.Parse()

	switch *num {
	case 1:
		demo1()
	case 2:
		demo2()
	}
}

func demo1() {
	pool := &pb.Pool{}
	first := pb.Full.New(1000).Set("prefix", "First ").SetMaxWidth(100)
	second := pb.Full.New(1000).Set("prefix", "Second").SetMaxWidth(100)
	third := pb.Full.New(1000).Set("prefix", "Third ").SetMaxWidth(100)
	if err := pool.Start(first, second, third); err == nil {
		defer pool.Stop()
	}

	wg := new(sync.WaitGroup)
	for _, bar := range []*pb.ProgressBar{first, second, third} {
		wg.Add(1)
		go func(cb *pb.ProgressBar) {
			for n := 0; n < 1000; n++ {
				cb.Increment()
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			}
			cb.Finish()
			wg.Done()
		}(bar)
	}
	wg.Wait()
}

// https://github.com/thbar/golang-playground/blob/master/download-files.go
/*
Downloading GB.zip 2.22 MiB / 2.22 MiB [--------------------------------] 100.00% 112.07 KiB p/s 20s
Downloading FR.zip 6.69 MiB / 6.69 MiB [--------------------------------] 100.00% 338.30 KiB p/s 20s
Downloading ES.zip 2.85 MiB / 2.85 MiB [--------------------------------] 100.00% 143.93 KiB p/s 20s

*/
func demo2() {
	pool := &pb.Pool{}
	countries := []string{"GB", "FR", "ES" /*, "DE", "CN", "CA", "ID", "US"*/}
	pbs := make([]*pb.ProgressBar, len(countries))
	for i := 0; i < len(countries); i++ {
		pbs[i] = pb.Full.New(0).Set("prefix", fmt.Sprintf("Downloading %s.zip", countries[i])).SetMaxWidth(100)
	}
	pool.Start(pbs...)
	defer pool.Stop()

	wg := sync.WaitGroup{}

	for i := 0; i < len(countries); i++ {
		wg.Add(1)
		go func(bar *pb.ProgressBar, url string) {
			defer wg.Done()
			defer bar.Finish()

			downloadFromUrl(url, bar)
		}(pbs[i], "http://download.geonames.org/export/dump/"+countries[i]+".zip")
	}

	wg.Wait()
}

func downloadFromUrl(url string, bar *pb.ProgressBar) {
	tokens := strings.Split(url, "/")
	fileName := tokens[len(tokens)-1]
	//fmt.Println("Downloading", url, "to", fileName)

	// TODO: check file existence first with io.IsExist
	output, err := os.Create(fileName)
	if err != nil {
		fmt.Println("Error while creating", fileName, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()
	contentLength := response.Header.Get("Content-Length")
	if contentLength != "" {
		total, _ := strconv.ParseInt(contentLength, 10, 64)
		bar.SetTotal(total)
	}

	defer bar.Finish()

	// create proxy reader
	barReader := bar.NewProxyReader(response.Body)
	if _, err := io.Copy(output, barReader); err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	//fmt.Println(n, "bytes downloaded.")
}
