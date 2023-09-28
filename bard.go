package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Chatbot struct {
	ReqID          int
	SNlM0e         string
	ConversationID string
	ResponseID     string
	ChoiceID       string
	Client         *http.Client
	Sessionid      string
}

type Response struct {
	Content           string
	ConversationID    string
	ResponseID        string
	FactualityQueries []interface{}
	TextQuery         string
	Choices           []Choice
}

type Choice struct {
	ID      string
	Content string
}

const HOST = "bard.google.com"
const Origin_url = "https://" + HOST
const BASE_URL = "https://" + HOST + "/"
const ASK_URL = BASE_URL + "_/BardChatUi/data/assistant.lamda.BardFrontendService/StreamGenerate"

func getHeader() http.Header {
	return http.Header{
		"Host":          []string{HOST},
		"X-Same-Domain": []string{"1"},
		"User-Agent":    []string{"Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"},
		"Content-Type":  []string{"application/x-www-form-urlencoded;charset=UTF-8"},
		"Origin":        []string{Origin_url},
		"Referer":       []string{BASE_URL},
	}
}

// NewChatbot() function returns a chatbot client which can be used to ask questions.
func NewChatbot(sessionID string) (*Chatbot, error) {

	client := &http.Client{
		Timeout: 100 * time.Second,
	}
	client.Jar, _ = cookiejar.New(nil)
	chatBot := &Chatbot{
		ReqID:          rand.Intn(100000),
		SNlM0e:         "",
		ConversationID: "",
		ResponseID:     "",
		ChoiceID:       "",
		Client:         client,
		Sessionid:      sessionID,
	}
	chatBot.setCookie()
	sNlM0e, err := chatBot.getSNlM0e()
	chatBot.SNlM0e = sNlM0e
	return chatBot, err
}

func (chatBot *Chatbot) setCookie() {
	url, _ := url.Parse(BASE_URL)
	cookie := &http.Cookie{Name: "__Secure-1PSID", Value: chatBot.Sessionid}
	chatBot.Client.Jar.SetCookies(url, []*http.Cookie{cookie})
}

func (c *Chatbot) getSNlM0e() (string, error) {
	resp, err := c.Client.Get(BASE_URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not get google bard")
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Failed to read the response body: ", err)
		return "", err
	}

	// Convert the response body to a string
	bodyString := string(body)

	re := regexp.MustCompile(`SNlM0e":"(.*?)"`)
	match := re.FindStringSubmatch(bodyString)
	if len(match) < 2 {
		return "", fmt.Errorf("Init failed, SNlM0e not found")
	}
	return match[1], nil
}

// Ask() function takes the message string and returns Bard response of type bard.Response
func (c *Chatbot) Ask(message string) (*Response, error) {
	params := url.Values{
		"bl":     {"boq_assistant-bard-web-server_20230514.20_p0"},
		"_reqid": {fmt.Sprintf("%d", c.ReqID)},
		"rt":     {"c"},
	}
	data := url.Values{
		"f.req": []string{fmt.Sprintf(`[null,"[[\"%s\"], null, [\"%s\",\"%s\",\"%s\"]]"]`, message, c.ConversationID, c.ResponseID, c.ChoiceID)},
		"at":    []string{c.SNlM0e},
	}
	url := ASK_URL + "?" + params.Encode()
	req, _ := http.NewRequest("POST", url, strings.NewReader(data.Encode()))
	req.Header = getHeader()
	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	respLines := strings.Split(buf.String(), "\n")
	respJSON := respLines[3]
	var jsonChatData []interface{}
	err = json.Unmarshal(json.RawMessage(respJSON), &jsonChatData)
	if err != nil {
		return nil, err
	}
	jsonChatData = jsonChatData[0].([]interface{})
	jsonChat := jsonChatData[2].(string)
	err = json.Unmarshal(json.RawMessage(jsonChat), &jsonChatData)
	if err != nil {
		return nil, err
	}
	choices := make([]Choice, 3)
	for i, item := range jsonChatData[4].([]interface{}) {
		choices[i] = Choice{ID: item.([]interface{})[0].(string), Content: item.([]interface{})[1].([]interface{})[0].(string)}
	}
	results := &Response{
		Content:        jsonChatData[4].([]interface{})[0].([]interface{})[1].([]interface{})[0].(string),
		ConversationID: jsonChatData[1].([]interface{})[0].(string),
		ResponseID:     jsonChatData[1].([]interface{})[1].(string),
		// FactualityQueries: jsonChatData[3].([]interface{}),
		// TextQuery:         jsonChatData[2].([]interface{})[0].([]interface{})[0].(string),
		Choices: choices,
	}
	c.ConversationID = results.ConversationID
	c.ResponseID = results.ResponseID
	c.ChoiceID = results.Choices[0].ID
	c.ReqID += 100000
	return results, nil
}
