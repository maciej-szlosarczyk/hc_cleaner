package cmd

import (
	"fmt"
	"icejam/hc_cleaner/honeycomb"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "hc_cleaner",
	Short: "A CLI to clean honeycomb datasets from stale columns.",
	Long:  `A CLI to clean honeycomb datasets from stale columns.`,
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number.",
	Long:  `All software has versions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HC cleaner v0.1.0")
	},
}

var inactiveCmd = &cobra.Command{
	Use:     "inactive",
	Short:   "Delete inactive columns.",
	Long:    `Delete columns that have not received data since 'days'.`,
	Example: "hc_cleaner inactive my_dataset --since=10 --api-key=superSecret",
	RunE: func(cmd *cobra.Command, args []string) error {
		dataset := args[0]
		since, err := cmd.Flags().GetInt64("since")

		if err != nil {
			return err
		}

		apiKey, err := cmd.Flags().GetString("api-key")

		if err != nil {
			return err
		}

		result := honeycomb.DeleteInactiveColumns(dataset, apiKey, since)
		return result

	},
	Args: cobra.ExactArgs(1),
}

var prefixCmd = &cobra.Command{
	Use:     "prefix",
	Short:   "Delete columns with specified name prefix.",
	Long:    `Delete columns with name that starts with 'prefix'.`,
	Example: "hc_cleaner prefix my_dataset --prefix=http.query_params --api-key=superSecret",
	RunE: func(cmd *cobra.Command, args []string) error {
		dataset := args[0]
		prefix, err := cmd.Flags().GetString("prefix")

		if err != nil {
			return err
		}

		apiKey, err := cmd.Flags().GetString("api-key")

		if err != nil {
			return err
		}

		result := honeycomb.DeleteColumnsWithPrefix(dataset, apiKey, prefix)
		return result

	},
	Args: cobra.ExactArgs(1),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringP("api-key", "a", "", "Honeycomb API key")
	rootCmd.MarkFlagRequired("api-key")

	inactiveCmd.PersistentFlags().Int64P("since", "s", 90, "Delete columns that have not received data since as many days.")
	inactiveCmd.MarkFlagRequired("since")

	prefixCmd.PersistentFlags().StringP("prefix", "p", "", "Delete columns with a given prefix.")
	prefixCmd.MarkFlagRequired("prefix")

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.AddCommand(versionCmd, inactiveCmd, prefixCmd)
}
