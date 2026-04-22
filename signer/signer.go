package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func main() {

}

var md5Mutex sync.Mutex

func SingleHash(in, out chan interface{}) {
	var wg sync.WaitGroup
	for data := range in {
		wg.Add(1)
		go func(data interface{}) {
			defer wg.Done()
			str := fmt.Sprint(data)
			var crcData, crcMd5 string
			var wg2 sync.WaitGroup
			wg2.Add(2)
			go func() {
				defer wg2.Done()
				crcData = DataSignerCrc32(str)
			}()

			go func() {
				defer wg2.Done()
				md5Mutex.Lock()
				md5Data := DataSignerMd5(str)
				md5Mutex.Unlock()
				crcMd5 = DataSignerCrc32(md5Data)
			}()
			wg2.Wait()

			out <- crcData + "~" + crcMd5
		}(data)

	}
	wg.Wait()
}

func MultiHash(in, out chan interface{}) {
	var wg sync.WaitGroup
	for data := range in {
		wg.Add(1)
		go func(data interface{}) {
			defer wg.Done()
			h := data.(string)
			m := make([]string, 6)
			var wg2 sync.WaitGroup
			wg2.Add(6)
			for th := 0; th < 6; th++ {
				th := th
				go func() {
					defer wg2.Done()
					predRes := strconv.Itoa(th) + h
					m[th] = DataSignerCrc32(predRes)
				}()
			}
			wg2.Wait()
			var result string
			for _, part := range m {
				result += part
			}
			out <- result
		}(data)
	}
	wg.Wait()
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
		out := make(chan interface{}, 10)
		wg.Add(1)
		go func(worker job, in, out chan interface{}) {
			defer wg.Done()
			worker(in, out)
			close(out)
		}(j, in, out)
		in = out
	}
	wg.Wait()
}
