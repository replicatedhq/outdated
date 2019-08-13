package cli

import (
	"fmt"

	"github.com/replicatedhq/outdated/pkg/outdated"
)

const (
	maxImageLength = 50
	maxTagLength   = 50
)

func headerLine(images []outdated.RunningImage) (string, int, int) {
	longestImageNameLength := 0
	longestTagNameLength := 0
	for _, image := range images {
		repo, img, tag, err := outdated.ParseImageName(image.Image)
		if err != nil {
			return "", 0, 0
		}

		imageName := fmt.Sprintf("%s/%s", repo, img)

		if len(imageName) > longestImageNameLength {
			longestImageNameLength = len(imageName)
		}
		if len(tag) > longestTagNameLength {
			longestTagNameLength = len(tag)
		}
	}

	imageColumnWidth := longestImageNameLength - len("Image")
	if imageColumnWidth < len("Image     ") {
		imageColumnWidth = len("Image    ")
	}
	tagColumnWidth := longestTagNameLength - len("Current")
	if tagColumnWidth < len("Current     ") {
		tagColumnWidth = len("Current     ")
	}

	if imageColumnWidth > maxImageLength {
		imageColumnWidth = maxImageLength
	}

	if tagColumnWidth > maxTagLength {
		tagColumnWidth = maxTagLength
	}

	imageColumnSpacer := ""
	for i := 0; i < imageColumnWidth; i++ {
		imageColumnSpacer = fmt.Sprintf("%s ", imageColumnSpacer)
	}
	tagColumnSpacer := ""
	for i := 0; i < tagColumnWidth; i++ {
		tagColumnSpacer = fmt.Sprintf("%s ", tagColumnSpacer)
	}

	headerLine := fmt.Sprintf("Image%sCurrent%sLatest%sBehind", imageColumnSpacer, tagColumnSpacer, tagColumnSpacer)

	return headerLine, imageColumnWidth, tagColumnWidth
}

func runningImage(image outdated.RunningImage, imageColumnWidth int, tagColumnWidth int) string {
	repo, img, tag, err := outdated.ParseImageName(image.Image)
	if err != nil {
		return ""
	}

	imageName := fmt.Sprintf("%s/%s", repo, img)
	truncatedImageName := imageName
	if len(truncatedImageName) > maxImageLength {
		truncatedImageName = fmt.Sprintf("%s...", truncatedImageName[0:maxImageLength-3])
	}

	truncatedTagName := tag
	if len(tag) > maxTagLength {
		truncatedTagName = fmt.Sprintf("%s...", truncatedTagName[0:maxTagLength-3])
	}

	imageColumnSpacer := ""
	for i := len(truncatedImageName); i < imageColumnWidth+5; i++ {
		imageColumnSpacer = fmt.Sprintf("%s ", imageColumnSpacer)
	}
	tagColumnSpacer := ""
	for i := 0; i < len(truncatedTagName); i++ {
		tagColumnSpacer = fmt.Sprintf("%s ", tagColumnSpacer)
	}

	return fmt.Sprintf("%s%s%s", truncatedImageName, imageColumnSpacer, truncatedTagName)
}

func completedImage(image outdated.RunningImage, checkResult *outdated.CheckResult, imageColumnWidth int, tagColumnWidth int) string {
	repo, img, tag, err := outdated.ParseImageName(image.Image)
	if err != nil {
		return ""
	}

	imageName := fmt.Sprintf("%s/%s", repo, img)
	truncatedImageName := imageName
	if len(truncatedImageName) > maxImageLength {
		truncatedImageName = fmt.Sprintf("%s...", truncatedImageName[0:maxImageLength-3])
	}

	truncatedTagName := tag
	if len(tag) > maxTagLength {
		truncatedTagName = fmt.Sprintf("%s...", truncatedTagName[0:maxTagLength-3])
	}

	imageColumnSpacer := ""
	for i := len(truncatedImageName); i < imageColumnWidth+5; i++ {
		imageColumnSpacer = fmt.Sprintf("%s ", imageColumnSpacer)
	}
	tagColumnSpacer := ""
	for i := len(truncatedTagName); i < tagColumnWidth+7; i++ {
		tagColumnSpacer = fmt.Sprintf("%s ", tagColumnSpacer)
	}

	truncatedLatestTagName := checkResult.LatestVersion
	if len(truncatedLatestTagName) > maxTagLength {
		truncatedLatestTagName = fmt.Sprintf("%s...", truncatedLatestTagName[0:maxTagLength-3])
	}
	latestTagColumnSpacer := ""
	for i := len(truncatedLatestTagName); i < tagColumnWidth+6; i++ {
		latestTagColumnSpacer = fmt.Sprintf("%s ", latestTagColumnSpacer)
	}

	return fmt.Sprintf("%s%s%s%s%s%s%d", truncatedImageName, imageColumnSpacer, truncatedTagName, tagColumnSpacer, truncatedLatestTagName, latestTagColumnSpacer, checkResult.VersionsBehind)

}

func erroredImage(image outdated.RunningImage, checkResult *outdated.CheckResult, imageColumnWidth int, tagColumnWidth int) string {
	repo, img, tag, err := outdated.ParseImageName(image.Image)
	if err != nil {
		return ""
	}

	imageName := fmt.Sprintf("%s/%s", repo, img)
	truncatedImageName := imageName
	if len(truncatedImageName) > maxImageLength {
		truncatedImageName = fmt.Sprintf("%s...", truncatedImageName[0:maxImageLength-3])
	}

	truncatedTagName := tag
	if len(tag) > maxTagLength {
		truncatedTagName = fmt.Sprintf("%s...", truncatedTagName[0:maxTagLength-3])
	}

	imageColumnSpacer := ""
	for i := len(truncatedImageName); i < imageColumnWidth+5; i++ {
		imageColumnSpacer = fmt.Sprintf("%s ", imageColumnSpacer)
	}
	tagColumnSpacer := ""
	for i := len(truncatedTagName); i < tagColumnWidth+7; i++ {
		tagColumnSpacer = fmt.Sprintf("%s ", tagColumnSpacer)
	}

	message := "Unable to get image data"
	if checkResult != nil {
		message = checkResult.CheckError
	}
	return fmt.Sprintf("%s%s%s%s%s", truncatedImageName, imageColumnSpacer, truncatedTagName, tagColumnSpacer, message)

}
