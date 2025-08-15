package db

func InitDb() (*DBConnections, error) {
	dbConnections, err := NewPostgres()
	if err != nil {
		return nil, err
	}

	return dbConnections, nil
}
