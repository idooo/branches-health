package core

import (
	"strings"
	"bytes"
	"log"
	"regexp"
	"time"
	"fmt"
	"strconv"
	"os"
	"os/exec"
	"io/ioutil"
	"path/filepath"
	"github.com/boltdb/bolt"
)

func runGitCommand(args []string) string {
	cmd := exec.Command("git", args...)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.Fatalf("Command failed %s : %s", strings.Join(args, " "), err)
	}
	return out.String()
}

func formatBranchData (repoName, branchName string, isMerged, isOutdated bool) Branch {
	output := runGitCommand([]string{"show", "--format=%ct,%an,%cn", branchName})
	info := strings.Split(strings.Split(output, "\n")[0], ",")

	lastAuthor := info[1]

	var lastUpdated time.Time
	lastUpdatedInt, errParseInt := strconv.ParseInt(info[0], 10, 64)
	if errParseInt != nil {
		fmt.Printf("Can't convert date from git log %s/%s - %s", repoName, branchName, info[0])
		lastUpdated = time.Now()
	} else {
		lastUpdated = time.Unix(lastUpdatedInt, 0)
	}

	return Branch{
		repoName,
		strings.Replace(branchName, "origin/", "", -1),
		repoName + "/" + branchName,
		isMerged,
		isOutdated,
		lastAuthor,
		lastUpdated,
	}
}

func GetBranchesInfoForRepo(repoName string) []Branch {

	// Change folder to the current to prevent errors
	// if previously working directory has been deleted
	currentDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	errChdir := os.Chdir(currentDir)
	if errChdir != nil {
		log.Fatalf("Can't change folder to %s", currentDir)
	}

	// Create a temporary folder
	tmpDir, errTmp := ioutil.TempDir(os.TempDir(), "branchitto")
	if errTmp != nil {
		log.Fatal(errTmp)
	}

	// Clone repo
	runGitCommand([]string{"clone", "-q", repoName, tmpDir})

	// Ignore some branches
	master := regexp.MustCompile("(origin/HEAD|origin/master)")
	errChdir = os.Chdir(tmpDir)
	if errChdir != nil {
		log.Fatalf("Can't change folder to %s", tmpDir)
	}

	// Get info about merged branches
	branches := make([]Branch, 0)
	merged := runGitCommand([]string{"branch", "-r", "--merged"})
	for _, branchName := range strings.Split(merged, "\n") {
		if master.MatchString(strings.TrimSpace(branchName)) { continue }
		if len(branchName) == 0 { continue }
		branches = append(branches, formatBranchData(repoName, strings.TrimSpace(branchName), true, true))
	}

	// Get info about not merged branches
	notMerged := runGitCommand([]string{"branch", "-r", "--no-merged"})
	for _, branchName := range strings.Split(notMerged, "\n") {
		if len(branchName) == 0 { continue }
		logLastMonth := runGitCommand([]string{"log", "-1", "--since='1 month ago'", "-s", strings.TrimSpace(branchName), "--oneline"})
		isOutdated := len(logLastMonth) == 0
		branches = append(branches, formatBranchData(repoName, strings.TrimSpace(branchName), false, isOutdated))
	}

	// Clean up
	os.RemoveAll(tmpDir)

	return branches
}


// Iterates through repositories, gets data and saves it to a database
func GetBranchesInfoForRepos (repositories []string, database *bolt.DB) {

	for _, repoName := range repositories {
		branches := GetBranchesInfoForRepo(repoName)

		for _, branch := range branches {
			log.Printf("Get information about: %s/%s", repoName, branch.Name)
			branch.Save(database)
		}
	}

}
