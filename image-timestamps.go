package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const dockerHubAPI = "https://hub.docker.com/v2/repositories"

type TagList struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []Tag  `json:"results"`
}

type Tag struct {
	Name        string `json:"name"`
	LastUpdated string `json:"last_updated"`
	Images      []struct {
		Digest string `json:"digest"`
	} `json:"images"`
}

func getTags(namespace, repository string) ([]Tag, error) {
	var tags []Tag
	url := fmt.Sprintf("%s/%s/%s/tags/", dockerHubAPI, namespace, repository)

	for {
		res, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		var tagList TagList
		if err := json.Unmarshal(body, &tagList); err != nil {
			return nil, err
		}

		tags = append(tags, tagList.Results...)

		if tagList.Next == "" {
			break
		}
		url = tagList.Next
	}

	return tags, nil
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <namespace> <repository>", os.Args[0])
	}
	namespace := os.Args[1]
	repository := os.Args[2]

	tags, err := getTags(namespace, repository)
	if err != nil {
		log.Fatal(err)
	}

	for _, tag := range tags {
		fmt.Printf("Tag: %s, Last Updated: %s, Digest: %s\n", tag.Name, tag.LastUpdated, tag.Images[0].Digest)
	}
}
