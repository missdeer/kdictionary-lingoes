package main

import (
	"io/ioutil"
	"log"
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
	cookie                  = `bwm_uid=YEX4mbDELdWZqf2tjoJ+0Q==; ul=en; PHPSESSID=jtkvf0mivt70ubiiho00op7of5`
	connection              = `keep-alive`
	upgradeInsecureRequests = `1`
)

var (
	client *http.Client
	wg     sync.WaitGroup
)

func downloadDictionary(u string) {
	wg.Add(1)
	defer wg.Done()
	retry := 0
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		log.Println("Could not parse downloadDictionary request:", err)
		return
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", accept)
	req.Header.Set("Accept-Language", acceptLanguage)
	req.Header.Set("Referer", rootURL)
	req.Header.Set("Cookie", cookie)
	req.Header.Set("Connection", connection)
	req.Header.Set("Upgrade-Insecure-Requests", upgradeInsecureRequests)
doPageRequest:
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Could not send downloadDictionary request:", err)
		retry++
		if retry < 3 {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("cannot read downloadDictionary content", err)
		retry++
		if retry < 3 {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}
	// save data to file
	ioutil.WriteFile("file.ld2", data, 0644)
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
doPageRequest:
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Could not send downloadDictionaryPage request:", err)
		retry++
		if retry < 3 {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("cannot read downloadDictionaryPage content", err)
		retry++
		if retry < 3 {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	if resp.StatusCode != 200 {
		log.Println("downloadDictionaryPage status", resp.Status)
		retry++
		if retry < 3 {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	idx := strings.Index(string(data), `href="http://www.lingoes.cn/download/dict/`)
	if idx > 0 {
		endIdx := strings.Index(string(data)[idx+6:], `"`)
		dict := string(data)[idx+6 : idx+5+endIdx]
		log.Println("downloading", dict)
	} else {
		log.Println("can't find dictionary on", u, string(data))
		os.Exit(1)
	}
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
doPageRequest:
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Could not send downloadCategory request:", err)
		retry++
		if retry < 3 {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("cannot read downloadCategory content", err)
		retry++
		if retry < 3 {
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
doPageRequest:
	resp, err := client.Do(req)
	if err != nil {
		log.Println("Could not send downloadRoot request:", err)
		retry++
		if retry < 3 {
			time.Sleep(3 * time.Second)
			goto doPageRequest
		}
		return
	}

	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("cannot read downloadRoot content", err)
		retry++
		if retry < 3 {
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
	client = &http.Client{
		Timeout: 60 * time.Second,
	}

	downloadRoot()
	wg.Wait()
}
