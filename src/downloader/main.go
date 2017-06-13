package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

const (
	rootURL                 = `http://www.lingoes.cn/zh/dictionary/index.html`
	categoryURLPattern      = `dict_cata\.php\?cata=\d+(\.\w+)?`
	dictionaryURLPattern    = `dict_down\.php\?id=[0-9A-Z]{32,}`
	dictionaryLinkPattern   = `http://www\.lingoes\.cn/download/dict/[^\.]+\.ld2`
	host                    = `www.lingoes.cn`
	userAgent               = ` Mozilla/5.0 (Windows NT 6.1; WOW64; rv:53.0) Gecko/20100101 Firefox/53.0`
	accept                  = `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`
	acceptLanguage          = `en-US,en;q=0.5`
	cookie                  = `bwm_uid=R0aTmS7Eskuxqf2tXoKq0Q==`
	connection              = `keep-alive`
	upgradeInsecureRequests = `1`
	retryCount              = math.MaxInt32
)

type Semaphore struct {
	c chan int
}

func NewSemaphore(n int) *Semaphore {
	s := &Semaphore{
		c: make(chan int, n),
	}
	return s
}

func (s *Semaphore) Acquire() {
	s.c <- 0
}

func (s *Semaphore) Release() {
	<-s.c
}

type Dictionary map[string]string

var (
	wg                sync.WaitGroup
	dictionaries      = make(map[string]Dictionary)
	dictionariesMutex sync.Mutex
	sema              = NewSemaphore(3)
)

