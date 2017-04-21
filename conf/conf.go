package conf

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"syscall"

	"github.com/umayr/lunch/endpoint"
	"golang.org/x/crypto/ssh/terminal"
)

type Conf struct {
	Token string
}

func homedir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func credentials() (string, string, error) {
	defer fmt.Println()
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter Username: ")
	username, err := reader.ReadString('\n')
	if err != nil {
		return "", "", err
	}

	fmt.Print("Enter Password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", "", err
	}

	password := string(bytePassword)
	return strings.TrimSpace(username), strings.TrimSpace(password), nil
}

func Update() (*Conf, error) {
	username, password, err := credentials()
	if err != nil {
		return nil, err
	}

	raw, err := endpoint.Auth(username, password)
	if err != nil {
		return nil, err
	}

	var response struct {
		Data struct {
			Attr struct {
				Token string `json:"access-token"`
			} `json:"attributes"`
		}
	}

	if err := json.Unmarshal(raw, &response); err != nil {
		return nil, err
	}

	cur, err := ioutil.ReadFile(fmt.Sprintf("%s/.lunchrc", homedir()))
	if err != nil {
		return nil, err
	}

	c := Conf{}
	if err := json.Unmarshal(cur, &c); err != nil {
		return nil, err
	}

	c.Token = response.Data.Attr.Token

	buf, err := json.Marshal(c)
	if err != nil {
		return nil, err
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/.lunchrc", homedir()), buf, 0644); err != nil {
		return nil, err
	}

	return &c, nil
}

func Get() (*Conf, error) {
	if _, err := os.Stat(fmt.Sprintf("%s/.lunchrc", homedir())); os.IsNotExist(err) {
		username, password, err := credentials()
		if err != nil {
			return nil, err
		}

		raw, err := endpoint.Auth(username, password)
		if err != nil {
			return nil, err
		}

		var response struct {
			Data struct {
				Attr struct {
					Token string `json:"access-token"`
				} `json:"attributes"`
			}
		}

		if err := json.Unmarshal(raw, &response); err != nil {
			return nil, err
		}

		c := Conf{response.Data.Attr.Token}
		buf, err := json.Marshal(c)
		if err != nil {
			return nil, err
		}

		if err := ioutil.WriteFile(fmt.Sprintf("%s/.lunchrc", homedir()), buf, 0644); err != nil {
			return nil, err
		}

		return &c, nil
	}

	raw, err := ioutil.ReadFile(fmt.Sprintf("%s/.lunchrc", homedir()))
	if err != nil {
		return nil, err
	}

	c := Conf{}
	if err := json.Unmarshal(raw, &c); err != nil {
		return nil, err
	}

	return &c, nil
}
