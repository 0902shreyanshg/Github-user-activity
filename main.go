/*

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

