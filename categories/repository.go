package categories

import "database/sql"

func GetAllCategories(db *sql.DB) ([]Categories, error) {
	rows, err := db.Query("SELECT id, name FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []Categories

	for rows.Next() {
		var c Categories
		if err := rows.Scan(&c.ID, &c.Name); err != nil {
			return nil, err
		}
		list = append(list, c)
	}

	return list, nil
}
