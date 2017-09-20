# How to create a job
In order to create Chronos job your code needs to have a Dockerfile that will
build a Docker container with your code.  Once you have that you will need to do
the following to prepare your job for publishing it to Chronos:
* Create chronos.json file in the root of your project.  For your convenience
you may want to commit this file with your project.  You will need the following
information in chronos.json:

```
{
  "version": "1.0.0",
  "name": "chronostest",
  "schedule": "@hourly",
  "repeat": 0,
  "callback": "http://localhost:3007/jobcallback/uno",
  "check_in_threshold": 240
}
```

* version - this is a string identifying version of your job
* name - this is the name of your job and it will be also used as a name for the docker container
* schedule - this string is the schedule specification for your job.  It follows cron tab
format but it supports seconds granulity.  See Cron Expression Format section below.
* repeat - specifies how many times to repeat the job.  -1 means always repeat, 0 run once and do not repeat, 1 means repeat once, 2 repeat twice and so on.
* callback - when job is completed Chronos will execute POST on the specified URL to notify of job completion.
* check_in_threshold - currently not used but it will be required in the next release to specify duration between heartbeats from the job to Chronos.  This field must be present in the configuration file (chronos.json).

# CRON Expression Format
A cron expression represents a set of times, using 6 space-separated fields.
```
Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Seconds      | Yes        | 0-59            | * / , -
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?
Note: Month and Day-of-week field values are case insensitive. "SUN", "Sun", and "sun" are equally accepted.
```
## Special Characters
Asterisk ( * )

The asterisk indicates that the cron expression will match for all values of the field; e.g., using an asterisk in the 5th field (month) would indicate every month.

Slash ( / )

Slashes are used to describe increments of ranges. For example 3-59/15 in the 1st field (minutes) would indicate the 3rd minute of the hour and every 15 minutes thereafter. The form "*\/..." is equivalent to the form "first-last/...", that is, an increment over the largest possible range of the field. The form "N/..." is accepted as meaning "N-MAX/...", that is, starting at N, use the increment until the end of that specific range. It does not wrap around.

Comma ( , )

Commas are used to separate items of a list. For example, using "MON,WED,FRI" in the 5th field (day of week) would mean Mondays, Wednesdays and Fridays.

Hyphen ( - )

Hyphens are used to define ranges. For example, 9-17 would indicate every hour between 9am and 5pm inclusive.

Question mark ( ? )

Question mark may be used instead of '*' for leaving either day-of-month or day-of-week blank.

## Predefined schedules
You may use one of several pre-defined schedules in place of a cron expression.
```
Entry                  | Description                                | Equivalent To
-----                  | -----------                                | -------------
@yearly (or @annually) | Run once a year, midnight, Jan. 1st        | 0 0 0 1 1 *
@monthly               | Run once a month, midnight, first of month | 0 0 0 1 * *
@weekly                | Run once a week, midnight on Sunday        | 0 0 0 * * 0
@daily (or @midnight)  | Run once a day, midnight                   | 0 0 0 * * *
@hourly                | Run once an hour, beginning of hour        | 0 0 * * * *
```
## Intervals
You may also schedule a job to execute at fixed intervals, starting at the time it's added or cron is run. This is supported by formatting the cron spec like this:
```
@every <duration>
```
where "duration" is a string of possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".

For example, "@every 1h30m10s" would indicate a schedule that activates immediately, and then every 1 hour, 30 minutes, 10 seconds.

## Note
Chronos takes into consideration time it takes to run a job.  So each scheduled time
is calculated on initial scheduling and then recalculated after the job is finished.

# Callback Server 
The callback_server directory contains an example program that gets executed when a job is completed. This program prints job information to the console. The information gets sent to the callback server in the body of a POST request. 

```go
func jobCallbackMethod(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	paramsID := params["id"]

	var jobInfo CallbackPostRequestBody
	if err := json.NewDecoder(req.Body).Decode(&jobInfo); err != nil {
		fmt.Println("Error decoding JSON: ", err)
		return
	}

	fmt.Printf("Received job completion at: %v\n", time.Now())
	fmt.Printf("Params ID: %s\n", paramsID)
	fmt.Printf("Job ID: %s\nJob Name: %s\nJob Container ID: %s\nJob Instance ID: %s\n",
						jobInfo.JobID, jobInfo.JobName, jobInfo.JobContainerID, jobInfo.JobInstanceID)
	fmt.Printf("State: %d\nStatus: %d\nStatus Description: %s\n",
						jobInfo.State, jobInfo.Status, jobInfo.StatusDescription)
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/jobcallback/{id}", jobCallbackMethod).Methods("POST")
	log.Fatal(http.ListenAndServe(":3007", router))
}
```
For the full program with imports and CallbackPostRequestBody struct definition, see `./callback-server/main.go`
