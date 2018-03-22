package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type Link struct {
	LinkID      int `json:"LinkID"`
	LinkContent int `json:"LinkContent"`
}

type LinkEngOrJap struct {
	linkID int `json:"linkID"`
	eng    int `json:"english"`
	jap    int `json:"japanese"`
}

type Sentence struct {
	SentenceID      int    `json:"SentenceID"`
	Lang            string `json:"Lang"`
	SentenceContent string `json:"SentenceContent"`
}

type SentenceWithAudio struct {
	SentenceID      int    `json:"SentenceID"`
	SentenceContent string `json:"SentenceContent"`
}

type Result struct {
	Eng string `json:"english"`
	Jap string `json:"japanese"`
	Auth     string `json:"author"`
	SoundID  int `json:"soundId"`
	Syllabes int `json:"syllabes"`
}

func getLink(filePath string) []Link {
	var links []Link

	// fmt.Println(filePath)
	csvFile, _ := os.Open(filePath)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.LazyQuotes = true

	// i := 0
	for {
		// i += 1
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		// fmt.Println(line)
		id, err := strconv.Atoi(line[0])
		if err == nil {
			content, err := strconv.Atoi(line[1])
			if err == nil {
				link := Link{
					LinkID:      id,
					LinkContent: content,
				}
				// fmt.Println(link)
				links = append(links, link)
			} else {
				fmt.Println(err)
			}
		} else {
			fmt.Println(err)
		}
		// break
		// if i == 1000 {
		// 	break
		// }
	}
	// fmt.Println(links)
	return links
}

func getSentence(filePath string) []Sentence {
	var sentences []Sentence

	// fmt.Println(filePath)
	csvFile, _ := os.Open(filePath)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.LazyQuotes = true
	// i := 0
	for {
		// i += 1
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		// fmt.Println(line)
		id, err := strconv.Atoi(line[0])
		if err == nil {
			sentence := Sentence{
				SentenceID:      id,
				Lang:            line[1],
				SentenceContent: line[2],
			}
			sentences = append(sentences, sentence)
		} else {
			fmt.Println(err)
		}
		// if i == 1000 {
		// 	break
		// }
	}
	// fmt.Println(sentences)
	return sentences
}

func getSentenceWithAudio(filePath string) []SentenceWithAudio {
	var sentences []SentenceWithAudio

	// fmt.Println(filePath)
	csvFile, _ := os.Open(filePath)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'
	reader.LazyQuotes = true

	// i := 0
	for {
		// i += 1
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		// fmt.Println(line)
		id, err := strconv.Atoi(line[0])
		if err == nil {
			sentence := SentenceWithAudio{
				SentenceID:      id,
				SentenceContent: line[3],
			}
			sentences = append(sentences, sentence)
		} else {
			fmt.Println(err)
		}
		// break
		// if i == 2 {
		// 	break
		// }
	}
	// fmt.Println(sentences)
	return sentences
}

func getEngAndJapSentence(sentences []Sentence) ([]Sentence, []Sentence) {
	var engSentences []Sentence
	var japSentences []Sentence

	for _, value := range sentences {
		if value.Lang == "eng" {
			engSentences = append(engSentences, value)
		} else if value.Lang == "jpn" {
			japSentences = append(japSentences, value)
		}
	}
	// fmt.Println(engSentences)
	// fmt.Println(japSentences)
	return engSentences, japSentences
}

func findASentenceByID(sentences []Sentence, id int) bool {
	// fmt.Println("find ", id, " on ", sentences)
	for _, value := range sentences {
		if value.SentenceID == id {
			return true
		}
	}
	return false
}

func getValueSentenceContentWithID(sentences []Sentence, sentenceID int) string {
	for _, value := range sentences {
		if value.SentenceID == sentenceID {
			return value.SentenceContent
		}
	}
	return ""
}

func findAudioSentenceByID(sentences []SentenceWithAudio, id int) bool {
	// fmt.Println("find ", id, " on ", sentences)
	for _, value := range sentences {
		if value.SentenceID == id {
			if value.SentenceContent != "" {
				return true
			}
		}
	}
	return false
}

func getEngAndJapLink(links []Link, engSentences []Sentence, japSentences []Sentence) []LinkEngOrJap {
	var engAndJapLink []LinkEngOrJap
	var langList []int
	var engID = -1
	var japID = -1
	currentID := 1

	for ID, value := range links {
		if currentID == value.LinkID { // get all of elements have same id
			langList = append(langList, value.LinkContent)
		} else { // process current list
			if len(langList) > 0 {
				for _, id := range langList {
					// fmt.Println(id, ": \n", engSentences)
					if findASentenceByID(engSentences, id) {
						engID = id
					}
					if findASentenceByID(japSentences, id) {
						japID = id
					}
				}
				if engID != -1 || japID != -1 {
					engAndJapLink = append(engAndJapLink, LinkEngOrJap{
						linkID: currentID,
						eng:    engID,
						jap:    japID,
					})
				}

			}
			// reset for next
			langList = nil
			engID = -1
			japID = -1
			currentID = ID
			langList = append(langList, value.LinkContent)
			// break
		}
	}

	return engAndJapLink
}

func buidResult(engAndJapLink []LinkEngOrJap, eng []Sentence, jap []Sentence, audioSentences []SentenceWithAudio) []Result {
	var result []Result

	for _, value := range engAndJapLink {
		var aResult Result
		if value.eng == -1 {	// has no eng
			continue
		}
		aResult.Eng = getValueSentenceContentWithID(eng, value.eng)
		if value.jap == -1 {	// has no japanese
			continue
		}
		aResult.Jap = getValueSentenceContentWithID(jap, value.jap)
		if !findAudioSentenceByID(audioSentences, value.eng) {	// has no audio
			continue
		}
		aResult.SoundID = value.eng
		aResult.Auth = "hai nguyen"
		aResult.Syllabes = strings.Count(aResult.Eng, " ") + 1
		result = append(result, aResult)
	}

	return result
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "\t")
	if err != nil {
		return in
	}
	return out.String()
}

func writeResult(resultPath string, content string) {
	f, err := os.Create(resultPath)
	if err == nil {
		defer f.Close()
		f.WriteString(content)
		f.Sync()
	}
}

func main() {
	basePath := "/Users/hainguyen/Desktop/CrewzTest"
	links := getLink(basePath + "/links.csv")
	sentences := getSentence(basePath + "/sentences.csv")
	sentencesWithAudio := getSentenceWithAudio(basePath + "/sentences_with_audio.csv")
	engSentences, japSentences := getEngAndJapSentence(sentences)
	engAndJapLink := getEngAndJapLink(links, engSentences, japSentences)
	result := buidResult(engAndJapLink, engSentences, japSentences, sentencesWithAudio)
	// fmt.Println(result)
	// json.Marshal(result)
	// fmt.Println(string(resultJson))

	resultJson, _ := json.Marshal(result)
	fmt.Println(jsonPrettyPrint(string(resultJson)))
	// writeResult(basePath + "/result.json", jsonPrettyPrint(string(resultJson)))
}
