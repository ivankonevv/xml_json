package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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

func ParseXML() Items {
	xmlFile, err := os.Open("test.xml")
	if err != nil {
		fmt.Println("error opening file.")
	}
	fmt.Println("Opened...")

	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)

	var items Items

	if err = xml.Unmarshal(byteValue, &items); err != nil {
		fmt.Println("error unmarshalling file.")
	}

	for i := 0; i < len(items.Items); i++ {
		fmt.Println("----------------------")
		fmt.Println("URL: " + items.Items[i].URL)
		fmt.Println("Location: " + items.Items[i].Location)
		fmt.Println("Phone: " + items.Items[i].PhoneNumber)
		for n := 0; n < len(items.Items[i].Menu); n++ {
			for x := 0; x < len(items.Items[i].Menu[n].Value); x++ {
				fmt.Println("- Menu: " + items.Items[i].Menu[n].Value[x].MenuName)
				fmt.Println("-- Price: " + strconv.Itoa(items.Items[i].Menu[n].Value[x].Price))
			}
		}
	}

	return items
}

func CreateJSON(items Items) {
	b, err := json.Marshal(items)
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = ioutil.WriteFile("test.json", b, 0644)

}

func main() {
	items := ParseXML()
	CreateJSON(items)
}