package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jygastaud/go-gitlab-client"
)

type Config struct {
	Host    string `json:"host"`
	ApiPath string `json:"api_path"`
	Token   string `json:"token"`
}

func main() {
	help := flag.Bool("help", false, "Show usage")

	file, e := ioutil.ReadFile("config.json")
	if e != nil {
		fmt.Printf("Config file error: %v\n", e)
		os.Exit(1)
	}

	var config Config
	json.Unmarshal(file, &config)
	fmt.Printf("Config: %+v\n", config)

	gitlab := gogitlab.NewGitlab(config.Host, config.ApiPath, config.Token)

	var method string
	flag.StringVar(&method, "m", "", "Specify method to retrieve projects infos, available methods:\n"+
		"  > -m users\n"+
		"  > -m users    			-search PATTERN\n"+
		"  > -m groups   			-search PATTERN\n"+
		"  > -m groups   			-ids GROUP_ID\n"+
		"  > -m team     			-ids GROUP_ID\n"+
		"  > -m new_member  	-ids GROUP_ID\n"+
		"  > -m sync_members 	-ids GROUP_ID -search PATTERN")

	var id, ids, query string
	flag.StringVar(&id, "id", "", "Specify repository id")
	flag.StringVar(&ids, "ids", "", "A list of id separated by comma")
	flag.StringVar(&query, "search", "", "Specify a pattern to search")

	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *help == true || method == "" {
		flag.Usage()
		return
	}

	startedAt := time.Now()
	defer func() {
		fmt.Printf("processed in %v\n", time.Now().Sub(startedAt))
	}()

	switch method {
	case "users":
		fmt.Println("Fetching users...")

		users, err := gitlab.Users(query, 1, 10000)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, user := range users {
			fmt.Printf("> [%d] %s (%s) %s\n", user.Id, user.Username, user.Name, user.Email)
		}

	case "groups":

		if ids == "" {
			fmt.Println("Fetching groups...")

			groups, err := gitlab.Groups()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			for _, group := range groups {
				fmt.Printf("> [%d] %s (%s) %s\n", group.Id, group.Name, group.Path, group.Description)
			}
		} else {
			listIds := strings.Split(ids, ",")
			for _, id := range listIds {
				projects, err := gitlab.GroupProjects(id)

				if err != nil {
					fmt.Println(err.Error())
					return
				}

				for _, project := range projects {
					fmt.Printf("> [%d] %s\n", project.Id, project.Name)
				}
			}
		}

	case "team":
		fmt.Println("Fetching project team members…")

		if id == "" {
			flag.Usage()
			return
		}

		members, err := gitlab.GroupMembers(id)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, member := range members {
			fmt.Printf("> [%d] %s (%s)\n", member.Id, member.Username, member.Name)
		}

	case "new_member":
		fmt.Println("Add new member…")

		if id == "" {
			flag.Usage()
			return
		}

		// @todo: Need to be dynamic
		user_id := "3"
		access_level := "30"

		err := gitlab.AddGroupMember(id, user_id, access_level)

		// @todo: Handle errors
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("> [%s] added\n", id)

	case "delete_member":

		fmt.Println("Delete member…")

		if id == "" {
			flag.Usage()
			return
		}

		// @todo: Need to be dynamic
		user_id := "3"

		err := gitlab.RemoveGroupMember(id, user_id)

		// @todo: Handle errors
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Printf("> [%s] removed\n", id)

	case "sync_members":

		if ids == "" {
			flag.Usage()
			return
		}

		users, err := gitlab.Users(query, 1, 10000)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if ids != "" {
			for _, user := range users {
				user_id := strconv.Itoa(user.Id)
				access_level := "30"

				listIds := strings.Split(ids, ",")

				fmt.Printf("> Try to add [%d] %s (%s) %s to groups (%s)\n", user.Id, user.Username, user.Name, user.Email, ids)

				for _, project_id := range listIds {

					/*err := */ gitlab.AddGroupMember(project_id, user_id, access_level)

					// @todo: Handle errors
					//if err != nil {
					//	fmt.Println(err.Error())
					//}

					fmt.Printf("> User [%d] %s (%s) %s added to group id [%s]\n", user.Id, user.Username, user.Name, user.Email, project_id)
				}
			}
		}

	}
}
