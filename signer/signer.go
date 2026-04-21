package main

import (
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func main() {

}

var md5Mutex sync.Mutex

func SingleHash(in, out chan interface{}) {

	for data := range in {
		str := data.(string)
		var crcData, crcMd5 string
		var wg sync.WaitGroup{}
		wg.Add(2)
		go func() {
			defer wg.Done()
			crcData = DataSignerCrc32(str)
		}()

		go func() {
			defer wg.Done()
			md5Mutex.Lock()
			md5Data := DataSignerMd5(str)
			md5Mutex.Unlock()
			crcMd5 = DataSignerCrc32(md5Data)
		}()
		wg.Wait()

		out <- crcData + "~" + crcMd5
	}

}

func MultiHash(in, out chan interface{}) {

	for data := range in {
		h := data.(string)
		m := make([]string, 6)
		var wg sync.WaitGroup
		wg.Add(6)
		for th := 0; th < 6; th++ {
			th := th
			go func () {
				defer wg.Done()
				predRes := strconv.Itoa(th) + h
				m[th] = DataSignerCrc32(predRes)
			}()
		}
		wg.Wait()
		var n string
		for _, finRes := range m {
			n += finRes
		}
		out <- n
	}

}

func CombineResults(in, out chan interface{}) {
	var m []string
	for data := range in {
		m = append(m, data.(string))
	}
	sort.Strings(m)
	out <- strings.Join(m, "_")

}

func ExecutePipeline(jobs ...job) {
	var wg sync.WaitGroup
	var in chan interface{}

	for _, j := range jobs {
		out := make(chan interface{})
		wg.Add(1)
		go func (worker job, in, out chan interface{}){
			defer wg.Done()
			worker(in, out)
			close(out)
		}(j, in, out)
		in = out
	}
	wg.Wait()
}

