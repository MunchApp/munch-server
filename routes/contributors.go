package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ContributorResponse struct {
	login             string
	id                int64
	nodeid            string
	avatarurl         string
	gravatarid        string
	url               string
	htmlurl           string
	followersurl      string
	followingurl      string
	gistsurl          string
	starredurl        string
	subscriptionsurl  string
	organizationsurl  string
	reposurl          string
	eventsurl         string
	receivedeventsurl string
	typeresponse      string
	siteadmin         bool
	contributions     int
}

// [
// 	{
// 		"login": "kftang",
// 		"id": 15274333,
// 		"node_id": "MDQ6VXNlcjE1Mjc0MzMz",
// 		"avatar_url": "https://avatars0.githubusercontent.com/u/15274333?v=4",
// 		"gravatar_id": "",
// 		"url": "https://api.github.com/users/kftang",
// 		"html_url": "https://github.com/kftang",
// 		"followers_url": "https://api.github.com/users/kftang/followers",
// 		"following_url": "https://api.github.com/users/kftang/following{/other_user}",
// 		"gists_url": "https://api.github.com/users/kftang/gists{/gist_id}",
// 		"starred_url": "https://api.github.com/users/kftang/starred{/owner}{/repo}",
// 		"subscriptions_url": "https://api.github.com/users/kftang/subscriptions",
// 		"organizations_url": "https://api.github.com/users/kftang/orgs",
// 		"repos_url": "https://api.github.com/users/kftang/repos",
// 		"events_url": "https://api.github.com/users/kftang/events{/privacy}",
// 		"received_events_url": "https://api.github.com/users/kftang/received_events",
// 		"type": "User",
// 		"site_admin": false,
// 		"contributions": 11
// 	  },
// 	  {
// 		"login": "Lmnorrell99",
// 		"id": 31517170,
// 		"node_id": "MDQ6VXNlcjMxNTE3MTcw",
// 		"avatar_url": "https://avatars2.githubusercontent.com/u/31517170?v=4",
// 		"gravatar_id": "",
// 		"url": "https://api.github.com/users/Lmnorrell99",
// 		"html_url": "https://github.com/Lmnorrell99",
// 		"followers_url": "https://api.github.com/users/Lmnorrell99/followers",
// 		"following_url": "https://api.github.com/users/Lmnorrell99/following{/other_user}",
// 		"gists_url": "https://api.github.com/users/Lmnorrell99/gists{/gist_id}",
// 		"starred_url": "https://api.github.com/users/Lmnorrell99/starred{/owner}{/repo}",
// 		"subscriptions_url": "https://api.github.com/users/Lmnorrell99/subscriptions",
// 		"organizations_url": "https://api.github.com/users/Lmnorrell99/orgs",
// 		"repos_url": "https://api.github.com/users/Lmnorrell99/repos",
// 		"events_url": "https://api.github.com/users/Lmnorrell99/events{/privacy}",
// 		"received_events_url": "https://api.github.com/users/Lmnorrell99/received_events",
// 		"type": "User",
// 		"site_admin": false,
// 		"contributions": 1
// 	  }
// 	]

func GetContributorsHandler(w http.ResponseWriter, r *http.Request) {

	//get users from the HTTP link
	resp, err := http.Get("https://api.github.com/repos/MunchApp/munchserver/contributors")
	if err != nil {
		//handle error
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	englishBody := string(body)
	fmt.Println(englishBody)

	type ContributorResponseTwo struct {
		Username string `json:"Login"`
		// Id                 int
		// Node_id            string
		// Avatar_url         string
		// Gravatar_id        string
		// Url                string
		// Html_url           string
		// Followers_url      string
		// Following_url      string
		// Gists_url          string
		// Starred_url        string
		// Subscriptions_url  string
		// Organizations_url  string
		// Repos_url          string
		// Events_url         string
		// Receivedevents_url string
		// Type_response      string 	`json:"type"`
		// Site_admin         bool
		Contributions int `json:"contributions`
	}

	var contributorresponses []ContributorResponseTwo
	jsonErr := json.Unmarshal(body, &contributorresponses)
	if jsonErr != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("ContributorResponses : %+v", contributorresponses)

}

func GetContributors2Handler(w http.ResponseWriter, r *http.Request) {
	var jsonBlob = []byte(`[
	{"Name": "Platypus", "Order": "Monotremata"},
	{"Name": "Quoll",    "Order": "Dasyuromorphia"}
]`)
	type Animal struct {
		Name  string
		Order string
	}
	var animals []Animal
	err := json.Unmarshal(jsonBlob, &animals)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", animals)
}
