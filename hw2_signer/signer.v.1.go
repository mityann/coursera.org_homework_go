package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)
//
func SingleHash(input int) string{
	result := ""
	i := 0
	data := fmt.Sprintf("%v", input)
	fmt.Println(data)
	fmt.Printf("%v SingleHash data %v\n", i, data)
	md5Date := DataSignerMd5(data)
	fmt.Printf("%v SingleHash md5(data) %v\n", i, md5Date)
	Crc32md5Date := DataSignerCrc32(md5Date)
	fmt.Printf("%v SingleHash crc32(md5(data)) %v\n", i, Crc32md5Date)
	Crc32Date := DataSignerCrc32(data)
	fmt.Printf("%v SingleHash crc32(data) %v\n", i, Crc32Date)
	result = Crc32Date + "~" + Crc32md5Date
	fmt.Printf("%v SingleHash result %v\n", i, result)
	return result
}

func MultiHash(input string) string{
		data := fmt.Sprintf("%v", input)
		result := ""
		for i := 0; i < 6; i++ {
			s:=i+1
			th := DataSignerCrc32(strconv.Itoa(i) + data)
			fmt.Printf("%v MultiHash: crc32(th+step%d)) %d %v\n", data,s,i, th)
			result += th
		}
		fmt.Printf("%v MultiHash result: %v\n", data, result)
		return result
}
func CombineResults(result []string) string{
	sort.Strings(result)
	return strings.Join(result, "_")
}
func main() {
	timeStart := time.Now()
	fmt.Println("0:",MultiHash(SingleHash(0)))
	timeEnd := time.Since(timeStart)
	fmt.Println("MultiHash-SingleHash 0 :", timeEnd)
	timeStart = time.Now()
	fmt.Println("1:",MultiHash(SingleHash(1)))
	timeEnd = time.Since(timeStart)
	fmt.Println("MultiHash-SingleHash 1 :", timeEnd)
	timeStart = time.Now()
	fmt.Println("CombineResults:",CombineResults([]string{MultiHash(SingleHash(0)),MultiHash(SingleHash(1))}))
	timeEnd = time.Since(timeStart)
	fmt.Println("CombineResults (0,1) :", timeEnd)
}