## LOGZY - The lightweight log notification handler for your service.

enable notifications for your specific log message

## How it works

Tailing the log files from your specified log file and send slack notifications using filtered log messages regex.

## Things currently developed

- Tailing the log files
- Maintain a log rotation
- Indexing the log files
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
