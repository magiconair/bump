// Copyright 2026 The Bump Authors. All rights reserved. See LICENSE.

package main

import (
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/magiconair/bump/git"

	"github.com/urfave/cli"
)

var version = buildVersion()

func buildVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok || info.Main.Version == "" || info.Main.Version == "(devel)" {
		return "devel"
	}
	return info.Main.Version
}

func tagAndPush(v git.Version, push bool) error {
	if err := git.Tag(v); err != nil {
		return err
	}
	log.Print(v)
	if push {
		return git.PushTag(v)
	}
	return nil
}

func main() {
	log.SetFlags(0)

	empty, err := git.IsEmptyRepository()
	if err != nil {
		log.Fatal(err)
	}
	if empty {
		log.Fatal("git repository is empty. Please create at least one commit")
	}

	app := cli.NewApp()
	app.HideVersion = true
	app.HideHelp = true
	app.Usage = "A tool for managing versions in git tags"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "s, service",
			Usage: "service prefix for tags (e.g. foo for foo/v1.2.3)",
		},
	}

	var versions []git.Version
	var cur git.Version
	app.Before = func(c *cli.Context) error {
		service := c.GlobalString("s")
		var err error
		versions, err = git.Tags(service)
		if err != nil {
			return err
		}
		if len(versions) == 0 {
			versions = append(versions, git.Version{Service: service, Prefix: "v"})
		}
		cur = versions[len(versions)-1]
		return nil
	}

	app.Commands = []cli.Command{
		{
			Name:  "cur",
			Usage: "print current version",
			Action: func(c *cli.Context) error {
				fmt.Println(cur)
				return nil
			},
		},
		{
			Name:  "list",
			Usage: "print all version",
			Action: func(c *cli.Context) error {
				for _, v := range versions {
					fmt.Println(v)
				}
				return nil
			},
		},
		{
			Name:  "next",
			Usage: "print next version",
			Action: func(c *cli.Context) error {
				fmt.Println(cur.Bump())
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "major",
					Usage: "print next major version",
					Action: func(c *cli.Context) error {
						fmt.Println(cur.BumpMajor())
						return nil
					},
				},
				{
					Name:  "minor",
					Usage: "print next minor version",
					Action: func(c *cli.Context) error {
						fmt.Println(cur.BumpMinor())
						return nil
					},
				},
				{
					Name:  "patch",
					Usage: "print next patch version",
					Action: func(c *cli.Context) error {
						fmt.Println(cur.BumpPatch())
						return nil
					},
				},
			},
		},
		{
			Name:  "tag",
			Usage: "tag with next major/minor/patch version",
			Subcommands: []cli.Command{
				{
					Name:  "major",
					Usage: "tag with next major version",
					Flags: []cli.Flag{cli.BoolFlag{Name: "push", Usage: "push tag to origin"}},
					Action: func(c *cli.Context) error {
						return tagAndPush(cur.BumpMajor(), c.Bool("push"))
					},
				},
				{
					Name:  "minor",
					Usage: "tag with next minor version",
					Flags: []cli.Flag{cli.BoolFlag{Name: "push", Usage: "push tag to origin"}},
					Action: func(c *cli.Context) error {
						return tagAndPush(cur.BumpMinor(), c.Bool("push"))
					},
				},
				{
					Name:  "patch",
					Usage: "tag with next patch version",
					Flags: []cli.Flag{cli.BoolFlag{Name: "push", Usage: "push tag to origin"}},
					Action: func(c *cli.Context) error {
						return tagAndPush(cur.BumpPatch(), c.Bool("push"))
					},
				},
			},
		},
		{
			Name:  "version",
			Usage: "print bump version",
			Action: func(c *cli.Context) error {
				fmt.Println(version)
				return nil
			},
		},
	}

	if err = app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
