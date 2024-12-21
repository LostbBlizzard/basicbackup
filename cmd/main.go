package main

import (
	"fmt"
	"os"
	"path"

	"github.com/alecthomas/kong"
	lua "github.com/yuin/gopher-lua"

	settingsmod "backup/internal/settings"
)

var CLI struct {
	Version struct {
		Info string `default:"all" enum:"all,name,version,githash,targetos,targetarch,builddate" help:""`
	} ` short:"v" cmd:""  help:"get program version" `

	Config struct {
		Get struct {
			ConfigPath struct{} `cmd:"" get the config path`
		} `cmd`
		Add struct {
			Project struct {
				Path string `required:"" short:"p"`
			} `cmd:""  help:"Set config value of AutoUpdate" `
		} `cmd`
		Ls struct {
			Project struct {
			} `cmd:""  help:"Set config value of AutoUpdate" `
		} `cmd`
	} `cmd:"" short:"c"`

	Backup struct {
		All struct {
		} `cmd:""  help:"Set config value of AutoUpdate" `
		Project struct {
			Value string `required:"" short:"v"`
		} `cmd:""  help:"Set config value of AutoUpdate" `
		Group struct {
			Value string `required:"" short:"v"`
		} `cmd:""  help:"Set config value of AutoUpdate" `
	} `cmd:"" short:"c"`
}

func backupproject(project settingsmod.Project, output string) {
	L := lua.NewState()
	defer L.Close()

	os.MkdirAll(output, os.ModePerm)
	// log := log.New(os.Stderr, "", 0)
	txt, err := settingsmod.LuaScriptPath()
	if err != nil {
		panic(err)
	}

	if err := L.DoFile(txt); err != nil {
		panic(err)
	}
	var callback = L.ToFunction(-1)
	L.Pop(1)

	L.Push(callback)

	L.Push(lua.LString(project.ProjectPath))

	settings := L.NewTable()
	{
		L.SetField(settings, "type", lua.LString(project.ProjectType))
	}
	L.Push(settings)
	L.Push(lua.LString(output))

	err = L.PCall(3, 0, nil)
	if err != nil {
		panic(err)
	}
}

func postbackup(_ []settingsmod.Project, input string, output string) {
	L := lua.NewState()
	defer L.Close()

	txt, err := settingsmod.AllLuaScriptPath()
	if err != nil {
		panic(err)
	}

	if err := L.DoFile(txt); err != nil {
		panic(err)
	}
	var callback = L.ToFunction(-1)
	L.Pop(1)

	L.Push(callback)

	L.Push(lua.LString(input))
	L.Push(lua.LString(output))

	err = L.PCall(2, 0, nil)
	if err != nil {
		panic(err)
	}

}
func main() {

	ctx := kong.Parse(&CLI)

	switch ctx.Command() {
	case "version":
		// type versionfile struct {
		// 	Name       string
		// 	Version    string
		// 	Githash    string
		// 	Targetos   string
		// 	Targetarch string
		// 	Builddate  string
		// }
		// versioninfo := versionfile{}
		//
		// err := yaml.Unmarshal([]byte(VersionFile), &versioninfo)
		// if err != nil {
		// 	panic(err)
		// }
		//
		// infotolog := CLI.Version.Info
		// switch infotolog {
		// case "all":
		// 	fmt.Printf("name:%s \n", versioninfo.Name)
		// 	fmt.Printf("version:%s \n", versioninfo.Version)
		// 	fmt.Printf("githash:%s \n", versioninfo.Githash)
		// 	fmt.Printf("targetos:%s \n", versioninfo.Targetos)
		// 	fmt.Printf("targetarch:%s \n", versioninfo.Targetarch)
		// 	fmt.Printf("builddate:%s \n", versioninfo.Builddate)
		//
		// case "name":
		// 	println(versioninfo.Name)
		// case "version":
		// 	println(versioninfo.Version)
		// case "githash":
		// 	println(versioninfo.Githash)
		// case "targetos":
		// 	println(versioninfo.Targetos)
		// case "targetarch":
		// 	println(versioninfo.Targetarch)
		// case "builddate":
		// 	println(versioninfo.Builddate)
		// }

	case "backup all":
		config, err := settingsmod.Getsettings()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		mydir, err := os.Getwd()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		mydirbackupdir := path.Join(mydir, "backup")

		os.RemoveAll(mydirbackupdir)

		for _, item := range config.Projects {
			newbackup := path.Join(mydirbackupdir, path.Base(item.ProjectPath))
			backupproject(item, newbackup)
		}
		postbackup(config.Projects, mydirbackupdir, "backupfile")

	case "backup project":
		var project = CLI.Backup.Project.Value

		config, err := settingsmod.Getsettings()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		mydir, err := os.Getwd()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		mydirbackupdir := path.Join(mydir, "backup")

		for _, item := range config.Projects {
			if item.ProjectPath == project {
				newbackup := path.Join(mydirbackupdir, path.Base(item.ProjectPath))
				backupproject(item, newbackup)

			}
		}

	case "backup group":
		var _ = CLI.Backup.Group.Value

		config, err := settingsmod.Getsettings()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		mydir, err := os.Getwd()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		mydirbackupdir := path.Join(mydir, "backup")

		for _, item := range config.Projects {
			newbackup := path.Join(mydirbackupdir, path.Base(item.ProjectPath))
			backupproject(item, newbackup)
		}

	case "config get config-path":
		path, err := settingsmod.GetSettingsPath()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}
		fmt.Println(path)

	case "config add project":

		var projectpath = CLI.Config.Add.Project.Path

		config, err := settingsmod.Getsettings()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		var newproject = settingsmod.Project{
			ProjectPath: projectpath,
			ProjectType: "git",
		}
		config.Projects = append(config.Projects, newproject)

		err = settingsmod.Savesettings(config)
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

	case "config ls project":

		config, err := settingsmod.Getsettings()
		if err != nil {
			fmt.Print(err)
			os.Exit(1)
		}

		for _, item := range config.Projects {
			println(item.ProjectPath)
		}

	default:
		panic(ctx.Command())
	}
}
