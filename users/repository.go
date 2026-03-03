package users

import "database/sql"

// Exporté : Repository
type Repository struct {
	DB *sql.DB
}

// Exporté : méthode GetOnlineUsers
func (r *Repository) GetOnlineUsers() ([]string, error) {
	rows, err := r.DB.Query(`
        SELECT userName
        FROM users
        WHERE userOnline = 1
    `)
	if err != nil {
		return nil, err
	}

	var users []string
	for rows.Next() {
		var name string
		rows.Scan(&name)
		users = append(users, name)
	}

	return users, nil
}
