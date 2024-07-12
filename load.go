package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
	"github.com/shadiestgoat/log"
	"gopkg.in/yaml.v3"
	"whotfislucy.com/encryption"
	"whotfislucy.com/parser"
)

type opts struct {
	FinaleLoc  []string `short:"f" long:"finale"   default:"finale"        description:"The location of the finale pages"`
	EnvLoc     []string `short:"e" long:"env"                              description:"The .env file"`
	SectionLoc string   `short:"s" long:"sections" default:"sections.yaml" description:"The sections file"`
}

func Load(conf *opts) *parser.SectionState {
	fParser := flags.NewParser(conf, flags.HelpFlag|flags.AllowBoolValues|flags.IgnoreUnknown|flags.PassDoubleDash)

	_, err := fParser.Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok && err.Type == flags.ErrHelp {
			log.PrintDebug(err.Error())
			os.Exit(0)
		} else {
			log.FatalIfErr(err, "parsing args")
		}
	}

	for _, env := range conf.EnvLoc {
		log.ErrorIfErr(godotenv.Load(env), "reading env file '%v'", env)
	}

	secret := os.Getenv("AUTH_SECRET_KEY")
	if secret == "" {
		log.Fatal("Env var 'AUTH_SECRET_KEY' is non existent, but is needed!")
	}

	log.FatalIfErr(encryption.Init(secret), "initializing the aes cipher")

	sd, err := os.ReadFile(conf.SectionLoc)
	log.FatalIfErr(err, "reading sections file '%v'", conf.SectionLoc)

	rawSections := []*parser.RawSection{}
	log.FatalIfErr(yaml.Unmarshal(sd, &rawSections), "parsing raw sections")
	sectionState := parser.Parse(rawSections)

	for _, finaleDir := range conf.FinaleLoc {
		dir, err := os.ReadDir(finaleDir)
		log.FatalIfErr(err, "reading dir finale '%v'", dir)

		for _, f := range dir {
			n := f.Name()
			fullPath := filepath.Join(finaleDir, n)

			if f.IsDir() {
				log.Warn("'%v' is a directory, will not descend!", fullPath)
				continue
			} else if !strings.HasSuffix(n, ".md") {
				log.Warn("'%v' is not a .md file, skipping...", fullPath)
				continue
			}

			mdData, err := os.ReadFile(fullPath)
			log.FatalIfErr(err, "reading file '%v' as a finale", fullPath)

			parser.ParseFinale(n[:len(n)-3], string(mdData))
		}
	}

	return sectionState
}
