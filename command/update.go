package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sort"

	"github.com/Asice-Cloud/tzgin2/util"
	"github.com/urfave/cli/v2"
	"golang.org/x/mod/semver"
)

func Update(c *cli.Context) error {
	stop := make(chan int, 1)
	go util.Loading(stop)
	path, err := checkExists()
	if err != nil {
		log.Printf("checkExists error: %v", err)
		return cli.Exit(err.Error(), 1)
	}
	ver, err := getLatestVer()
	if err != nil {
		log.Printf("getLatestVer error: %v", err)
		return cli.Exit(err.Error(), 1)
	}

	if *ver == "v"+c.App.Version {
		return cli.Exit("\nAlready the newest version", 2)
	}
	// beng, err!= nil has existed
	// if len(path) != 0 && err == nil {
	if len(path) != 0 {
		cmd := exec.Command(path, "install", fmt.Sprintf("github.com/xjtu-tenzor/tz-gin@%s", *ver))

		_, err := cmd.StdoutPipe()
		if err != nil {
			log.Printf("StdoutPipe error: %v", err)
			return cli.Exit(err.Error(), 1)
		}
		_, err = cmd.StderrPipe()
		if err != nil {
			log.Printf("StderrPipe error: %v", err)
			return cli.Exit(err.Error(), 1)
		}

		err = cmd.Start()
		if err != nil {
			log.Printf("cmd.Start error: %v", err)
			return cli.Exit(err.Error(), 1)
		}

		err = cmd.Wait()
		if err != nil {
			log.Printf("cmd.Wait error: %v", err)
			return cli.Exit(err.Error(), 1)
		}

		stop <- 1
		util.SuccessMsg(fmt.Sprintf("\nSuccessfully update to %s\n", *ver))
		return nil
	}

	return err
}

func getLatestVer() (*string, error) {
	apiUrl := "https://api.github.com/repos/Asice-Cloud/tzgin2/tags"
	response, err := http.Get(apiUrl)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch tags")
	}

	var releaseInfo []struct {
		Name string `json:"name"`
	}

	err = json.NewDecoder(response.Body).Decode(&releaseInfo)

	if err != nil {
		return nil, err
	}

	sort.Slice(releaseInfo, func(i, j int) bool {
		return semver.Compare(releaseInfo[i].Name, releaseInfo[j].Name) == 1
	})

	return &releaseInfo[0].Name, nil
}

func checkExists() (string, error) {
	path, err := exec.LookPath("go")
	if err != nil {
		fmt.Printf("cannot find command\"go\"")
		return "", err
	} else {
		// fmt.Printf("\"go\" executable is in '%s'\n", path)
		return path, nil
	}
}
