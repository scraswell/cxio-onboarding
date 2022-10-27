package onboarding

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/customerio/go-customerio/v3"
	"github.com/go-errors/errors"
	"gopkg.in/yaml.v3"
)

type Config struct {
	csvPath     string `yaml:"csvPath"`
	siteID      string `yaml:"siteID"`
	trackAPIKey string `yaml:"trackAPIKey"`
}

type Person struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Title string `json:"title"`
	Dept  string `json:"department"`
	Since int    `json:"created_at"`
}

type FieldDescriptor struct {
	Name      string
	FieldType reflect.Kind
}

func getAttr(obj interface{}, fieldName string) reflect.Value {
	pointToStruct := reflect.ValueOf(obj)

	curStruct := pointToStruct.Elem()
	if curStruct.Kind() != reflect.Struct {
		panic("not struct")
	}

	curField := curStruct.FieldByName(fieldName)
	if !curField.IsValid() {
		panic("not found:" + fieldName)
	}
	return curField
}

func createFieldMap() map[int]FieldDescriptor {
	fieldMap := make(map[int]FieldDescriptor)

	idField := FieldDescriptor{
		Name:      "Id",
		FieldType: reflect.Int64,
	}

	nameField := FieldDescriptor{
		Name:      "Name",
		FieldType: reflect.String,
	}

	emailField := FieldDescriptor{
		Name:      "Email",
		FieldType: reflect.String,
	}

	titleField := FieldDescriptor{
		Name:      "Title",
		FieldType: reflect.String,
	}

	deptField := FieldDescriptor{
		Name:      "Dept",
		FieldType: reflect.String,
	}

	sinceField := FieldDescriptor{
		Name:      "Since",
		FieldType: reflect.Int64,
	}

	fieldMap[0] = idField
	fieldMap[1] = nameField
	fieldMap[2] = emailField
	fieldMap[3] = titleField
	fieldMap[4] = deptField
	fieldMap[5] = sinceField

	return fieldMap
}

func createPersonList(data [][]string) []Person {
	var fieldMap = createFieldMap()
	var personList []Person

	for i, line := range data {
		if i > 0 {
			var p Person
			for j, field := range line {
				fieldName := fieldMap[j].Name
				fieldType := fieldMap[j].FieldType

				if fieldType == reflect.Int64 {
					value, err := strconv.ParseInt(field, 10, 64)
					if err == nil {
						getAttr(&p, fieldName).SetInt(value)
					}
				} else if fieldType == reflect.String {
					getAttr(&p, fieldName).SetString(field)
				}
			}

			personList = append(personList, p)
		}
	}

	return personList
}

func GetPersonListFromCsv(csvPath string) []Person {
	f, err := os.Open(csvPath)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	csvReader := csv.NewReader(f)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	return createPersonList(data)
}

func readConf(filename string) (*Config, error) {
	buf, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// I tried to unmarshall this directly to the struct; but it failed silently... :(
	data := make(map[string]string)
	err = yaml.Unmarshal(buf, data)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %w", filename, err)
	}

	c := Config{
		csvPath:     data["csvPath"],
		trackAPIKey: data["trackAPIKey"],
		siteID:      data["siteID"],
	}

	return &c, err
}

func LoadPeople(configPath string) {
	log.Printf("Loading configuration from: %s...", configPath)

	cfg, err := readConf(configPath)
	if err != nil {
		log.Printf("here")
		log.Fatal(fmt.Println(err.(*errors.Error).ErrorStack()))
	}

	log.Printf("Configuration loaded from: %s.", configPath)
	log.Printf("\tCSV path is: %s", cfg.csvPath)

	track := customerio.NewTrackClient(
		cfg.siteID,
		cfg.trackAPIKey,
		customerio.WithRegion(customerio.RegionUS))

	personList := GetPersonListFromCsv(cfg.csvPath)

	for _, person := range personList {
		log.Printf("Invoking Identify for user %d...", person.Id)
		if err := track.Identify(fmt.Sprint(person.Id), map[string]interface{}{
			"email":      person.Email,
			"created_at": time.Now().Unix(),
			"first_name": person.Name,
			"title":      person.Title,
			"department": person.Dept,
		}); err != nil {
			log.Fatal(err)
		}
	}
}
