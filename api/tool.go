package api

import (
	"encoding/json"
	"fmt"
	"github.com/sqzxcv/aaPanelGoApi/tool"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"

	"time"
)

type Jar struct {
	lk      sync.Mutex
	cookies map[string][]*http.Cookie
}

func NewJar() *Jar {
	jar := new(Jar)
	jar.cookies = make(map[string][]*http.Cookie)
	return jar
}

// SetCookies handles the receipt of the cookies in a reply for the
// given URL.  It may or may not choose to save the cookies, depending
// on the jar's policy and implementation.
func (jar *Jar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.lk.Lock()
	jar.cookies[u.Host] = cookies
	jar.lk.Unlock()
}

// Cookies returns the cookies to send in a request for the given URL.
// It is up to the implementation to honor the standard cookie use
// restrictions such as in RFC 6265.
func (jar *Jar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

/**
 * 构造带有签名的关联数组
 */
func GetKeyData() map[string]interface{} {
	tempMap := make(map[string]interface{})
	now_time := time.Now().Unix()
	s_now_time := strconv.Itoa(int(now_time))
	log.Println(s_now_time)
	tempMap["request_token"] = tool.GetStringMd5(s_now_time + tool.GetStringMd5(BtConfig.Key))
	tempMap["request_time"] = s_now_time //int(now_time)
	return tempMap
}

//HttpPostCookie

func HttpPostCookie(murl string, cookies map[string]interface{}) (string, error) {

	//PostForm
	jar := NewJar()

	fullpost := make(map[string][]string)
	for name, v := range cookies {
		//change v to string
		s, ok := v.(string)
		if !ok {
			vs, vok := v.(bool)
			if vok {
				s = strconv.FormatBool(vs)
			} else {
				s = strconv.Itoa(v.(int))
			}
		}
		fullpost[name] = []string{s}
	}

	client := http.Client{Transport: nil, CheckRedirect: nil, Jar: jar, Timeout: 99999999999992}
	resp, err := client.PostForm(murl, fullpost)
	if err != nil {
		log.Println(err.Error())
		return "", err
	}

	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Println(map2ToJson(fullpost))
	return string(b), nil
}

func DoReq(murl string, param map[string]interface{}) (result string, err error) {
	client := &http.Client{Timeout: 60 * time.Second}

	form := url.Values{}
	for key, value := range param {
		form.Add(key, fmt.Sprintf("%v", value))
	}

	req, err := http.NewRequest("POST", murl, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

//func DoReq(murl string, param map[string]interface{}) (result string, err error) {
//
//	client := &http.Client{Timeout: 60 * time.Second}
//	reqURL := murl
//	req, err1 := http.NewRequest("POST", reqURL, nil)
//	if err1 != nil {
//		err = err1
//
//		return result, err
//	}
//	//if len(cookieStr) != 0 {
//	//	req.Header.Add("Cookie", cookieStr)
//	//}
//	resp, err2 := client.Do(req)
//	if err2 != nil {
//		err = err2
//		return result, err
//	}
//	defer resp.Body.Close()
//
//	data, err3 := ioutil.ReadAll(resp.Body)
//	if err3 != nil {
//		err = err3
//		return result, err
//	}
//
//	//glog.Info("请求返回cookie信息:", cookie)
//	return string(data), err
//}

func map2ToJson(param map[string][]string /*interface{}*/) string {
	dataType, _ := json.Marshal(param)
	dataString := string(dataType)
	return dataString
}
