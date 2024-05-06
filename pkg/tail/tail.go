package tail

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"logzy/config"
	"logzy/entity"
	logger "logzy/pkg/log"
	restClient "logzy/pkg/restclient"

	"github.com/hpcloud/tail"
	"go.uber.org/zap"
)

var (
	log *zap.SugaredLogger = logger.GetLogger().Sugar()
	mtx sync.Mutex
)

// readLastKnownOffset returns the last known offset from a file.
func readLastKnownOffset(filePath string) (int64, error) {
	data, err := os.ReadFile(filePath)
	if err != nil || len(data) == 0 {
		return 0, nil
	}
	return strconv.ParseInt(string(data), 10, 64)
}

// writeLastKnownOffset writes the last known offset to a file.
func writeLastKnownOffset(filePath string, offset int64) error {
	return os.WriteFile(filePath, []byte(strconv.FormatInt(offset, 10)), 0644)
}

// TailLogFile handles the tailing of a log file, capturing and processing error logs.
func TailLogFile(logFilePath, offsetFilePath string, wg *sync.WaitGroup) {
	defer wg.Done()

	if lastOffset, err := manageOffset(logFilePath, offsetFilePath); err != nil {
		return
	} else {
		mtx.Lock() // lock the mutex to prevent other goroutines from entering
		captureErrors(logFilePath, offsetFilePath, lastOffset)
		mtx.Unlock() // Ensure the mutex is unlocked when the function exits
	}
}

// manageOffset retrieves and updates the last known offset based on file path.
func manageOffset(logFilePath, offsetFilePath string) (int64, error) {
	_, err := os.Stat(logFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warnf("Log file %s does not exists, starting from the begining", logFilePath)
			return 0, nil
		}
		log.Errorf("Unable to access log file %s: %v", logFilePath, err)
		return 0, err
	}

	lastOffset, err := readLastKnownOffset(offsetFilePath)
	if err != nil {
		log.Errorf("Error reading last offset for %s: %v", logFilePath, err)
		return 0, err
	}
	log.Infof("Starting to tail from offset %d for file %s", lastOffset, logFilePath)
	return lastOffset, nil
}

// captureErrors starts capturing errors from the log file starting from the last offset.
func captureErrors(logFilePath, offsetFilePath string, lastOffset int64) {
	var fileSize int64

	t, err := tail.TailFile(logFilePath, tail.Config{
		Follow:    true,
		Location:  &tail.SeekInfo{Offset: lastOffset, Whence: io.SeekStart},
		MustExist: true,
		ReOpen:    true,
	})
	if err != nil {
		log.Errorf("Error tailing file %s: %v", logFilePath, err)
		return
	}

	fileInfo, err := os.Stat(logFilePath)
	if err == nil {
		fileSize = fileInfo.Size()
	}

	var errorDescription []string
	capturing := false

	for line := range t.Lines {
		currentSize := checkFileSize(logFilePath, &fileSize)
		if currentSize == -1 { // If file is rotated
			processLine(line, &capturing, &errorDescription, offsetFilePath, t)
			continue // Skip to the next iteration to avoid processing lines from an old file instance
		}

		processLine(line, &capturing, &errorDescription, offsetFilePath, t)
	}
}

// check file size
func checkFileSize(filePath string, lastKnowsSize *int64) int64 {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Errorf("Cannot stat the log file: %s", err)
		return *lastKnowsSize // Return the last known size if error occurs
	}
	currentSize := fileInfo.Size()
	if currentSize < *lastKnowsSize {
		log.Info("Log file has been rotated.")
		*lastKnowsSize = 0
		return -1
	}
	*lastKnowsSize = currentSize
	return currentSize
}

// processLine processes each line of the log, checking for errors and managing state.
func processLine(line *tail.Line, capturing *bool, errorDescription *[]string, offsetFilePath string, t *tail.Tail) {
	regexPatterns := config.Parse()

	startingRegexps, err := MatchingStartRegexPatterns(regexPatterns.AppConfig.StartErrorPatterns)
	if err != nil {
		log.Errorf("error compiling regex patterns. %v", err)
		return
	}

	// identify pattern one by one with line text
	for _, pattern := range startingRegexps {
		if pattern.MatchString(line.Text) {
			// Start capturing if the line starts a new error
			*capturing = true
			*errorDescription = []string{line.Text}
			log.Info("Starting new error capture")
		}
	}

	if len(*errorDescription) > 0 {
		log.Info("Sending error trace to Slack")
		sendRestCall(errorDescription)
	} else {
		log.Warn("No error description to send")
	}

	// After processing each line, update the offset to ensure correct continuation after restarts
	if err := updateOffset(offsetFilePath, t); err != nil {
		log.Errorf("Failed to update offset file: %v", err)
	}
}

func sendRestCall(errorDescription *[]string) {
	cfg := config.Parse()
	timeDuration := time.Duration(cfg.AppConfig.ResponseTimeout * int(time.Second))
	client := restClient.GetInstance(&timeDuration)

	// Convert slice of strings to a single string
	var textContent string
	if errorDescription != nil && len(*errorDescription) > 0 {
		textContent = strings.Join(*errorDescription, "\n") // Joining all error descriptions with newlines
	}

	postData := entity.SlackMessage{
		Text: "Error Report",
		Blocks: []entity.Block{
			{
				Type: "divider",
			},
			{
				Type: "header",
				Text: &entity.HeaderText{
					Type:  "plain_text",
					Text:  fmt.Sprintf("Error on %s", cfg.AppConfig.TargetService),
					Emoji: true,
				},
			},
			{
				Type: "rich_text",
				Element: []entity.Element{
					{
						Type:   "rich_text_preformatted",
						Border: 0,
						Elements: []entity.TextObject{
							{
								Type: "text",
								Text: &textContent,
							},
						},
					},
				},
			},
			{
				Type: "section",
				Fields: []entity.MarkdownText{
					{
						Type: "mrkdwn",
						Text: fmt.Sprintf("*Captured At:*\n%s", time.Now().Format("2006-01-02 15:04:05")),
					},
				},
			},
			{
				Type: "divider",
			},
		},
	}

	_, errRestClient := client.Post(cfg.AppConfig.SlackUri, postData)
	if errRestClient != nil {
		log.Errorf("Failed to send error notification: %v", errRestClient)
		return
	}
}

// updateOffset updates the current file offset in the offset file.
func updateOffset(offsetFilePath string, t *tail.Tail) error {
	pos, err := t.Tell()
	if err != nil {
		log.Errorf("Error obtaining file position for %s: %v", offsetFilePath, err)
		return err
	}
	if err := writeLastKnownOffset(offsetFilePath, pos); err != nil {
		log.Errorf("Error writing offset for %s: %v", offsetFilePath, err)
		return err
	}
	return nil
}

// matchingStartRegexPatterns matches the suitable regex pattern
func MatchingStartRegexPatterns(patterns []string) ([]*regexp.Regexp, error) {
	var regexps []*regexp.Regexp
	for _, pattern := range patterns {
		compiled, err := regexp.Compile(pattern)
		if err != nil {
			return nil, err
		}
		regexps = append(regexps, compiled)
	}
	return regexps, nil
}
