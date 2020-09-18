package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)
//
type tjob func(in string) string

func chanDataSignerMd5(data string, out chan string)  {
	out<-DataSignerMd5(data)
}
func chanDataSignerCrc32(data string, out chan string)  {
	out<-DataSignerCrc32(data)
}
func SingleHash(input string) string{
	result := ""
	data := fmt.Sprintf("%v", input)

	arrChan00 := make(chan string,0)
	arrChan01 := make(chan string,0)
	arrChan02 := make(chan string,0)

	go chanDataSignerMd5(data, arrChan00)
	go chanDataSignerCrc32(data, arrChan01)

	md5Date :=""
	Crc32Date := ""
	Crc32md5Date := ""
LOOP:
	for {
		select {
		case md5 := <-arrChan00:
			//md5Date := DataSignerMd5(data)
			md5Date = md5
			go chanDataSignerCrc32(md5, arrChan02)
		case Crc32Date = <-arrChan01:
			//Crc32Date := DataSignerCrc32(data)
		case Crc32md5Date = <-arrChan02:
			//Crc32md5Date := DataSignerCrc32(md5Date)
		}
		if Crc32Date != "" && Crc32md5Date !="" {
			result = Crc32Date + "~" + Crc32md5Date
			break LOOP
		}
	}
	close(arrChan00)
	close(arrChan01)
	close(arrChan02)

	fmt.Printf("%v SingleHash md5(data) %v\n", data, md5Date)
	fmt.Printf("%v SingleHash crc32(md5(data)) %v\n", data, Crc32md5Date)
	fmt.Printf("%v SingleHash crc32(data) %v\n", data, Crc32Date)
	fmt.Printf("%v SingleHash result %v\n", data, result)
	return result
}

func MultiHash(input string) string{
	data := fmt.Sprintf("%v", input)
	result := ""
	arrChanMH := make([]chan string, 6) // массив каналов
	tht := make([]string, 6)
	for i := 0; i < 6; i++ {
		arrChanMH[i] = make(chan string,3)
		go chanDataSignerCrc32(strconv.Itoa(i) + data, arrChanMH[i])
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

	for i := 0; i < 6; i++ {
		s:=i+1
		th := tht[i]
		fmt.Printf("%v MultiHash: crc32(th+step%d)) %d %v\n", data,s,i, th)
		result += th
	}
	fmt.Printf("%v MultiHash result: %v\n", data, result)
	return result
}
func CombineResults(input string) string{
	result := make([]string,0)
	result = append(result, input)
	sort.Strings(result)
	return strings.Join(result, "_")
}
func ExecutePipeline(tjobs ...tjob) {
	out := "0"
	for i, jobi := range tjobs {
		fmt.Printf("jobi adress: %v %v\n", i, jobi)
		out = jobi(out)
	}
}
func main() {
	timeStart := time.Now()
	hashSignJobs := []tjob{
		tjob(SingleHash),
		tjob(MultiHash),
		tjob(CombineResults),
	}
	ExecutePipeline(hashSignJobs...)
	timeEnd := time.Since(timeStart)
	fmt.Println(" time :", timeEnd)
}