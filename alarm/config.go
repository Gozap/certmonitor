package alarm

type Config struct {
	Type    string   `yml:"type"`
	Targets []string `yml:"targets"`
}

func ExampleConfig() []*Config {
	return []*Config{
		{
			Type: "smtp",
			Targets: []string{
				"mritd1234@gmail.com",
			},
		},
		{
			Type: "webhook",
			Targets: []string{
				"https://google.com",
			},
		},
	}
}
