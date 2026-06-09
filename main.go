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

package main

import (
	"os"
	"net/http"
	"encoding/json"
	"io"
	"fmt"
)


// * STRUCTS 
// 
// 		Struct definition is just the shape of the data before it can parse it, the data comes from API
//		REST API endpoints for Github events
// 		- https://docs.github.com/en/rest/activity/events?apiVersion=2022-11-28 > List events for the authenticated user > response schema
// 		OR
//		- https://api.github.com/users/torvalds/events 	// The link is just the GitHub API endpoint with a real username plugged in so you get actual data back instead of an empty response
// 		
//		TO DISPLAY
// 		- Pushed 3 commits to kamranahmedse/developer-roadmap 			- type, payload.size, repo.name
//		- Opened a new issue in kamranahmedse/developer-roadmap			- type, payload.action, repo.name
// 		- Starred kamranahmedse/developer-roadmap						- type, repo.name
//		- Created a new branch in torvalds/ScrollWheel					- type, payload.ref_type, repo.name
//
// 		ONE EVENT OBJECT: 
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


// 