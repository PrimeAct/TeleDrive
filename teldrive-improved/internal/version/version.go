package version

import "sync"

type Info struct {
	Version   string
	Commit    string
	BuildTime string
}

var (
	info Info
	once sync.Once
)

func Set(version, commit, buildTime string) {
	once.Do(func() {
		info = Info{
			Version:   version,
			Commit:    commit,
			BuildTime: buildTime,
		}
	})
}

func Get() Info {
	return info
}
