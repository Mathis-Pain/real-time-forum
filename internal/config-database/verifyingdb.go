package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
)

// SELECT name FROM sqlite_master WHERE type='table' AND name='ma_table';
func InspectDbIntegrity(db *sql.DB) {

	fmt.Println("Cheking database...")

	// Lecture du fichier SQL de création des tables.
	fileCreatingTables, err := os.ReadFile("./internal/config-database/001_create_tables.sql")
	if err != nil {
		log.Fatalln("Error during reading File Create Table for checking tables integrity : %w", err)
	}

	// Extraction des noms des tables. Utilisation du Regex pour trouver dans le fichier sql.
	searchOccurenceInCreatingFile := regexp.MustCompile(`(?i)CREATE\s+TABLE\s+(?:IF\s+NOT\s+EXISTS\s+)?([a-zA-Z0-9_]+)`)
	matches := searchOccurenceInCreatingFile.FindAllStringSubmatch(string(fileCreatingTables), -1) // -1 c'est pour dire : donne moi tous les résultats correspondants.

	// Rangement dans un tableau de string.
	var tablesInFile []string
	for _, m := range matches {
		tablesInFile = append(tablesInFile, m[1]) // 1 permet de cibler le groupe de mot et non pas l'espace entre les groueps de mots.
	}

	// Lecture de la BDD.
	rows, err := db.Query("SELECT name FROM sqlite_master WHERE type='table'")
	if err != nil {
		log.Fatalln("Error during reading SQL Database for table integrity cheking : %w", err)
	}
	defer rows.Close()

	// Rangement dans un tableau de string.
	var tablesInDB []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		// Dès qu'il y a une table SQL avec AUTOINCREMENT, SQL créé une table
		// de son côté pour sa gestion interne. Dans cette vérification de nos tables,
		// nous ne la prenons pas en compte.
		if name != "sqlite_sequence" {
			tablesInDB = append(tablesInDB, name)
		}
	}

	// Vérification si il y a le même nombre de tables.
	if len(tablesInFile) != len(tablesInDB) {
		log.Fatalf(
			"They haven't the same amount of tables in : \n - Database (%d) \n  -Table Create File (%d)",
			len(tablesInFile),
			len(tablesInDB),
		)
	}

	// Vérification si les noms des tables sont :
	// -Dans le même ordre.
	// -Les mêmes.
	for i := 0; i < len(tablesInFile); i++ {
		if tablesInFile[i] != tablesInDB[i] {
			log.Fatalf(
				"Tables don't have the same name in : \n -Database (%s) \n -Table Create File (%s) ",
				tablesInFile[i],
				tablesInDB[i],
			)
		}
	}
	fmt.Println("Database integrity positive...")
}
