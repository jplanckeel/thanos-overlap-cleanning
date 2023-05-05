package cmd

import (
	"fmt"
	"os"

	"github.com/jplanckeel/thanos-overlap-cleanning/pkg"
	"github.com/spf13/cobra"
)

var dryrun bool
var accessKey string
var secretKey string
var bucketName string
var region string
var maxTime string
var minTime string
var labelsSelector string
var cacheDir string

func init() {
	rootCmd.PersistentFlags().BoolVar(&dryrun, "dryrun", false, "enable dry-run mode")
	rootCmd.PersistentFlags().StringVar(&accessKey, "access-key", "", "access key for bucket account")
	rootCmd.PersistentFlags().StringVar(&secretKey, "secret-key", "", "secret key for bucket account")
	rootCmd.PersistentFlags().StringVar(&bucketName, "bucket-name", "", "bucket name to check overlapping")
	rootCmd.PersistentFlags().StringVar(&region, "region", "", "region for bucket account")
	rootCmd.PersistentFlags().StringVar(&minTime, "min-time", "", "start of time range limit 9999-12-31T23:59:59Z")
	rootCmd.PersistentFlags().StringVar(&maxTime, "max-time", "", "end of time range limit 9999-12-31T23:59:59Z")
	rootCmd.PersistentFlags().StringVar(&labelsSelector, "labels-selector", "", "label selector to find overlapping")
	rootCmd.PersistentFlags().StringVarP(&cacheDir, "cache-dir", "", "./data", "cache dir to stock metadata")
}

var rootCmd = &cobra.Command{
	Use:   "toc",
	Short: "Thanos Overlaps Cleaning ",
	Long:  `A cli to cleaning Overlaps in Thanos S3 Bucket`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg.CheckOverlap(dryrun, accessKey, secretKey, bucketName, region, maxTime, minTime, labelsSelector, cacheDir)
		fmt.Println("script finish")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}