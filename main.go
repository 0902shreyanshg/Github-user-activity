/*
	main.go : 	github user activity CLI
	Usage: 		github-activity <username>
*/


// ----------------------------------------------------------------------------------------------------------------------------------
// NOTE: CLI
// 
// 		USER INPUT:
// 		./github-activity <username>
//		GO receives this via `os.args` - a SLICE of strings (SIMILAR to a LIST in python)
//		- os.args[0] : 			program name 
// 		- os.args[1] : 			username
// 
//
// NOTE: API
//
// 		HTTP:
// 		- net/http : 			GO's built-in HTTP package - no external library needed
//		- http.Client : 		HTTP client object that sends HTTP requests to a server
// 		- Timeout : 			max time to wait for a response before giving up
// 		- resp.StatusCode : 	number in response telling you the outcome
// 			- 200 : 				success
// 			- 404 : 				user not found
// 			- 403 : 				rate limited
//		- resp.Body : 			actual content of the response (JSON data)
// 
//
// NOTE: FILE-SYSTEM
//
// 		JSON:
// 		GO doesnot parse JSON into dicitionaries like python does. 
//		We define STRUCT that mirror the shape of JSON
// 		- json.Unmarshal() : 	maps JSON fields into STRUCT fields
// 		- json : "type" : 		tells GO which JSON key maps to which field
// 
// ----------------------------------------------------------------------------------------------------------------------------------


// * PACKAGE DECLARATION & IMPORTS
//
// 		Python files are just scripts - you run them directly with `python3 file.py`
// 		In GO, every file starts with a package declaration - it's an executable program not a library
// 		GO is a compiled language - could contain multiple files and packages, hence, GO needs to which package is the entry point 
// 		`package main` is that signal
// 		
// 		os : 					CLI 
// 		net/http : 				API
// 		encoding/json : 		JSON parsing
// 		io : 					reading `resp.Body`
// 		fmt : 					printing output
// 		time : 					
//		strings : 				

package main

import (
	"os"
	"net/http"
	"encoding/json"
	"io"
	"fmt"
	"time"
	"strings"
)


// * STRUCTS 
// 
// 		Struct definition is just the shape of the data before it can parse it, the data comes from API
// 		
//	*	TO DISPLAY
// 		- Pushed 3 commits to kamranahmedse/developer-roadmap 			- type, Payload.Size, Repo.Name
//		- Opened a new issue in kamranahmedse/developer-roadmap			- type, Payload.Action, Repo.Name
// 		- Starred kamranahmedse/developer-roadmap						- type, Repo.Name
//		- Created a new branch in torvalds/ScrollWheel					- type, Payload.ref_type, Repo.Name
//
//	*	REST API endpoints for Github events
// 		- https://docs.github.com/en/rest/activity/events?apiVersion=2022-11-28 > List events for the authenticated user > response schema
// 		OR
//		- https://api.github.com/users/torvalds/events 	// The link is just the GitHub API endpoint with a real username plugged in so you get actual data back instead of an empty response
// 
// 	*	ONE EVENT OBJECT : 
//  	{
//     		"id":         ...
//     		"type":       "PushEvent"          		- // // NEED: tells us what kind of event
//     		"actor":      { ... }              		- SKIP: we already know the username
//     		"repo":       { ... }              		- // // NEED: which repo — but go inside it
//     		"payload":    { ... }              		- // // NEED: event details — but go inside it
//     		"public":     true                 		- SKIP
//     		"created_at": ...                  		- SKIP for now
// 		}
// 		"repo": {
//     			"id":   ...                        	- SKIP
//     			"name": "torvalds/GuitarPedal"     	- // // NEED: this is what we display
//     			"url":  ...                        	- SKIP
// 		}
// 		PushEvent payload
// 		"payload": {
//     			"push_id": ...                     	- SKIP
//     			"ref":     ...                     	- SKIP
//     			"size":    3                       	- // // NEED: number of commits pushed
//     			"commits": [ ... ]                 	- SKIP
// 		}
// 		IssueCommentEvent / PullRequestEvent payload  
// 		"payload": {
//     			"action": "created"               	- // // NEED: opened/closed/created
// 		}
// 		CreateEvent payload
// 		"payload": {
//     			"ref_type": "branch"              	- // // NEED: what was created
// 		}
//
// * 	SYNTAX
//		

type Event struct {
	Type 		string		`json:"type"`
	Repo 		Repo 		`json:"repo"`
	Payload		Payload		`json:"payload"`
}

type Repo struct {
	Name 		string 		`json:"name"`
}

type Payload struct {
	Size 		int 		`json:"size"`
	Action 		string		`json:"action"`
	RefType		string		`json:"ref_type"`
}


