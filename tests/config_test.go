package tests

import (
	"testing"

	"5tk.dev/braza"
)

func TestCfgAppFromJson(t *testing.T) {
	app := braza.NewApp(
		braza.NewConfigFromFile("cfg_loadCfg.json"),
	)

	if app.Env != "test" {
		t.Errorf("app.Env got %v: - want 'test'", app.Env)
	}
	if !app.Silent {
		t.Errorf("app.Silent got %v: - want 'true'", app.Silent)
	}
	if app.LogFile != "test.log" {
		t.Errorf("app.LogFile got %v: - want 'test.log'", app.LogFile)
	}
	if app.SecretKey != "supersecret" {
		t.Errorf("app.SecretKey got %v: - want 'supersecret'", app.SecretKey)
	}
	if app.Servername != "localhost" {
		t.Errorf("app.Servername got %v: - want 'localhost'", app.Servername)
	}
	if app.StaticFolder != "front/static" {
		t.Errorf("app.StaticFolder got %v: - want 'front/static'", app.StaticFolder)
	}
	if app.StaticUrlPath != "/cdn" {
		t.Errorf("app.StaticUrlPath got %v: - want '/cdn'", app.StaticUrlPath)
	}
	if app.TemplateFolder != "front/html" {
		t.Errorf("app.TemplateFolder got %v: - want 'front/html'", app.TemplateFolder)
	}
	if app.SessionName != "foo" {
		t.Errorf("app.SessionName got %v: - want 'foo'", app.SessionName)
	}
	if !app.EnableDocFull {
		t.Errorf("app.EnableDocFull got %v: - want 'true'", app.EnableDocFull)
	}
	if !app.ListeningInTLS {
		t.Errorf("app.ListeningInTLS got %v: - want 'true'", app.ListeningInTLS)
	}
	if !app.DisableStatic {
		t.Errorf("app.DisableStatic got %v: - want 'true'", app.DisableStatic)
	}
	if !app.DisableWarnOn405 {
		t.Errorf("app.DisableWarnOn405 got %v: - want 'true'", app.DisableWarnOn405)
	}
	if !app.DisableFileWatcher {
		t.Errorf("app.DisableFileWatcher got %v: - want 'true'", app.DisableFileWatcher)
	}
	if !app.DisableParseFormBody {
		t.Errorf("app.DisableParseFormBody got %v: - want 'true'", app.DisableParseFormBody)
	}
	if !app.DisableTemplateReloader {
		t.Errorf("app.DisableTemplateReloader got %v: - want 'true'", app.DisableTemplateReloader)
	}
}

func TestCfgAppFromYaml(t *testing.T) {
	app := braza.NewApp(
		braza.NewConfigFromFile("cfg_loadCfg.yaml"),
	)

	if app.Env != "test" {
		t.Errorf("app.Env got %v: - want 'test'", app.Env)
	}
	if !app.Silent {
		t.Errorf("app.Silent got %v: - want 'true'", app.Silent)
	}
	if app.LogFile != "test.log" {
		t.Errorf("app.LogFile got %v: - want 'test.log'", app.LogFile)
	}
	if app.SecretKey != "supersecret" {
		t.Errorf("app.SecretKey got %v: - want 'supersecret'", app.SecretKey)
	}
	if app.Servername != "localhost" {
		t.Errorf("app.Servername got %v: - want 'localhost'", app.Servername)
	}
	if app.StaticFolder != "front/static" {
		t.Errorf("app.StaticFolder got %v: - want 'front/static'", app.StaticFolder)
	}
	if app.StaticUrlPath != "/cdn" {
		t.Errorf("app.StaticUrlPath got %v: - want '/cdn'", app.StaticUrlPath)
	}
	if app.TemplateFolder != "front/html" {
		t.Errorf("app.TemplateFolder got %v: - want 'front/html'", app.TemplateFolder)
	}
	if app.SessionName != "foo" {
		t.Errorf("app.SessionName got %v: - want 'foo'", app.SessionName)
	}
	if !app.EnableDocFull {
		t.Errorf("app.EnableDocFull got %v: - want 'true'", app.EnableDocFull)
	}
	if !app.ListeningInTLS {
		t.Errorf("app.ListeningInTLS got %v: - want 'true'", app.ListeningInTLS)
	}
	if !app.DisableStatic {
		t.Errorf("app.DisableStatic got %v: - want 'true'", app.DisableStatic)
	}
	if !app.DisableWarnOn405 {
		t.Errorf("app.DisableWarnOn405 got %v: - want 'true'", app.DisableWarnOn405)
	}
	if !app.DisableFileWatcher {
		t.Errorf("app.DisableFileWatcher got %v: - want 'true'", app.DisableFileWatcher)
	}
	if !app.DisableParseFormBody {
		t.Errorf("app.DisableParseFormBody got %v: - want 'true'", app.DisableParseFormBody)
	}
	if !app.DisableTemplateReloader {
		t.Errorf("app.DisableTemplateReloader got %v: - want 'true'", app.DisableTemplateReloader)
	}
}
