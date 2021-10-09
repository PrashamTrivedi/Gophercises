module urlShortner

go 1.15

require (
	github.com/go-yaml/yaml v2.1.0+incompatible // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	pht/urlshortnerhandler v0.0.0-00010101000000-000000000000
)

replace pht/urlshortnerhandler => ./urlShortnerHandler
