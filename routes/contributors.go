package routes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ContributorResponse struct {
	Username      string `json:"Login"`
	Contributions int    `json:"contributions`
}

type ReturnResponse struct {
	Username      string `json: "login"`
	Contributions int    `json:"contributions"`
	Issues        int
}

type IssueResponse struct {
	User struct {
		Login string `json: "login"`
	} `json: "user"`
	IssueNumber int `json: "number"`
}

func GetContributorsHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println()
	fmt.Println("---------------------------------")
	fmt.Println("creating contributors response...")
	fmt.Println("---------------------------------")
	fmt.Println()

	//TODO: cache maybe?

	///////////////////////////////////////
	// GETTING CONTRIBUTIONS FROM SERVER //
	///////////////////////////////////////

	fmt.Println("getting munchserver contributions...")

	//get users from the HTTP link
	resp, err := http.Get("https://api.github.com/repos/MunchApp/munchserver/contributors")
	if err != nil {
		fmt.Println("error:", err)
	}

	//close the response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	//Create an array and print the contents of the array
	var contributorResponsesServer []ContributorResponse
	jsonErr := json.Unmarshal(body, &contributorResponsesServer)
	if jsonErr != nil {
		fmt.Println("error:", err)
	}

	//////////////////////////////////////////
	// GETTING CONTRIBUTIONS FROM munch-app //
	//////////////////////////////////////////

	fmt.Println("getting munch-app contributions...")

	//get users from the HTTP link
	respApp, errApp := http.Get("https://api.github.com/repos/MunchApp/munch-app/contributors")
	if errApp != nil {
		fmt.Println("error:", err)
	}

	//close the response
	defer respApp.Body.Close()
	bodyApp, errApp := ioutil.ReadAll(respApp.Body)

	//Create an array and print the contents of the array
	var contributorResponsesApp []ContributorResponse
	jsonErrApp := json.Unmarshal(bodyApp, &contributorResponsesApp)
	if jsonErrApp != nil {
		fmt.Println("error:", err)
	}

	//////////////////////////////////////
	// GETTING ALL OF THE CLOSED ISSUES //
	//////////////////////////////////////

	fmt.Println("getting munch-app closed issues...")

	respClosed, errClosed := http.Get("https://api.github.com/repos/MunchApp/munch-app/issues?state=closed")
	if errClosed != nil {
		fmt.Println("Error getting Github Closed Issues: ", errClosed)
	}

	defer respClosed.Body.Close()
	bodyClosed, errClosed := ioutil.ReadAll(respClosed.Body)

	var contributorAppClosedIssues []IssueResponse
	closedIssueErr := json.Unmarshal(bodyClosed, &contributorAppClosedIssues)
	if closedIssueErr != nil {
		fmt.Println("Error creating closed issue array...")
	}

	fmt.Println("getting munchserver closed issues...")

	respClosedServer, errClosed := http.Get("https://api.github.com/repos/MunchApp/munchserver/issues?state=closed")
	if errClosed != nil {
		fmt.Println("Error getting Github Closed Issues: ", errClosed)
	}

	defer respClosed.Body.Close()
	bodyClosedServer, errClosed := ioutil.ReadAll(respClosedServer.Body)

	var contributorServerClosedIssues []IssueResponse
	closedIssueErrServer := json.Unmarshal(bodyClosedServer, &contributorServerClosedIssues)
	if closedIssueErrServer != nil {
		fmt.Println("Error creating closed issue array...")
	}

	/////////////////////////////////////
	// GETTING THE LIST OF OPEN ISSUES //
	/////////////////////////////////////

	fmt.Println("getting munch-app open issues...")

	respOpen, errOpen := http.Get("https://api.github.com/repos/MunchApp/munch-app/issues")
	if errOpen != nil {
		fmt.Println("Error getting Open Issues HTTP Request...")
	}

	defer respOpen.Body.Close()
	bodyOpen, errOpen := ioutil.ReadAll(respOpen.Body)

	var contributorAppOpenIssues []IssueResponse
	openIssueErr := json.Unmarshal(bodyOpen, &contributorAppOpenIssues)
	if openIssueErr != nil {
		fmt.Println("Error during unmarshal Open issues JSON...")
	}

	fmt.Println("getting munchserver open issues...")

	respOpenServer, errOpenServer := http.Get("https://api.github.com/repos/MunchApp/munchserver/issues")
	if errOpenServer != nil {
		fmt.Println("Error getting Open Issues HTTP Request...")
	}

	defer respOpen.Body.Close()
	bodyOpenServer, errOpen := ioutil.ReadAll(respOpenServer.Body)

	var contributorServerOpenIssues []IssueResponse
	openIssueErrServer := json.Unmarshal(bodyOpenServer, &contributorServerOpenIssues)
	if openIssueErrServer != nil {
		fmt.Println("Error during unmarshal Open issues JSON...")
	}

	//////////////////////////////////
	// CREATING THE RETURN RESPONSE //
	//////////////////////////////////

	yasira := newReturnResponse("yasirayounus", contributorResponsesServer, contributorResponsesApp)
	yasira.Issues = newIssueCount("yasirayounus", contributorAppClosedIssues, contributorAppOpenIssues)
	yasira.Issues += newIssueCount("yasirayounus", contributorServerClosedIssues, contributorServerOpenIssues)

	kenny := newReturnResponse("kftang", contributorResponsesServer, contributorResponsesApp)
	kenny.Issues = newIssueCount("kftang", contributorAppClosedIssues, contributorAppOpenIssues)
	kenny.Issues += newIssueCount("kftang", contributorServerClosedIssues, contributorServerOpenIssues)

	luke := newReturnResponse("Lmnorrell99", contributorResponsesServer, contributorResponsesApp)
	luke.Issues = newIssueCount("Lmnorrell99", contributorAppClosedIssues, contributorAppOpenIssues)
	luke.Issues += newIssueCount("Lmnorrell99", contributorServerClosedIssues, contributorServerOpenIssues)

	janine := newReturnResponse("janinebar", contributorResponsesServer, contributorResponsesApp)
	janine.Issues = newIssueCount("janinebar", contributorAppClosedIssues, contributorAppOpenIssues)
	janine.Issues += newIssueCount("janinebar", contributorServerClosedIssues, contributorServerOpenIssues)

	syed := newReturnResponse("Majjalpee", contributorResponsesServer, contributorResponsesApp)
	syed.Issues = newIssueCount("Majjalpee", contributorAppClosedIssues, contributorAppOpenIssues)
	syed.Issues += newIssueCount("Majjalpee", contributorServerClosedIssues, contributorServerOpenIssues)

	rafael := newReturnResponse("RafaelHerrejon", contributorResponsesServer, contributorResponsesApp)
	rafael.Issues = newIssueCount("RafaelHerrejon", contributorAppClosedIssues, contributorAppOpenIssues)
	rafael.Issues += newIssueCount("RafaelHerrejon", contributorServerClosedIssues, contributorServerOpenIssues)

	andrea := newReturnResponse("ngynandrea", contributorResponsesServer, contributorResponsesApp)
	andrea.Issues = newIssueCount("ngynandrea", contributorAppClosedIssues, contributorAppOpenIssues)
	andrea.Issues += newIssueCount("ngynandrea", contributorServerClosedIssues, contributorServerOpenIssues)

	var returnResponses [7]ReturnResponse

	returnResponses[0] = *andrea
	returnResponses[1] = *janine
	returnResponses[2] = *kenny
	returnResponses[3] = *luke
	returnResponses[4] = *rafael
	returnResponses[5] = *syed
	returnResponses[6] = *yasira

	js, errJS := json.Marshal(returnResponses)
	if errJS != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in decoding mongo document: %v", err)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

	fmt.Println()
	fmt.Println("contributors' response posted!")
	fmt.Println()

	return
}

func newReturnResponse(username string, serverResponse []ContributorResponse, appResponse []ContributorResponse) *ReturnResponse {

	r := ReturnResponse{Username: username}
	r.Contributions = 0

	for i := 0; i < len(serverResponse); i++ {
		if serverResponse[i].Username == r.Username {
			r.Contributions += serverResponse[i].Contributions
		}
	}

	for i := 0; i < len(appResponse); i++ {
		if appResponse[i].Username == r.Username {
			r.Contributions += appResponse[i].Contributions
		}
	}

	return &r

}

func newIssueCount(username string, closedIssues []IssueResponse, openIssues []IssueResponse) int {

	issuenumber := 0

	for i := 0; i < len(closedIssues); i++ {
		if closedIssues[i].User.Login == username {
			issuenumber++
		}
	}

	for i := 0; i < len(openIssues); i++ {
		if openIssues[i].User.Login == username {
			issuenumber++
		}
	}

	return issuenumber

}
