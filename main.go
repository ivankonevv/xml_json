package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Items struct {
	XMLName xml.Name `xml:"items" json:"-"`
	Items   []Item   `xml:"item" json:"items"`
}

type Item struct {
	XMLName     xml.Name `xml:"item" json:"-"`
	URL         string   `xml:"url" json:"url"`
	PhoneNumber string   `xml:"phone_number" json:"phone_number"`
	Location    string   `xml:"location" json:"location"`
	Menu        []Menu   `xml:"menu" json:"menu"`
}

type Menu struct {
	XMLName xml.Name `xml:"menu" json:"-"`
	Value   []Value  `xml:"value" json:"value"`
}

type Value struct {
	XMLName  xml.Name `xml:"value" json:"-"`
	MenuName string   `xml:"menu_name" json:"menu_name"`
	Price    int      `xml:"price" json:"price"`
}

func (items Items) ParseXML(filename string) (*Items, error) {
	xmlFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open `%s` file: %s\n", filename, err)
	}

	defer func(xmlFile *os.File) {
		err := xmlFile.Close()
		if err != nil {
			fmt.Printf("failed to close `%s`: %s\n", filename, err)
		}
	}(xmlFile)

	byteValue, _ := ioutil.ReadAll(xmlFile)

	if err = xml.Unmarshal(byteValue, &items); err != nil {
		return nil, fmt.Errorf("error unmarshalling `%s`: %s\n", filename, err)
	}

	for i := range items.Items {
		fmt.Println("----------------------")
		fmt.Println("URL: " + items.Items[i].URL)
		fmt.Println("Location: " + items.Items[i].Location)
		fmt.Println("Phone: " + items.Items[i].PhoneNumber)

		for n := range items.Items[i].Menu {
			for x := range items.Items[i].Menu[n].Value {
				fmt.Println("- Menu: " + items.Items[i].Menu[n].Value[x].MenuName)
				fmt.Println("-- Price: " + strconv.Itoa(items.Items[i].Menu[n].Value[x].Price))
			}
		}
	}
	return &items, nil
}

func (items Items) CreateJSON(filename string) error {
	b, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("error marshalling items: %s\n", err)
	}
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return fmt.Errorf("error writing file `%s`: %s\n", filename, err)
	}

	return nil
}

func main() {
	var items Items
	err := godotenv.Load(".env.dev")
	if err != nil {
		fmt.Printf("error loading .env file")
		return
	}

	item, err := items.ParseXML(os.Getenv("XML_FILE_NAME"))
	if err != nil {
		fmt.Printf("an error occurred in ParseXML(): %s\n", err)
		return
	}

	if err := item.CreateJSON(os.Getenv("JSON_FILE_NAME")); err != nil {
		fmt.Printf("an error occurred in CreateJSON(): %s\n", err)
		return
	}
}
