package chaincode

import (
	"github.com/spf13/cobra"
	"github.com/fsouza/go-dockerclient"
	"fmt"
)

var chaincodeDeleteCmd *cobra.Command

const deleteCmdName = "delete"

// deleteCmd returns the cobra command for Chaincode Delete
func deleteCmd() *cobra.Command {
	chaincodeDeleteCmd = &cobra.Command{
		Use:       deleteCmdName,
		Short:     "Delete chaincode.",
		Long:      "Delete all chaincode containers & images run by peers under this namespace.",
		ValidArgs: []string{"1"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return delete(cmd)
		},
	}
	flagList := []string{
		"namespace",
	}
	attachFlags(chaincodeDeleteCmd, flagList)

	return chaincodeDeleteCmd
}

func delete(cmd *cobra.Command) error {
	// delete all chaincode containers & images under this namespace
	endpoint := "unix:///host/var/run/docker.sock"
	client, err := docker.NewClient(endpoint)
	if err != nil {
		return fmt.Errorf("cannot connect to docker daemon when delete")
	}
	containers, err := client.ListContainers(docker.ListContainersOptions{All: true, Filters: map[string][]string{"name": {namespace}}})
	if err != nil {
		return fmt.Errorf("cannot list containers to be deleted when delete")
	}
	for _, container := range containers {
		client.RemoveContainer(docker.RemoveContainerOptions{ID: container.ID})
		client.RemoveImage(container.Image)
	}
	return nil
}