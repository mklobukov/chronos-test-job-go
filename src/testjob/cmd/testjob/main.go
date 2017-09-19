//Test job for Chronos Go SDK
//Provide path to config.json as a command line argument
//Initializes counter to the arguments from the job descriptor
//and repeatedly updates provided job's status with the incremented counter 
package main

import (
  "encoding/json"
  "fmt"
  "os"
  "time"
  "github.com/iris-platform/chronos-go-sdk"
  "strconv"
)

func main() {
  var config chronossdk.Config
  if (len(os.Args) < 2) {
    fmt.Println("Provide path to config.json as an argument")
    return
  }
  pathToConfig := os.Args[1]
  configFile, err := os.Open(pathToConfig)
  defer configFile.Close()
  if err != nil {
    fmt.Println(err.Error())
    return
  }

  jsonParser := json.NewDecoder(configFile)
  jsonParser.Decode(&config)

  //initialize counter to argument from the job description
  counter := initializeCounter(&config)
  fmt.Println("Counter initialized to: ", counter)

  //update custom job status every 250ms for a few seconds
  ticker := time.NewTicker(time.Millisecond * 250)
  go func() {
    for _ = range ticker.C {
      chronossdk.UpdateJobStatus(&config, "Custom status = " + strconv.Itoa(counter))
      counter = counter + 1
      fmt.Println("Updating custom status to: ", counter)
    }
  }()
  time.Sleep(time.Millisecond * 3300)
  ticker.Stop()
}

func initializeCounter(config *chronossdk.Config) (int) {
  argsString, err := chronossdk.GetJobArgs(config)
  if err != nil {
    fmt.Println("Could not get job args: ", err.Error())
    return 0
  } else {
    //unmarshal args string into an object
    argsJSON := make(map[string]interface{})
    err = json.Unmarshal([]byte(argsString), &argsJSON)
    if err != nil {
      fmt.Println("Could not retrieve initial counter value. Initializing to 0.")
      return 0
    }
    counter := int(argsJSON["counterInit"].(float64))
    return counter
  }
}
