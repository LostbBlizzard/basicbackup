package settings

import (
	"errors"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

type Project struct {
	ProjectPath string
	ProjectType string
}
type Settings struct {
	Projects []Project
}

func NewSettings() *Settings {
	return &Settings{
		Projects: make([]Project, 0),
	}
}

func GetSettingsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".config", "basicbackup"), nil

}

func LuaScriptPath() (string, error) {
	dir, err := GetSettingsDir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, "main.lua"), nil
}
func AllLuaScriptPath() (string, error) {
	dir, err := GetSettingsDir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, "all-project.lua"), nil
}

func GetSettingsPath() (string, error) {
	dir, err := GetSettingsDir()
	if err != nil {
		return "", err
	}

	return path.Join(dir, "config.yaml"), nil
}

func Getsettings() (*Settings, error) {
	settingsfilepath, err := GetSettingsPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(settingsfilepath); errors.Is(err, os.ErrNotExist) {
		return NewSettings(), nil
	} else {
		filetext, err := os.ReadFile(settingsfilepath)
		if err != nil {
			return nil, err
		}

		settings := NewSettings()

		err = yaml.Unmarshal([]byte(filetext), settings)
		if err != nil {
			return nil, err
		}

		return settings, nil
	}
}

func Savesettings(settings *Settings) error {
	settingsfilepath, err := GetSettingsPath()
	if err != nil {
		return err
	}

	d, err := yaml.Marshal(&settings)
	if err != nil {
		return err
	}

	dirpath, err := GetSettingsDir()
	if err != nil {
		return err
	}

	if _, err := os.Stat(dirpath); errors.Is(err, os.ErrNotExist) {
		err = os.MkdirAll(dirpath, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(settingsfilepath, []byte(d), 0644)
	if err != nil {
		return err
	}
	return nil
}
