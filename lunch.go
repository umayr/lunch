package lunch

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/umayr/lunch/conf"
	"github.com/umayr/lunch/endpoint"
)

var c *conf.Conf

func init() {
	var err error

	c, err = conf.Get()
	if err != nil {
		panic(err)
	}
}

type response struct {
	Data []struct {
		ID   int `json:"id"`
		Attr struct {
			Name     string `json:"menu-item"`
			Date     Date   `json:"lunch-date"`
			Likes    int    `json:"likes-count"`
			Dislikes int    `json:"dislikes-count"`
		} `json:"attributes"`
	} `json:"data"`
}

type config struct {
	Token string
}

type Item struct {
	ID       int
	Name     string
	Date     time.Time
	Likes    int
	Dislikes int
}

type Lunches map[string]Item

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(b []byte) error {
	str := strings.TrimSpace(string(b))
	str = str[1 : len(str)-1]
	if str == "" {
		d.Time = time.Time{}
		return nil
	}

	t, err := time.Parse("2006-01-02T15:04:05", str)
	if err != nil {
		return err
	}

	d.Time = t
	return nil
}

func parse(raw []byte) (Lunches, error) {
	var r response
	if err := json.Unmarshal(raw, &r); err != nil {
		return nil, err
	}

	l := make(Lunches, len(r.Data))
	for _, d := range r.Data {
		l[d.Attr.Date.Format("2006-01-02")] = Item{d.ID, d.Attr.Name, d.Attr.Date.Time, d.Attr.Likes, d.Attr.Dislikes}
	}

	return l, nil
}

func Find() (Lunches, error) {
	raw, err := endpoint.Lunch(c.Token)
	if err != nil {
		switch e := err.(type) {
		case endpoint.HttpError:
			if e.Code != 401 {
				return nil, err
			}

			c, err = conf.Update()
			if err != nil {
				return nil, err
			}

			raw, err = endpoint.Lunch(c.Token)
			if err != nil {
				return nil, err
			}
		default:
			return nil, err
		}
	}
	return parse(raw)
}

func Today() (Item, error) {
	l, err := Find()
	if err != nil {
		return Item{}, err
	}

	return l[time.Now().Format("2006-01-02")], nil
}
