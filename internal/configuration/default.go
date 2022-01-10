package configuration

func DefaultConfig() *Config {
	return &Config{
		Profiles: []Profile{
			{
				ID: "default",
				Name: "default",
				// todo
			},
		},
		SamplingInterval: "60s",
		FruitingChambers: []FruitingChamber{
			{
				ID: "default",
				Name: "Martha",
				Profile: "default",
			},
		},
	}
}