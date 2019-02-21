package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/user"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/boltdb/bolt"
)

func runGitCommand(args []string) string {
	cmd := exec.Command("git", args...)

	var out bytes.Buffer
	var error bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &error

	err := cmd.Run()
	if err != nil {
		fmt.Errorf("Command failed %s : %s", strings.Join(args, " "), err)
	}
	if len(error.String()) != 0 {
		fmt.Printf("Command failed %s : %s\n", strings.Join(args, " "), error.String())
	}
	return out.String()
}

func formatBranchData(repoName, branchName string, isMerged, isOutdated bool) Branch {
	output := runGitCommand([]string{"show", "--format=%ct,%an,%cn", branchName})
	info := strings.Split(strings.Split(output, "\n")[0], ",")

	lastAuthor := info[1]

	var lastUpdated time.Time
	lastUpdatedInt, errParseInt := strconv.ParseInt(info[0], 10, 64)
	if errParseInt != nil {
		fmt.Errorf("Can't convert date from git log %s/%s - %s", repoName, branchName, info[0])
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

func resetDir() error {
	usr, err := user.Current()
    if err != nil {
        log.Fatal( err )
    }
	return os.Chdir(usr.HomeDir )
}

func GetBranchesInfoForRepo(repoName string, ignore *regexp.Regexp) []Branch {
	fmt.Printf("Getting information about %s\n", repoName)
	branches := make([]Branch, 0)

	// Change folder to the current to prevent errors
	// if previously working directory has been deleted
	if errChdir := resetDir(); errChdir != nil {
		fmt.Errorf("Can't reset folder %s", errChdir)
		return branches
	}

	// Create a temporary folder
	tmpDir, errTmp := ioutil.TempDir(os.TempDir(), "branches-health")
	if errTmp != nil {
		fmt.Errorf("Can't create temporary folder %s - %s", tmpDir, errTmp)
		return branches
	}

	// Clone repo and open its folder
	runGitCommand([]string{"clone", "-q", repoName, tmpDir})
	if errChdir := os.Chdir(tmpDir); errChdir != nil {
		fmt.Errorf("Can't change folder to %s", tmpDir)
		return branches
	}

	// Get info about merged branches
	merged := runGitCommand([]string{"branch", "-r", "--merged"})
	for _, branchName := range strings.Split(merged, "\n") {
		if len(branchName) == 0 || ignore.MatchString(strings.TrimSpace(branchName)) {
			continue
		}
		branches = append(branches, formatBranchData(repoName, strings.TrimSpace(branchName), true, true))
	}

	// Get info about not merged branches
	notMerged := runGitCommand([]string{"branch", "-r", "--no-merged"})
	for _, branchName := range strings.Split(notMerged, "\n") {
		if len(branchName) == 0 || ignore.MatchString(strings.TrimSpace(branchName)) {
			continue
		}
		logLastMonth := runGitCommand([]string{"log", "-1", "--since='1 month ago'", "-s", strings.TrimSpace(branchName), "--oneline"})
		isOutdated := len(logLastMonth) == 0
		branches = append(branches, formatBranchData(repoName, strings.TrimSpace(branchName), false, isOutdated))
	}

	// Clean up
	os.RemoveAll(tmpDir)

	return branches
}

// Iterates through repositories, gets data and saves it to a database
func GetBranchesInfoForRepos(repositories []string, branchesToIgnore []string, database *bolt.DB) {

	if err := CleanBranches(database); err != nil {
		fmt.Errorf("Can't clean branches database: %s", err)
	}

	// Generate ignore string to skip some branches
	branchesToIgnoreStr := "(" + strings.Join(append(branchesToIgnore, "origin/HEAD", "origin/master"), "|") + ")"
	branchesPatternToIngore, regexpErr := regexp.Compile(branchesToIgnoreStr)
	if regexpErr != nil {
		log.Fatalf("Invalid list of ignored branches, failed to create regexp %s", branchesToIgnoreStr)
	}

	for _, repoName := range repositories {
		branches := GetBranchesInfoForRepo(repoName, branchesPatternToIngore)

		for _, branch := range branches {
			log.Printf(">>> %s/%s", repoName, branch.Name)
			branch.Save(database)
		}
	}

}
