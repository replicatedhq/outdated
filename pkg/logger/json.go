package logger

type JSONResult struct {
	Repo  string `json:"repository"`
	Image string `json:"image"`
	Tag   string `json:"tag"`

	LatestVersion  string `json:"latestVersion"`
	VersionsBehind int64  `json:"versionsBehind"`

	Error *string `json:"error,omitempty"`
}
