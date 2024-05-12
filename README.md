## LOGZY - The lightweight log notification handler for your service.

enable notifications for your specific log message

## How it works

Tailing the log files from your specified log file and send slack notifications using filtered log messages regex.

## Things currently developed

- Tailing the log files
- Adjust tailing log files according to the log rotation
- Indexing the log files and identified current state in the log file
- Filtering the error logs
- push notifications to the Slack channel

## Plans to release on next versions

- Make support for other notification platforms like PagerDuty, Discord.
- improve logs filtering functionality using more custom regex filter functions.
- Testing more for the race conditions.
- develop an agent to gather metrics as a background task with low usage.
- develop dashboard to visualize the metrics and logs issues.

## Option 01: Make a executable file and run as a system demon

- First you need to make a executable file using go build command.

```bash
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o logzy .
```

- Then you need to configure this executable file run as system demon.

- After setup you need to copy all the configuration files that included in `ops` directory.

- Make sure to create a log file named as `logzy.log` as mentioned log configuration file paths in ops config log file.

## Option 02: Run as a dockerfile

- We have provided a Dockerfile and you can customize that as you need

```bash
docker build -t logzy && docker run -d logzy
```

## Things need to setup before run the project

- First you need to change the Slack incoming webhook URL property named as `slack-uri` inside `ops/app/config-local.yaml` file.

- Then you need to mention the log file name you need to tail inside the `ops/app/config-local.yaml` file property named as `log-file-name`.

- There is two log locations you need to add the setup file. first one is `log-location` property. This is the log file location you need to tail.

- Second one is `logger-file-location`. This is the log file location related to your application logs. If there is no log file location mentioned in the `outputPaths` property in `ops/log/config-{ENV}.yaml` then you need to add a location with file name like `logs/logzy.log` and check if there is a file in your mentioned place.

- for locally run the project you can run `ENV=local go run main.go` command. Install go setup before run the command.