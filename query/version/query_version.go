package queryversion

const (
	QueryGetVersion = `
		SELECT os,version FROM ss_version WHERE os = $1;
	`
)
