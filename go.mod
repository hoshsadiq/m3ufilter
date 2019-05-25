module github.com/hoshsadiq/m3ufilter

require (
	github.com/grafov/m3u8 v0.6.1
	github.com/kr/pretty v0.1.0 // indirect
	github.com/maja42/goval v1.0.0
	github.com/maja42/no-comment v0.0.0-20180113082502-512948848672
	github.com/mitchellh/go-homedir v1.1.0
	github.com/sirupsen/logrus v1.4.2
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.2
)

replace (
	github.com/grafov/m3u8 => github.com/hoshsadiq/m3u8 v0.0.0-20190514185311-cb08e59df8fe
	github.com/maja42/goval => github.com/hoshsadiq/goval v1.0.1-0.20190525223338-f1ea9f026acd
)
