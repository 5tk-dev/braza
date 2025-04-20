package tests

import (
	"os"
	"testing"

	"5tk.dev/braza"
)

func TestConfigDevSetupFromJson(t *testing.T) {
	c := &braza.Config{}
	c.SetupFromFile("src/cfgs/test_for_reader.json")

	if c.Env != "dev" {
		t.Errorf("config.Env got %q, want \"dev\"", c.Env)
	}
	if c.Servername != "localhost" {
		t.Errorf("config.Env got %q, want \"localhost\"", c.Env)
	}
	if c.SecretKey != "******" {
		t.Errorf("config.SecretKey got %q, want \"******\"", c.SecretKey)
	}
	if !c.ListeningInTLS {
		t.Errorf("config.ListenInTLS got %v, want true", c.ListeningInTLS)
	}
	if c.TemplateFolder != "src/front/html/" {
		t.Errorf("config.TemplateFolder got %q, want \"src/front/html\"", c.TemplateFolder)
	}
	if !c.DisableParseFormBody {
		t.Errorf("config.DisableParseFormBody got %v, want true", c.DisableFileWatcher)
	}
	if !c.DisableTemplateReloader {
		t.Errorf("config.DisableTemplateReloader %v, want true", c.DisableTemplateReloader)
	}
	if c.StaticFolder != "src/front/assets/" {
		t.Errorf("config.StaticFolder %q, want \"src/front/assets\"", c.StaticFolder)
	}
	if c.StaticUrlPath != "assets/" {
		t.Errorf("config.StaticUrlPath %q, want \"assets\"", c.StaticFolder)
	}
	if !c.DisableStatic {
		t.Errorf("config.DisableStatic %v, want true", c.DisableStatic)
	}
	if !c.Silent {
		t.Errorf("config.Silent %v, want true", c.Silent)
	}
	if c.LogFile != "src/logs/server.log" {
		t.Errorf("config.EnvFile %q, want \"src/logs/server.log\"", c.EnvFile)
	}
	if c.EnvFile != "src/.env" {
		t.Errorf("config.EnvFile %q, want \"src/.env\"", c.EnvFile)
	}
	if c.EnvFileTest != "src/.env.test" {
		t.Errorf("config.EnvFileTest %q, want \"src/.env.test\"", c.EnvFileTest)
	}
	if c.EnvFileProd != "src/.env.prod" {
		t.Errorf("config.EnvFileProd %q, want \"src/.env.prod\"", c.EnvFileTest)
	}
	if !c.DisableFileWatcher {
		t.Errorf("config.DisableFileWatcher %v, want true", c.DisableFileWatcher)
	}
}

func TestConfigDevSetupFromYml(t *testing.T) {
	c := &braza.Config{}
	e := c.SetupFromFile("src/cfgs/test_for_reader.yml")
	if e != nil {
		t.Error(e)
		return
	}

	if c.Env != "dev" {
		t.Errorf("config.Env got %q, want \"dev\"", c.Env)
	}
	if c.Servername != "localhost" {
		t.Errorf("config.Servername got %q, want \"localhost\"", c.Servername)
	}
	if c.SecretKey != "******" {
		t.Errorf("config.SecretKey got %q, want \"******\"", c.SecretKey)
	}
	if !c.ListeningInTLS {
		t.Errorf("config.ListenInTLS got %v, want true", c.ListeningInTLS)
	}
	if c.TemplateFolder != "src/front/html/" {
		t.Errorf("config.TemplateFolder got %q, want \"src/front/html\"", c.TemplateFolder)
	}
	if !c.DisableParseFormBody {
		t.Errorf("config.DisableParseFormBody got %v, want true", c.DisableFileWatcher)
	}
	if !c.DisableTemplateReloader {
		t.Errorf("config.DisableTemplateReloader %v, want true", c.DisableTemplateReloader)
	}
	if c.StaticFolder != "src/front/assets/" {
		t.Errorf("config.StaticFolder %q, want \"src/front/assets\"", c.StaticFolder)
	}
	if c.StaticUrlPath != "assets/" {
		t.Errorf("config.StaticUrlPath %q, want \"assets\"", c.StaticFolder)
	}
	if !c.DisableStatic {
		t.Errorf("config.DisableStatic %v, want true", c.DisableStatic)
	}
	if !c.Silent {
		t.Errorf("config.Silent %v, want true", c.Silent)
	}
	if c.LogFile != "src/logs/server.log" {
		t.Errorf("config.EnvFile %q, want \"src/logs/server.log\"", c.EnvFile)
	}
	if c.EnvFile != "src/.env" {
		t.Errorf("config.EnvFile %q, want \"src/.env\"", c.EnvFile)
	}
	if c.EnvFileTest != "src/.env.test" {
		t.Errorf("config.EnvFileTest %q, want \"src/.env.test\"", c.EnvFileTest)
	}
	if c.EnvFileProd != "src/.env.prod" {
		t.Errorf("config.EnvFileProd %q, want \"src/.env.prod\"", c.EnvFileTest)
	}
	if !c.DisableFileWatcher {
		t.Errorf("config.DisableFileWatcher %v, want true", c.DisableFileWatcher)
	}
}

func TestConfigTest1_envfile_dev(t *testing.T) {
	cfg := &braza.Config{}
	cfg.SetupFromFile("src/cfgs/config_test1.yml")
	app := braza.NewApp(cfg)

	app.Env = "" // "", d, dev, development
	app.Build()
	if v := os.Getenv("DATABASE_URI"); v != "dev1" {
		t.Errorf("got %q, wants dev1", v)
	}
}

func TestConfigTest1_envfile_test(t *testing.T) {
	cfg := &braza.Config{}
	cfg.SetupFromFile("src/cfgs/config_test1.yml")
	app := braza.NewApp(cfg)

	app.Env = "t" // t, test, testing,
	app.Build()
	if v := os.Getenv("DATABASE_URI"); v != "test1" {
		t.Errorf("got %q, wants test1", v)
	}
}

func TestConfigTest1_envfile_prod(t *testing.T) {
	cfg := &braza.Config{}
	cfg.SetupFromFile("src/cfgs/config_test1.yml")
	app := braza.NewApp(cfg)

	app.Env = "p" // p,prod, production
	app.Build()
	if v := os.Getenv("DATABASE_URI"); v != "prod1" {
		t.Errorf("got %q, wants prod1", v)
	}
}

func TestConfigTest1_logfile(t *testing.T) {
	cfg := &braza.Config{}
	cfg.SetupFromFile("src/cfgs/config_test1.yml")
	app := braza.NewApp(cfg)

	app.Env = "d" // p,prod, production
	app.Build()

}
