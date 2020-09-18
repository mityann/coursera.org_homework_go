package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

func SingleHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for input := range in {
		data := fmt.Sprintf("%v", input)
		wg.Add(1)
		go func(mu *sync.Mutex, data string, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			resultChCrc32Date := make(chan string)
			go func(out chan string) {
				out <- DataSignerCrc32(data)
			}(resultChCrc32Date)
			resultChmd5 := make(chan string)
			go func() {
				mu.Lock()
				md5 := DataSignerMd5(data)
				mu.Unlock()
				resultChmd5 <- DataSignerCrc32(md5)
			}()
			Crc32Date := ""
			Crc32md5Date := ""
		LOOP:
			for {
				select {
				case Crc32md5Date = <-resultChmd5:
				case Crc32Date = <-resultChCrc32Date:
				}
				if Crc32Date != "" && Crc32md5Date !="" {
					break LOOP
				}
			}
			close(resultChmd5)
			close(resultChCrc32Date)
			out <- Crc32Date + "~" + Crc32md5Date
		}(mu,data,out,wg)
	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	wg := &sync.WaitGroup{}
	for input := range in {
		data := fmt.Sprintf("%v", input)
		wg.Add(1)
		go func(data string, out chan<- interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			arrChanMH := make([]chan string, 6)
			tht := make([]string, 6)
			for i := 0; i < 6; i++ {
				arrChanMH[i] = make(chan string)
				go func(i int) {
					arrChanMH[i] <- DataSignerCrc32(strconv.Itoa(i) + data)
				}(i)
			}
		LOOP:
			for {
				select {
				case tht[0] =  <-arrChanMH[0]:
				case tht[1] =  <-arrChanMH[1]:
				case tht[2] =  <-arrChanMH[2]:
				case tht[3] =  <-arrChanMH[3]:
				case tht[4] =  <-arrChanMH[4]:
				case tht[5] =  <-arrChanMH[5]:
				}
				if tht[0] != "" &&  tht[1] != "" && tht[2] != "" && tht[3] != "" && tht[4] != "" && tht[5] != "" {
					break LOOP
				}
			}
			close(arrChanMH[0])
			close(arrChanMH[1])
			close(arrChanMH[2])
			close(arrChanMH[3])
			close(arrChanMH[4])
			close(arrChanMH[5])
			out <- tht[0]+tht[1]+tht[2]+tht[3]+tht[4]+tht[5]
		}(data, out, wg)
	}
	wg.Wait()
}
func CombineResults(in, out chan interface{}) {
	timeStart := time.Now()
	result := make([]string,0)
	for input := range in {
		data := input.(string)
		result = append(result, data)
	}
	sort.Strings(result)
	out <- strings.Join(result, "_")
	timeEnd := time.Since(timeStart)
	fmt.Println("CombineResults time:", timeEnd)
}

func ExecutePipeline(jobs ...job) {
	out := make(chan interface{})
	for i, jobi := range jobs {
		in := out
		out = make(chan interface{},100)
		go func(val job, in, out chan interface{}) {
			val(in, out)
			close(out)
		}(jobi,in, out)
		if i == 0 {
			close(in)
		}
	}
}
func main() {
	inputData := []int{0, 1, 1, 2, 3, 5, 8}
	hashSignJobs := []job{
		job(func(in, out chan interface{}) {
			fmt.Println("jobi i:", "0","Sender", "in:", in, "out:",out)
			for _, fibNum := range inputData {
				fmt.Printf("Send %v to %v\n",fibNum,out)
				out <- fibNum
			}
		}),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(func(in, out chan interface{}) {
			fmt.Println("jobi i:", "4","Printer", "in:", in, "out:",out)
			fmt.Printf("Total result: %v\n", <-in)
		}),
	}
	ExecutePipeline(hashSignJobs...)
	fmt.Scanln()
}
