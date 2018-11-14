// Package config contrain all informations concerning configuration.
// Parameters can be accessed after using iniflags.Parse()
// Parameters can be set using command line or ini file.
// See -help or https://github.com/vharitonsky/iniflags/blob/master/README.md for details.
package config

import (
	"flag"
)

var (
	Port     = flag.String("port", "8080", "Port d'écoute du serveur")
	DbFile   = flag.String("database", "database.db", "Fichier de base SQLite")
	Firstrun = flag.Bool("init", false, "Création admin et 1er voucher")
)
