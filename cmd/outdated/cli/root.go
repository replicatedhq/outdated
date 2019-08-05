package cli

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/replicatedhq/outdated/pkg/logger"
	"github.com/replicatedhq/outdated/pkg/outdated"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "outdated",
		Short:         "",
		Long:          `.`,
		SilenceErrors: true,
		SilenceUsage:  true,
		PreRun: func(cmd *cobra.Command, args []string) {
			viper.BindPFlags(cmd.Flags())
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			v := viper.GetViper()
			log := logger.NewLogger()
			log.Info("")

			o := outdated.Outdated{}

			log.Info("Finding images in cluster")
			images, err := o.ListImages(v.GetString("kubeconfig"))
			if err != nil {
				log.Error(err)
				log.Info("")
				os.Exit(1)
				return nil
			}

			log.Info("")
			head, imageColumnWidth, tagColumnWidth := headerLine(images)
			log.Header(head)

			for _, image := range images {
				log.StartImageLine(runningImage(image, imageColumnWidth, tagColumnWidth))
				checkResult, err := o.ParseImage(image.Image, image.PullableImage)
				if err != nil {
					log.Error(err)
					log.Info("")
					os.Exit(1)
					return nil
				}

				if checkResult.VersionsBehind != -1 {
					log.FinalizeImageLine(checkResult.VersionsBehind, completedImage(image, checkResult, imageColumnWidth, tagColumnWidth))
				} else {
					log.FinalizeImageLineWithError(erroredImage(image, checkResult, imageColumnWidth, tagColumnWidth))
				}
			}
			return nil
		},
	}

	cobra.OnInitialize(initConfig)

	cmd.Flags().String("kubeconfig", path.Join(homeDir(), ".kube", "config"), "path to the kubeconfig to use")

	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	return cmd
}

func InitAndExecute() {
	if err := RootCmd().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	viper.SetEnvPrefix("OUTDATED")
	viper.AutomaticEnv()
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
