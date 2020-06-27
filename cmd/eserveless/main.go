package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func executeServelessCmd() {

}

func main() {
	cmd := &cobra.Command{
		Use:          "eserveless",
		Short:        "A serverless platform cli tool!",
		SilenceUsage: true,
	}

	deployCmd := &cobra.Command{
		Use: "deploy",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 || !strings.Contains(args[0], "/") {
				return errors.New("Please specify the github repo: user/repo")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			repo := fmt.Sprintf("https://github.com/%s", args[0])
			cmd.Println(fmt.Sprintf("Asking serveless to check in ... [%s]", repo))

			apiURL := "http://localhost:8000/api/projects"
			jsonBody := []byte(fmt.Sprintf(`{"repository":"%s"}`, repo))
			req, _ := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return errors.New("error connecting to eserveless platform server")
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return errors.New("something went wrong with your request")
			}
			body, _ := ioutil.ReadAll(resp.Body)
			var data map[string]interface{}
			err = json.Unmarshal(body, &data)
			if err != nil {
				return errors.New("something went wrong with your request")
			}

			cmd.Println("[ok] ", data["message"])
			if data["functions"] == nil {
				return nil
			}
			httpFns := data["functions"].([]interface{})
			if len(httpFns) != 0 {
				cmd.Println("List of http api endpoints")
				for _, v := range httpFns {
					cmd.Println("- ", fmt.Sprintf("http://localhost:8000%s", v.(string)))
				}
			}
			return nil
		},
	}

	cmd.AddCommand(deployCmd)
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
