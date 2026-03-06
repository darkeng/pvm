package commands

import (
	"bufio"
	"fmt"
	"hjbdev/pvm/theme"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Laravel() {
	theme.Title("Configuring PHP for Laravel")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalln(err)
	}

	batPath := filepath.Join(homeDir, ".pvm", "bin", "php.bat")

	if _, err := os.Stat(batPath); os.IsNotExist(err) {
		theme.Error("No active PHP version found. Please run 'pvm use <version>' first.")
		return
	}

	// Read php.bat to find the active version path
	file, err := os.Open(batPath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	var phpExePath string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "set filepath=") {
			phpExePath = strings.TrimPrefix(line, "set filepath=")
			phpExePath = strings.Trim(phpExePath, "\"")
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalln(err)
	}

	if phpExePath == "" {
		theme.Error("Could not determine the active PHP version from php.bat")
		return
	}

	phpDir := filepath.Dir(phpExePath)
	phpIniPath := filepath.Join(phpDir, "php.ini")
	phpIniDevPath := filepath.Join(phpDir, "php.ini-development")

	// If php.ini doesn't exist, create it from php.ini-development
	if _, err := os.Stat(phpIniPath); os.IsNotExist(err) {
		if _, err := os.Stat(phpIniDevPath); os.IsNotExist(err) {
			theme.Error("Could not find php.ini or php.ini-development in the active PHP directory.")
			return
		}

		theme.Info("php.ini not found. Copied from php.ini-development.")
		err = copyFile(phpIniDevPath, phpIniPath)
		if err != nil {
			log.Fatalln("Failed to copy php.ini-development:", err)
		}
	} else {
		theme.Info("Found existing php.ini")
	}

	// Read and modify php.ini
	iniBytes, err := os.ReadFile(phpIniPath)
	if err != nil {
		log.Fatalln(err)
	}

	iniContent := string(iniBytes)
	// Some php.ini files use \r\n
	lines := strings.Split(iniContent, "\n")

	// Extensions and settings to enable
	extensions := []string{
		"openssl", "sodium", "fileinfo", "curl", "mbstring", "pdo_mysql", "gd", "zip", "intl",
	}

	var newLines []string

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Enable extension_dir = "ext"
		if trimmedLine == `;extension_dir = "ext"` || trimmedLine == `; extension_dir = "ext"` {
			// preserve any carriage returns for this line
			suffix := ""
			if strings.HasSuffix(line, "\r") {
				suffix = "\r"
			}
			newLines = append(newLines, `extension_dir = "ext"`+suffix)
			continue
		}

		// Enable required extensions
		matchedExt := false
		for _, ext := range extensions {
			if trimmedLine == fmt.Sprintf(";extension=%s", ext) || trimmedLine == fmt.Sprintf("; extension=%s", ext) {
				suffix := ""
				if strings.HasSuffix(line, "\r") {
					suffix = "\r"
				}
				newLines = append(newLines, fmt.Sprintf("extension=%s", ext)+suffix)
				matchedExt = true
				break
			}
		}

		if !matchedExt {
			newLines = append(newLines, line)
		}
	}

	// Write the modified content back
	newIniContent := strings.Join(newLines, "\n")
	err = os.WriteFile(phpIniPath, []byte(newIniContent), 0644)
	if err != nil {
		log.Fatalln("Failed to write to php.ini:", err)
	}

	theme.Success("Successfully configured php.ini for Laravel.")
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	return err
}
