package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
)

const apiURL = "https://languagetool.org/api/v2/check"
const defaultParams = "text=%s&language=%s&motherTongue=de-DE&enabledCategories=TYPOS%%2CGRAMMAR&enabledOnly=true"
const langDE = "de-DE"
const langEN = "en-US"

const whiteListFile = "whitelist.json"

var whitelist = make([]string, 0)

// Response represents the language tool API response
type Response struct {
	Language struct {
		DetectedLanguage struct {
			Code string
		}
	}
	Matches []struct {
		Offset       int
		Length       int
		Replacements []struct {
			Value string
		}
	}
	Error *appError
}

// CheckText queries the LanguageTools spell check API for typos and spelling mistakes
func CheckText(text string) map[string]string {
	var deRes, enRes Response
	var wg sync.WaitGroup
	wg.Add(2)
	go queryAPI(text, langDE, &deRes, &wg)
	go queryAPI(text, langEN, &enRes, &wg)
	wg.Wait()
	if deRes.Error != nil {
		deRes.Error.log()
		return nil
	} else if enRes.Error != nil {
		enRes.Error.log()
		return nil
	}
	results := make(map[string]string)
	for di := 0; di < len(deRes.Matches); di++ {
		m := deRes.Matches[di]
		do := m.Offset
		ei := sort.Search(len(enRes.Matches), func(i int) bool { return enRes.Matches[i].Offset >= do })
		if ei < len(enRes.Matches) && enRes.Matches[ei].Offset == do {
			// Match was found for both langs
			word := strings.Fields(text[m.Offset:])[0]
			// Remove punctation
			if lastChar := word[len(word)-1]; lastChar == ',' || lastChar == '.' || lastChar == '-' {
				word = word[:len(word)-1]
			}
			// Check whitelist
			wi := sort.SearchStrings(whitelist, strings.ToLower(word))
			if wi < len(whitelist) && strings.EqualFold(whitelist[wi], word) {
				continue
			}
			// Use English replacement if detected
			if deRes.Language.DetectedLanguage.Code == langEN {
				m = enRes.Matches[ei]
			}
			var replacement string
			if len(m.Replacements) > 0 {
				replacement = m.Replacements[0].Value
			} else {
				replacement = "??"
			}
			results[word] = replacement
		}
	}
	if len(results) == 0 {
		return nil
	}
	return results
}

func queryAPI(text, lang string, res *Response, wg *sync.WaitGroup) {
	defer wg.Done()
	params := fmt.Sprintf(defaultParams, url.QueryEscape(text), lang)
	fmt.Println(params)
	postResponse, err := http.Post(apiURL, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(params)))
	if err != nil {
		res.Error = &appError{"POST failed", err}
		return
	}
	defer postResponse.Body.Close()

	body, _ := ioutil.ReadAll(postResponse.Body)
	fmt.Println(string(body))
	if err := json.Unmarshal(body, &res); err != nil {
		res.Error = &appError{"JSON failed", err}
		return
	}
}

// LoadDictionary loads the dictionary from the disk
func LoadDictionary() {
	data, err := ioutil.ReadFile(whiteListFile)
	if err != nil {
		return
	}
	if json.Unmarshal(data, &whitelist) != nil {
		whitelist = make([]string, 0)
	}
}

// AddToDictionary adds a word to the whitelist
func AddToDictionary(word string) bool {
	added := false
	added, whitelist = AddToSortedStringSet(whitelist, word)
	if added {
		// Commit the JSON to disk
		data, err := json.Marshal(whitelist)
		if err != nil {
			return false
		}
		return ioutil.WriteFile(whiteListFile, data, 0644) == nil
	}
	return false
}

// RemoveFromDictionary removes a word from the whitelist
func RemoveFromDictionary(word string) bool {
	removed := false
	removed, whitelist = RemoveFromSortedStringSet(whitelist, word)
	if removed {
		// Commit the JSON to disk
		data, err := json.Marshal(whitelist)
		if err != nil {
			return false
		}
		return ioutil.WriteFile(whiteListFile, data, 0644) == nil
	}
	return false
}

// QueryWhitelist returns the whitelist as a string
func QueryWhitelist() string {
	return fmt.Sprintln(whitelist)
}

type appError struct {
	Message string
	Error   error
}

func (a *appError) log() {
	if a.Error == nil {
		log.Println(a.Message)
	}
	log.Println(a.Message + ": " + a.Error.Error())
}
