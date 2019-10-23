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

func GetContributorsHandler(w http.ResponseWriter, r *http.Request) {

	//TODO: cache maybe?

	///////////////////////////////////////
	// GETTING CONTRIBUTIONS FROM SERVER //
	///////////////////////////////////////

	//get users from the HTTP link
	resp, err := http.Get("https://api.github.com/repos/MunchApp/munchserver/contributors")
	if err != nil {
		fmt.Println("error:", err)
	}

	//close the response
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	//print english for terminal debug
	// englishBody := string(body)
	// fmt.Println(englishBody)

	//Create an array and print the contents of the array
	var contributorResponsesServer []ContributorResponse
	jsonErr := json.Unmarshal(body, &contributorResponsesServer)
	if jsonErr != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("ContributorResponses : %+v", contributorResponsesServer)

	//////////////////////////////////////////
	// GETTING CONTRIBUTIONS FROM munch-app //
	//////////////////////////////////////////

	//get users from the HTTP link
	respApp, errApp := http.Get("https://api.github.com/repos/MunchApp/munch-app/contributors")
	if errApp != nil {
		fmt.Println("error HEEHRHR:", err)
	}

	//close the response
	defer respApp.Body.Close()
	bodyApp, errApp := ioutil.ReadAll(respApp.Body)

	//print english for terminal debug
	englishBodyApp := string(bodyApp)
	fmt.Println(englishBodyApp)

	//Create an array and print the contents of the array
	var contributorResponsesApp []ContributorResponse
	jsonErrApp := json.Unmarshal(bodyApp, &contributorResponsesApp)
	if jsonErrApp != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Printf("ContributorResponses : %+v", contributorResponsesApp)
	}

	//////////////////////////////////
	// CREATING THE RETURN RESPONSE //
	//////////////////////////////////

	fmt.Println()
	fmt.Println("TEST - printing: ", contributorResponsesApp[0])

	yasira := newReturnResponse("yasirayounus", contributorResponsesServer, contributorResponsesApp)
	var kenny = newReturnResponse("kftang", contributorResponsesServer, contributorResponsesApp)
	var luke = newReturnResponse("Lmnorrell99", contributorResponsesServer, contributorResponsesApp)
	var janine = newReturnResponse("janinebar", contributorResponsesServer, contributorResponsesApp)
	var syed = newReturnResponse("majjalpee", contributorResponsesServer, contributorResponsesApp)
	var rafael = newReturnResponse("RafaelHerrejon", contributorResponsesServer, contributorResponsesApp)
	var andrea = newReturnResponse("ngynandrea", contributorResponsesServer, contributorResponsesApp)

	var returnResponses [7]ContributorResponse

	returnResponses[0] = *andrea
	returnResponses[1] = *janine
	returnResponses[2] = *kenny
	returnResponses[3] = *luke
	returnResponses[4] = *rafael
	returnResponses[5] = *syed
	returnResponses[6] = *yasira

	fmt.Println()
	fmt.Printf("returnResponses : %+v", returnResponses)

	js, errJS := json.Marshal(returnResponses)
	if errJS != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, "Error in decoding mongo document: %v", err)
		return
	}
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return
}

func newReturnResponse(username string, serverResponse []ContributorResponse, appResponse []ContributorResponse) *ContributorResponse {

	r := ContributorResponse{Username: username}
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
