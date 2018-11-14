package api

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"goophr/concierge/common"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type payload struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

type document struct {
	Doc   string `json:"-"`
	Title string `json:"title"`
	DocID string `json:"DocID"`
}

type token struct {
	Line   string `json:"-"`
	Token  string `json:"token"`
	Title  string `json:"title"`
	DocID  string `json:"doc_id"`
	LIndex int    `json:"line_index"`
	Index  int    `json:"token_index"`
}

type dMsg struct {
	DocID string
	Ch    chan document
}

type lMsg struct {
	LIndex int
	DocID  string
	Ch     chan string
}

type lMeta struct {
	LIndex int
	DocID  string
	Line   string
}

type dAllMsg struct {
	Ch chan []document
}

var done chan bool

var dGetCh chan dMsg

var lGetCh chan lMsg

var lStoreCh chan lMeta

var iAddCh chan token

var dStoreCh chan document

var dProcessCh chan document

var dGetAllCh chan dAllMsg

var pProcessCh chan payload

func StartFeederSystem() {
	done = make(chan bool)

	dGetCh = make(chan dMsg, 8)
	dGetAllCh = make(chan dAllMsg)

	iAddCh = make(chan token, 8)
	pProcessCh = make(chan payload, 8)

	dStoreCh = make(chan document, 8)
	dProcessCh = make(chan payload, 8)
	lGetCh = make(chan lMsg)
	lStoreCh = make(chan lMata, 8)

	for i := 0; i < 4; i++ {
		go indexAdder(iAddCh, done)
		go docProcessor(pProcessCh, dStoreCh, dProcessCh, done)
		go indexProcessor(dProcessCh, lStoreCh, iAddCh, done)
	}

	go docStore(dStoreCh, dGetCh, dGetAllCh, done)
	go lineStore(lStoreCh, lGetch, done)
}

func indexAdder(ch chan token, done chan bool) {
	for {
		select {
		case tok := <-ch:
			fmt.Println("adding to librarian:", tok.Token)
		case <-done:
			common.Log("Existing indexAdder.")
			return
		}
	}
}

func docProcessor(in chan payload, dStoreCh chan document, dProcessCh chan document, done chan bool) {
	for {
		select {
		case newDoc := <-in:
			var err error
			doc := ""

			if doc, err = getFile(newDoc.URL); err != nil {
				common.Warn(err.Error())
				continue
			}

			titleID := getTitleHash(newDoc.Title)
			msg := document{
				Doc:   doc,
				DocID: titleID,
				Title: newDoc.Title,
			}

			dStoreCh <- msg
			dProcessCh <- msg
		case <-done:
			common.Log("Existing docProcessor.")
			return
		}
	}
}

func indexProcessor(ch chan document, lStoreCh chan lMeta, iAddCh chan token, done chan bool) {
	for {
		select {
		case doc := <-ch:
			docLines := strings.Split(doc.Doc, "\n")

			lin := 0
			for _, line := range docLines {
				if strings.TrimSpace(line) == "" {
					continue
				}

				lStoreCh <- lMeta{
					LIndex: lin,
					Line:   line,
					DocID:  doc.DocID,
				}

				index := 0
				words := strings.Fields(line)
				for _, word := range words {
					if tok, valid := common.SimplifyToken(work); valid {
						iAddCh <- token{
							Token:  tok,
							LIndex: lin,
							Line:   line,
							Index:  index,
							DocID:  doc.DocID,
							Title:  doc.Title,
						}
						index++
					}
				}

				lin++
			}

		case <-done:
			common.Log("Existing indexProcessor.")
			return
		}
	}
}

func docStore(add chan document, get chan dMsg, dGetAllCh chan dAllMsg, done chan bool) {
	store := map[string]document{}

	for {
		select {
		case doc := <-add:
			store[doc.DocID] = doc
		case m := <-get:
			m.Ch <- store[m.DocID]
		case ch := <-dGetAllCh:
			docs := []documents{}
			for _, doc := range store {
				docs = append(docs, doc)
			}
			ch.Ch <- docs
		case <-done:
			common.Log("Existing docStore.")
			return
		}
	}
}

func lineStore(ch chan lMeta, callback chan lMsg, done chan bool) {
	store := map[string]string{}
	for {
		select {
		case line := <-ch:
			id := fmt.Sprintf("%s-%d", line.DocID, line.LIndex)
			store[id] = line.Line
		case ch := <-callback:
			line := ""
			id := fmt.Sprintf("%s-%d", ch.DocID, ch.LIndex)
			if l, exists := store[id]; exists {
				line = l
			}
			ch.Ch <- line
		case <-done:
			common.Log("Existing docStore.")
			return
		}
	}
}

func getTitleHash(title string) string {
	hash := sha1.New()
	title = strings.ToLower(title)

	str := fmt.Sprintf("%s-%s", time.Now(), title)
	hash.Write([]byte(str))

	hByte := hash.Sum(nil)

	return fmt.Sprintf("%x", hByte)
}

func getFile(URL string) (string, error) {
	var res *http.Response
	var err error

	if res, err = http.Get(URL); err != nil {
		errMsg := fmt.Errorf("Unable to retrieve URL: %s.\nError: %s", URL, err)

		return "", errMsg
	}
	if res.StatusCode > 200 {
		errMsg := fmt.Sprintf("Unable to retrieve URL: %s.\nStatus Code: %d")

		return "", errMsg
	}

	body, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()

	if err != nil {
		errMsg := fmt.Errorf("Error while reading response: URL: %s.\nError: %s", URL, res)

		return "", errMsg
	}

	return string(body), nil
}

func FeedHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		ch := make(chan []document)
		dGetAllCh <- dAllMsg{Ch: ch}
		docs := <-ch
		close(ch)

		if serializedPayload, err := json.Marshal(docs); err != nil {
			w.Write(serializedPayload)
		} else {
			common.Warn("Unable to serialize all docs: " + err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"code": 500, "msg": "Error occurred while trying to receive documents."}`))
		}
		return
	} else if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"code": 405, "msg": "Method Not Allowed."}`))
		return
	}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var newDoc payload
	decoder.Decode(&newDoc)
	pProcessCh <- newDoc

	w.Write([]byte(`{"code": 200, "msg": "Request is being processed."}`))
}
