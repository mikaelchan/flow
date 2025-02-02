package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/mikaelchan/hamster/internal/domain/library"
	"github.com/mikaelchan/hamster/internal/domain/shared"
	"github.com/mikaelchan/hamster/pkg/util"
	"github.com/spf13/cobra"
)

// hamster lib create xxx --path /path/to/library --type movie --quality 1080p;5.2 --naming-template "{{.Title}} ({{.Year}}) {{.Quality}}"
// hamster lib scan xxx
// hamster lib archive xxx
// hamster lib delete xxx
// hamster lib rename xxx --name new_name
// hamster lib set-naming-template xxx --template new_template
// hamster lib set-quality-preference xxx --quality new_quality
// hamster lib set-location xxx --location new_location

var libraryCmd = &cobra.Command{
	Use:     "lib",
	Aliases: []string{"library"},
	Short:   "Library commands",
}

var libraryCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a library",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path, err := cmd.Flags().GetString("path")
		if err != nil {
			return err
		}
		mediaTypeStr, err := cmd.Flags().GetString("type")
		if err != nil {
			return err
		}
		mediaType, err := shared.FromString(mediaTypeStr)
		if err != nil {
			return err
		}
		quality, err := cmd.Flags().GetString("quality")
		if err != nil {
			return err
		}
		namingTemplate, err := cmd.Flags().GetString("naming-template")
		if err != nil {
			return err
		}

		createLibraryCmd := library.CreateLibrary{
			Name:              name,
			MediaType:         mediaType,
			Location:          shared.StorageLocation{Path: path},
			QualityPreference: shared.Match(mediaType, quality),
			NamingTemplate:    shared.NamingTemplate(namingTemplate),
		}
		if err := createLibraryCmd.Validate(); err != nil {
			return err
		}
		fmt.Printf("%s/commands/%s\n", cfg.Server.BaseURL, library.CreateLibraryContract)
		resp, err := util.PostJSON(fmt.Sprintf("%s/commands/%s", cfg.Server.BaseURL, library.CreateLibraryContract), createLibraryCmd, time.Duration(cfg.Server.Timeout)*time.Second)

		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
			return fmt.Errorf("failed to create library: %s", resp.Status)
		}
		fmt.Println("Library created successfully")
		return nil
	},
}

func init() {
	libraryCmd.PersistentFlags().StringP("path", "p", "", "The path to the library")
	libraryCmd.PersistentFlags().StringP("type", "t", "", "The type of the library")
	libraryCmd.PersistentFlags().StringP("quality", "q", "", "The quality of the library")
	libraryCmd.PersistentFlags().StringP("naming-template", "n", "", "The naming template of the library")
	libraryCmd.AddCommand(libraryCreateCmd)
}
