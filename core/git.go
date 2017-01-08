package core

import (
	"strings"
	"os/exec"
	"bytes"
	"log"
	"io/ioutil"
	"os"
	"regexp"
	"time"
	"fmt"
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

func getBranchData (repoName, branchName string, isMerged, isOutdated bool) Branch {
	info := runGitCommand([]string{"show", "--format=%cI,%an,%cn", branchName})
	info2 := strings.Split(strings.Split(info, "\n")[0], ",")

	lastUpdated, err := time.Parse(time.RFC3339, info2[0])

	if err != nil {
		fmt.Printf("Can't convert date from git log %s/%s - %s", repoName, branchName, info[0])
		lastUpdated = time.Now()
	}

	return Branch{
		repoName,
		branchName,
		repoName + "/" + branchName,
		isMerged,
		isOutdated,
		info2[1],
		lastUpdated,
	}
}


func GetInfoFromGit(repoName string) []Branch {

	tmpDir, errTmp := ioutil.TempDir(os.TempDir(), "branchitto")
	if errTmp != nil {
		log.Fatal(errTmp)
	}

	runGitCommand([]string{"clone", "-q", repoName, tmpDir})

	master := regexp.MustCompile("(origin/HEAD|origin/master)")
	errChdir := os.Chdir(tmpDir)
	if errChdir != nil {
		log.Fatalf("Can't change folder to %s", tmpDir)
	}

	branches := make([]Branch, 0)
	merged := runGitCommand([]string{"branch", "-r", "--merged"})
	for _, branchName := range strings.Split(merged, "\n") {
		if master.MatchString(strings.TrimSpace(branchName)) { continue }
		if len(branchName) == 0 { continue }
		branches = append(branches, getBranchData(repoName, strings.TrimSpace(branchName), true, true))
	}

	notMerged := runGitCommand([]string{"branch", "-r", "--no-merged"})
	for _, branchName := range strings.Split(notMerged, "\n") {
		if len(branchName) == 0 { continue }
		logLastMonth := runGitCommand([]string{"log", "-1", "--since='1 month ago'", "-s", strings.TrimSpace(branchName), "--oneline"})
		isOutdated := len(logLastMonth) == 0
		branches = append(branches, getBranchData(repoName, strings.TrimSpace(branchName), false, isOutdated))
	}

	defer os.RemoveAll(tmpDir)

	return branches
}
