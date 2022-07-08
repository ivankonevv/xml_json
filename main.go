package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"

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

func (items *Items) ParseXML(filename string) error {

	xmlFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open `%s` file: %v", filename, err)
	}

	byteValue, _ := ioutil.ReadAll(xmlFile)

	if err = xml.Unmarshal(byteValue, &items); err != nil {
		return fmt.Errorf("error unmarshalling `%s`: %v\n", filename, err)
	}

	for i := range items.Items {
		fmt.Println("----------------------")
		fmt.Println("URL: " + items.Items[i].URL)
		fmt.Println("Location: " + items.Items[i].Location)
		fmt.Println("Phone: " + items.Items[i].PhoneNumber)

		for n := range items.Items[i].Menu {
			for x := range items.Items[i].Menu[n].Value {
				fmt.Println("- Menu: " + items.Items[i].Menu[n].Value[x].MenuName)
				fmt.Println("-- Price: ", strconv.Itoa(items.Items[i].Menu[n].Value[x].Price))
			}
		}
	}
	err = xmlFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close `%s`: %v", filename, err)
	}

	return nil
}

func (items Items) CreateJSON(filename string) error {
	b, err := json.Marshal(items)
	if err != nil {
		return fmt.Errorf("error marshalling items: %v\n", err)
	}
	err = ioutil.WriteFile(filename, b, 0644)
	if err != nil {
		return fmt.Errorf("error writing file `%s`: %v\n", filename, err)
	}

	return nil
}

func main() {
	var items Items

	// Logger init
	XMLLogger := log.New().WithField("method", "ParseXML()")
	CreateJSONLogger := log.New().WithField("method", "CreateJSON()")
	MainLogger := log.New().WithField("function", "main()")

	// Load env variables
	if err := godotenv.Load(".env.dev"); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			MainLogger.Error("init env variables error occurred: ", err)
			return
		}

		MainLogger.Error("os.ErrNotExist occurred while loading env variables: ", err)
		return
	}
	MainLogger.Info("env file opened successfully")

	// Parse xml file
	if err := items.ParseXML(os.Getenv("XML_FILE_NAME")); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			XMLLogger.Error("an error occurred: ", err)
			return
		}

		XMLLogger.Error("os.ErrNotExist occurred: ", err)
		return
	}
	XMLLogger.Info("xml parsing process finished successfully.")

	// Write to JSON
	if err := items.CreateJSON(os.Getenv("JSON_FILE_NAME")); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			CreateJSONLogger.Error("an error occurred.")
			return
		}

		CreateJSONLogger.Error("os.ErrNotExist occurred: ", err)
		return
	}
	CreateJSONLogger.Info("writing to JSON finished successfully.")
}