func exists(f string) bool {
	stat, err := os.Stat(f)
	if err == nil {
		if stat.Mode()&os.ModeType == 0 {
			return true
		}
		return false
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func downloadDictionary(u string, m Dictionary) {
	sema.Acquire()
	defer func() {
		sema.Release()
		wg.Done()
	}()
	if !exists(m["id"]) {
		os.MkdirAll(m["id"], 0777)
	}
	slashIdx := strings.LastIndex(u, "/")
	name := u[slashIdx:]
	filePath := m["id"] + name
	if !exists(filePath) {
		retry := 0
		req, err := http.NewRequest("GET", u, nil)
		if err != nil {
			log.Println("Could not parse downloadDictionary request:", err)
			return
		}

		req.Header.Set("User-Agent", userAgent)
		req.Header.Set("Accept", accept)
		req.Header.Set("Accept-Language", acceptLanguage)
		req.Header.Set("Referer", m["referer"])
		req.Header.Set("Cookie", cookie)
		req.Header.Set("Connection", connection)
		req.Header.Set("Upgrade-Insecure-Requests", upgradeInsecureRequests)
		client := &http.Client{}
	doPageRequest:
		resp, err := client.Do(req)
		if err != nil {
			log.Println("Could not send downloadDictionary request:", err)
			retry++
			if retry < retryCount {
				time.Sleep(3 * time.Second)
				goto doPageRequest
			}
			return
		}

		defer resp.Body.Close()
		if resp.StatusCode != 200 {
			// log.Println("downloadDictionaryPage status", resp.Status)
			retry++
			if retry < retryCount {
				time.Sleep(3 * time.Second)
				goto doPageRequest
			}
			return
		}

		if resp.Header.Get("Content-Type") != `application/octet-stream` {
			retry++
			if retry < retryCount {
				time.Sleep(3 * time.Second)
				goto doPageRequest
			}
			return
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("cannot read downloadDictionary content", err)
			retry++
			if retry < retryCount {
				time.Sleep(3 * time.Second)
				goto doPageRequest
			}
			return
		}
		// save data to file
		ioutil.WriteFile(filePath, data, 0644)
		log.Println(u, "is saved to", m["id"]+name)
	}

	readme := m["id"] + "/readme.txt"
	if !exists(readme) {
		var s []string
		for k, v := range m {
			s = append(s, fmt.Sprintf("%s: %s", k, v))
		}
		ioutil.WriteFile(readme, []byte(strings.Join(s, "\n")), 0644)
	}
}

func downloadDictionaryPage(u string) {
	wg.Add(1)
	defer wg.Done()
	retry := 0
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println("Could not parse downloadDictionaryPage request:", err)
		return
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", accept)
	req.Header.Set("Accept-Language", acceptLanguage)
	req.Header.Set("Referer", rootURL)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Connection", connection)
	req.Header.Set("Upgrade-Insecure-Requests", upgradeInsecureRequests)
	client := &http.Client{}
doPageRequest:
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Could not send downloadDictionaryPage request:", err)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		// log.Println("downloadDictionaryPage status", resp.Status)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("cannot read downloadDictionaryPage content", err)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	content := string(data)
	idx := strings.Index(content, `href="http://www.lingoes.cn/download/dict/ld2`)
	if idx < 0 {
		log.Println("can't find dictionary on", u)
		return
	}
	endIdx := strings.Index(content[idx+6:], `"`)
	if endIdx < 0 {
		log.Println("can't find dictionary on", u)
		return
	}
	dict := content[idx+6 : idx+6+endIdx]
	m := make(Dictionary)
	m["referer"] = u

	idIdx := strings.Index(u, "id=")
	id := u[idIdx+3:]
	m["id"] = id

	titlePattern := `<div title="ID: [0-9A-Z]{32,}" style="font\-size: 16px; color:#07519A;"><b>([^<]+)</b>`
	regexTitle := regexp.MustCompile(titlePattern)
	ss := regexTitle.FindAllSubmatch(data, -1)
	for _, match := range ss {
		m["title"] = string(match[1])
		break
	}

	descriptionLeadings := `<div style="margin: 10px 0 10px 0; line-height: 130%">`
	idx = strings.Index(content, descriptionLeadings)
	if idx > 0 {
		endIdx = strings.Index(content[idx+len(descriptionLeadings):], "</div>")
		if endIdx > 0 {
			m["description"] = content[idx+len(descriptionLeadings) : idx+len(descriptionLeadings)+endIdx-1]
		}
	}

	languageLeadings := `<td width="80" valign="top"><font color="#333"><b>语言:</b></font></td>`
	idx = strings.Index(content, languageLeadings)
	if idx > 0 {
		p := `<td valign="top">([^<]+)</td>`
		regexLang := regexp.MustCompile(p)
		ss := regexLang.FindAllSubmatch(data, -1)
		for _, match := range ss {
			m["language"] = strings.Trim(string(match[1]), "\r\n ")
			break
		}
	}

	dictionariesMutex.Lock()
	dictionaries[dict] = m
	dictionariesMutex.Unlock()
}

func downloadCategory(u string) {
	wg.Add(1)
	defer wg.Done()
	retry := 0
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println("Could not parse downloadCategory request:", err)
		return
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", accept)
	req.Header.Set("Accept-Language", acceptLanguage)
	req.Header.Set("Referer", rootURL)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Connection", connection)
	req.Header.Set("Upgrade-Insecure-Requests", upgradeInsecureRequests)
	client := &http.Client{}
doPageRequest:
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Could not send downloadCategory request:", err)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		// log.Println("downloadDictionaryPage status", resp.Status)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("cannot read downloadCategory content", err)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	regexDict := regexp.MustCompile(dictionaryURLPattern)
	ss := regexDict.FindAllSubmatch(data, -1)
	for _, match := range ss {
		dict := string(match[0])
		log.Println("found dictionary", dict, "on category", u)
		go downloadDictionaryPage(`http://www.lingoes.cn/zh/dictionary/` + dict)
	}
}

func downloadRoot() {
	wg.Add(1)
	defer wg.Done()
	retry := 0
	req, err := http.NewRequest("GET", rootURL, nil)
	if err != nil {
		log.Println("Could not parse downloadRoot request:", err)
		return
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", accept)
	req.Header.Set("Accept-Language", acceptLanguage)
	req.Header.Set("Referer", rootURL)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Connection", connection)
	req.Header.Set("Upgrade-Insecure-Requests", upgradeInsecureRequests)
	client := &http.Client{}
doPageRequest:
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Could not send downloadRoot request:", err)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		// log.Println("downloadDictionaryPage status", resp.Status)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("cannot read downloadRoot content", err)
		retry++
		if retry < retryCount {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	regexCategory := regexp.MustCompile(categoryURLPattern)
	ss := regexCategory.FindAllSubmatch(data, -1)
	for _, match := range ss {
		cate := string(match[0])
		log.Println("found category", cate, "on root")
		go downloadCategory(`http://www.lingoes.cn/zh/dictionary/` + cate)
	}

	regexDict := regexp.MustCompile(dictionaryURLPattern)
	ss = regexDict.FindAllSubmatch(data, -1)
	for _, match := range ss {
		dict := string(match[0])
		log.Println("found dictionary", dict, "on root")
		go downloadDictionaryPage(`http://www.lingoes.cn/zh/dictionary/` + dict)
	}
}

func main() {
	downloadRoot()
	wg.Wait()

	log.Println("total dictionary count:", len(dictionaries))
	wg.Add(len(dictionaries))
	for dict, m := range dictionaries {
		go downloadDictionary(dict, m)
	}
	wg.Wait()
}
