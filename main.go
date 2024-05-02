package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jlaffaye/ftp"
)

type FtpData struct {
	Host     string
	User     string
	Password string
}

func main() {
	if len(os.Args) < 5 {
		printUsage("Invalid number of arguments")
		os.Exit(1)
	}

	ftpData := FtpData{
		Host:     os.Args[1],
		User:     os.Args[2],
		Password: os.Args[3],
	}

	maxAgeInDays, err := strconv.Atoi(os.Args[4])
	if err != nil {
		printUsage(fmt.Sprintf(`Invalid <max_age_in_days> argument: "%s"`, os.Args[4]))
		os.Exit(1)
	}

	fileRegexArg := os.Args[5]

	conn, err := ftp.Dial(ftpData.Host, ftp.DialWithTimeout(10*time.Second))
	if err != nil {
		fmt.Println("Connection error...", err)
		os.Exit(1)
	}

	err = conn.Login(ftpData.User, ftpData.Password)
	if err != nil {
		fmt.Println("Login failed", err)
		os.Exit(1)
	}

	conn.Type(ftp.TransferTypeBinary)

	fileList, err := conn.List("/")
	if err != nil {
		fmt.Println("Error while listing files...", err)
		os.Exit(1)
	}

	fileRegex := regexp.MustCompile(fileRegexArg)
	regexNames := fileRegex.SubexpNames()

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	var filesToDelete []string
	for _, fileEntry := range fileList {
		if fileRegex.MatchString(fileEntry.Name) {
			subMatches := fileRegex.FindAllSubmatch([]byte(fileEntry.Name), 1)

			fileDateValues := make(map[string]int)
			for i, name := range regexNames {
				value, _ := strconv.Atoi(string(subMatches[0][i]))
				fileDateValues[name] = value
			}

			fileDate := time.Date(
				fileDateValues["year"], time.Month(fileDateValues["month"]), fileDateValues["day"],
				0, 0, 0, 0,
				time.UTC,
			)

			if today.Sub(fileDate).Hours()/24 > float64(maxAgeInDays) {
				filesToDelete = append(filesToDelete, fileEntry.Name)
			}
		}
	}

	if len(filesToDelete) > 0 {
		fmt.Println("Files to be deleted...")
		for _, fileName := range filesToDelete {
			fmt.Printf("  %s", fileName)
			err = conn.Delete(fileName)
			if err != nil {
				fmt.Println(" Error!")
			} else {
				fmt.Println(" OK")
			}
		}
	} else {
		fmt.Println("No files found to be deleted")
	}

	if err := conn.Quit(); err != nil {
		fmt.Println("Disconnection error...")
		os.Exit(1)
	}
}

func printUsage(message string) {
	fmt.Println(message)
	fmt.Println("Usage: ftp-remove <host:port> <user> <password> <max_age_in_days> <file_regex>")
}
