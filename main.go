package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "cbor-gen-for",
		Usage: "Generate CBOR encoders for types",
		Action: func(c *cli.Context) error {
			filename := os.Getenv("GOFILE")
			path, _ := os.Getwd()
			err := Generator{
				Filename:   filename,
				Path:       path,
				Package:    os.Getenv("GOPACKAGE"),
				GenStructs: c.Args().Slice(),
			}.GenerateCborTypes()
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
