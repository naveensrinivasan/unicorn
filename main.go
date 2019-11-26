package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/99designs/keyring"
	"github.com/atotto/clipboard"
	"github.com/jinzhu/configor"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type configuration struct {
	Domain      string
	UserName    string
	EmailDomain string
}
type handle []struct {
	Aliases []struct {
		Address          string      `json:"address"`
		AddressDisplay   string      `json:"address_display"`
		ForwardsTo       []string    `json:"forwards_to"`
		PermittedSenders interface{} `json:"permitted_senders"`
		Required         bool        `json:"required"`
	} `json:"aliases"`
	Domain string `json:"domain"`
}

func main() {
	config := configuration{}
	e := configor.Load(&config, "settings.yaml")
	if e != nil {
		log.Fatal(e)
	}
	ring, _ := keyring.Open(keyring.Config{})
	password, e := ring.Get(config.UserName)
	if e != nil {
		log.Fatal(e)
	}
	aliases := getEmailAliases(config, string(password.Data))
	result, e := generateRandomEmail(aliases, config, string(password.Data))
	if e != nil {
		log.Fatal(e)
	}
	_ = clipboard.WriteAll(result)
	log.Println(result)
}

func generateRandomEmail(aliases map[string]interface{}, config configuration, password string) (string, error) {

	for i := 0; i < 5; i++ {
		uuid := pseudoUUID()
		if _, ok := aliases[uuid]; !ok {
			log.Println("Can generate a new email address", uuid)
			client := &http.Client{}
			body := strings.NewReader(fmt.Sprintf(`address=%s@%s&forwards_to=%s`,
				uuid,
				config.EmailDomain,
				config.UserName))
			req, err := http.NewRequest("POST",
				fmt.Sprintf("https://%s/admin/mail/aliases/add",
					config.Domain),
				body)
			fmt.Println(password)
			req.SetBasicAuth(config.UserName, password)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}
			bodyText, err := ioutil.ReadAll(resp.Body)
			log.Printf("%s\n", string(bodyText))
			return fmt.Sprintf("%s@%s", uuid, config.EmailDomain), nil
		}
	}
	return "", errors.New("could not generate random email address")
}

func getEmailAliases(config configuration, password string) map[string]interface{} {
	url := fmt.Sprintf("https://%s/admin/mail/aliases?format=json", config.Domain)
	spaceClient := http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.SetBasicAuth(config.UserName, password)
	res, getErr := spaceClient.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}
	handles := handle{}
	jsonErr := json.Unmarshal(body, &handles)
	if jsonErr != nil {
		log.Fatal(jsonErr)
	}
	aliases := make(map[string]interface{})
	for _, y := range handles {
		for _, z := range y.Aliases {
			aliases[z.Address] = z
		}
	}
	return aliases
}

func pseudoUUID() string {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Fatal(err)
	}
	return fmt.Sprintf("%X%X%X%X%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
