package twitter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//1865462090661515756

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

type RetweetLookup struct {
	Data []User `json:"data"`
}

type TwitterClient struct {
	BearerToken string
}

func (c TwitterClient) GetRetweetedFromPostId(postId string) ([]string, error) {
	url := fmt.Sprintf("https://api.x.com/2/tweets/%s/retweeted_by", postId)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", c.BearerToken))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", body)

	var data RetweetLookup
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	fmt.Printf("%v\n", data)

	users := make([]string, 0, len(data.Data))

	for _, user := range data.Data {
		users = append(users, user.Username)
	}

	return users, nil
}