// * MAIN FUNCTION FLOW 
// 
// INPUT :
//		I. CLI ARGUMENT 
// 		II. BUILD REQUEST URL
// 		III. HTTP CLIENT & REQUEST
// 		IV. RESPONSE VALIDATION
// 
// REPONSE : 
// 		V. READ & PARSE BODY
//			
// OUTPUT :
// 		VI. DISPLAY

func main() {

	// * I. CLI ARGUMENT 
	// 		EDGE CASE : need atleast 2 elements in os.Args
	// 		- fmt.Println :							fmt is printing output package & Println is the function inside it
	//		- os.Exit(1) :							Stop the program immediately (1 - program ended due to an error, 0 - clean exit; GO requires you to be explicit)
	//		- username := os.Args[1] : 				:= DECLARES a new variable & ASSIGNS a value to it in one step

	if len(os.Args) < 2 {
		fmt.Println("Please provide a github username")
		os.Exit(1)
	}
	username := os.Args[1]


	// * II. BUILD REQUEST URL

	url := "https://api.github.com/users/" + username + "/events"


	// * III. HTTP CLIENT & REQUEST
	//		i. create http.Client() with Timeout (10 s)
	// 			- http.Client{} : 						A struct (struct can store stuff apart from API data) that holds settings for how to make new requests
	//		ii. build GET request
	// 			- http.NewRequest("GET", url, nil) : 	NewRequest returns 2 things request & error; GET requests don't send any data to the server - nil
	//		iii. add authorization header from environment variable
	// 			- os.Getenv("GITHUB_TOKEN") : 			reads "GITHUB_TOKEN" from your system, attaches it to the request as authorization header; used to raise rate limit from 60 to 5,000 requests/hour
	//			- req.Header.Set("Authorization", 
	// 							"Bearer "+token)   :	attach token to the request header before sending; setting up one header named "Authorisation" with the value "Bearer"+token
	//		iv. response
	//			client, req, resp all come together
	// 			- client.Do(req) : 						returns response & error
	// 		v. closing body
	// 			- defer resp.Body.Close() : 			always close body when done reading

	client := http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request: ", err)
		os.Exit(1)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt. Println("Error making request: ", err)
		os.Exit(1)
	}
	defer resp.Body.Close()


	// * IV. RESPONSE VALIDATION 

	if resp.StatusCode == 404 {
		fmt.Println("User not found")
		os.Exit(1)
	}
	else if resp.StatusCode == 403 {
		fmt.Println("Rate limit exceeded. Set GITHUB_TOKEN environment variable to increase limit.")
		os.Exit(1)
	}
	else if resp.StatusCode == 200 {
		fmt.Println("API Error: ", resp.StatusCode)
		os.Exit(1)
	}


	// * V. READ & PARSE BODY
	// 		i. Reading raw bytes from response
	// 			- body, err := io.ReadAll(resp.Body) : 		body gets declared, err gets reassigned (not declared, it was declared above)
	//		ii. mapping response to []Event struct
	// 			- err = json.Unmarshal(body, &events) : 	we reuse err as json.Unmarshal can also fail if json is malformed or doesn't match struct

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response: ", err)
		os.Exit(1)
	}

	var events []Event
	err = json.Unmarshal(body, &events)
	if err != nil {
		fmt.Println("Error parsing response", err)
		os.Exit(1)
	}


	// * VI. DISPLAY
	// 		- for _, event := range events : 				"_" ; "range events" loops over the SLICE

	for _, event := range events {
		fmt.Println(formatEvent(event))
	}
}


// * FORMAT EVENT (HELPER FUNCTION)
// 		formatEvent(event Event) : 							the function takes one argument event of type Event (struct)
// 		switch event.Type: 									
//		fmt.Sprintf() : 									
//		both action
// 		capitalize(s string) string
// 		strings.toUpper(s[:1] + s[1:]

func formatEvent(event Event) string {
	switch event.Type {
	case "PushEvent":
		return fmt.Sprintf("Pushed %d commits to %s", event.Payload.Size, event.Repo.Name)
	case "IssuesEvent":
		return fmt.Sprintf("%s is an issue in %s", capitalize(event.Payload.Action), event.Repo.Name)
	case "PullRequestEvent":
		return fmt.Sprintf("%s a pull request in %s", capitalize(event.Payload.Action), event.Repo.Name)
	case "WatchEvent":
		return fmt.Sprintf("Starred %s", event.Repo.Name)
	case "CreateEvent":
		return fmt.Sprintf("Created a new %s in %s", event.Payload.RefType, event.Repo.Name)
	case "ForkEvent":
		return fmt.Sprintf("Forked %s", event.Repo.Name)
	default:
		return fmt.Sprintf("%s in %s", event.Type, event.Repo.Name)
	}
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.toUpper(s[:1] + s[1:])
}