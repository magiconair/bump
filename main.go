package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/magiconair/bump/git"
	"github.com/urfave/cli"
)

var errNoVersion = errors.New("no version")

func main() {
	log.SetFlags(0)

	versions, err := git.Tags()
	if err != nil {
		log.Fatal(err)
	}
	if len(versions) == 0 {
		log.Fatal("no versions")
	}
	cur := versions[len(versions)-1]

	app := cli.NewApp()
	app.HideVersion = true
	app.HideHelp = true
	app.Usage = "A tool for managing versions in git tags"
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
			Usage: "tag with next version",
			Action: func(c *cli.Context) error {
				v := cur.Bump()
				git.Tag(v)
				log.Print(v)
				return nil
			},
			Subcommands: []cli.Command{
				{
					Name:  "major",
					Usage: "tag with next major version",
					Action: func(c *cli.Context) error {
						v := cur.BumpMajor()
						git.Tag(v)
						log.Print(v)
						return nil
					},
				},
				{
					Name:  "minor",
					Usage: "tag with next minor version",
					Action: func(c *cli.Context) error {
						v := cur.BumpMinor()
						git.Tag(v)
						log.Print(v)
						return nil
					},
				},
				{
					Name:  "patch",
					Usage: "tag with next patch version",
					Action: func(c *cli.Context) error {
						v := cur.BumpPatch()
						git.Tag(v)
						log.Print(v)
						return nil
					},
				},
			},
		},
	}

	if err = app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
