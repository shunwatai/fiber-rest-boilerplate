package config

func (c *Config) setLocalDefault() {
	/* Set server default */
	c.Vpr.SetDefault("server", map[string]string{
		"env":  "local",
		"port": "7000",
	})
	/* Set database default */
	c.Vpr.SetDefault("database", map[string]any{
		"engine": "sqlite",
		"sqlite": map[string]string{
			"user":     "user",
			"password": "user",
			"database": "fiber-starter",
		},
	})
	/* Set log default */
	c.Vpr.SetDefault("logging", map[string]any{
		"level": 0,
		"type":  []string{"database"},
		"zap": map[string]any{
			"output":   []string{"file"},
			"filename": "requests.log",
		},
		"debugSymbol": nil,
	})
	/* Set rbmq default */
	c.Vpr.SetDefault("rbmq", map[string]any{
		"host": "localhost",
		"port": "5672",
		"user": "user",
		"pass": "password",
		"queues": map[string]string{
			"logQueue":   "log_queue",
			"emailQueue": "email_queue",
			"testQueue":  "test_queue",
		},
	})
}

func (c *Config) setDockerDefault() {
	/* Set server default */
	c.Vpr.SetDefault("server", map[string]string{
		"env":  "local",
		"port": "7000",
	})
	/* Set database default */
	c.Vpr.SetDefault("database", map[string]any{
		"engine": "sqlite",
		"sqlite": map[string]string{
			"user":     "user",
			"password": "user",
			"database": "fiber-starter",
		},
	})
	/* Set log default */
	c.Vpr.SetDefault("logging", map[string]any{
		"level": 0,
		"type":  []string{"database"},
		"zap": map[string]any{
			"output":   []string{"file"},
			"filename": "requests.log",
		},
		"debugSymbol": nil,
	})
	/* Set rbmq default */
	c.Vpr.SetDefault("rbmq", map[string]any{
		"host": "rabbitmq-dev",
		"port": "5672",
		"user": "user",
		"pass": "password",
		"queues": map[string]string{
			"logQueue":   "log_queue",
			"emailQueue": "email_queue",
			"testQueue":  "test_queue",
		},
	})
}

func (c *Config) setK3sDefault() {
	/* Set server default */
	c.Vpr.SetDefault("server", map[string]string{
		"env":  "local",
		"port": "7000",
	})
	/* Set database default */
	c.Vpr.SetDefault("database", map[string]any{
		"engine": "postgres",
		"sqlite": map[string]string{
			"user":     "user",
			"password": "password",
			"database": "postgres.database.svc.cluster.local",
		},
	})
	/* Set log default */
	c.Vpr.SetDefault("logging", map[string]any{
		"level": 0,
		"type":  []string{"database"},
		"zap": map[string]any{
			"output":   []string{"file"},
			"filename": "requests.log",
		},
		"debugSymbol": nil,
	})
	/* Set rbmq default */
	c.Vpr.SetDefault("rbmq", map[string]any{
		"host": "rbmq.database.svc.cluster.local",
		"port": "5672",
		"user": "user",
		"pass": "password",
		"queues": map[string]string{
			"logQueue":   "log_queue",
			"emailQueue": "email_queue",
			"testQueue":  "test_queue",
		},
	})
}
