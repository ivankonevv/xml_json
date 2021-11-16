package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
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

func (items *Items) ParseXML(filename string) error {

	defer runRecover()
	xmlFile, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open `%s` file: %v", filename, err)
	}

	err = xmlFile.Close()
	if err != nil {
		return fmt.Errorf("failed to close `%s`: %v", filename, err)
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
				fmt.Println("-- Price: " + strconv.Itoa(items.Items[i].Menu[n].Value[x].Price))
			}
		}
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

func HandleErrors(err error, errText string) {
	defer runRecover()
	if errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s:\n%v", errText, err)
		return
	}
	fmt.Printf("%v. Recovering.\n", err)
	panic(err)

}

func runRecover() {
	if r := recover(); r != nil {
		fmt.Println("Recovered:", r)
	}
}

func main() {
	var items Items
	// Load env variables
	HandleErrors(godotenv.Load(".env.dev"), "error loading env variables")

	// Parse xml file
	HandleErrors(items.ParseXML(os.Getenv("XML_FILE_NAME")), "an error occurred in ParseXML()")

	// Write to JSON
	HandleErrors(items.CreateJSON(os.Getenv("JSON_FILE_NAME")), "an error occurred in CreateJSON()")

}
