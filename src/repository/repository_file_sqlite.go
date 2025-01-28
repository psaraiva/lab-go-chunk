package repository

import (
	"database/sql"
	"lab/src/model"

	_ "github.com/mattn/go-sqlite3"
)

type RepositoryFileSqlite struct{}

func (rfs RepositoryFileSqlite) Create(file model.File) (int64, error) {
	var id int64
	db, err := getConectionSqlite()
	if err != nil {
		return id, err
	}
	defer db.Close()

	query := `INSERT INTO files (name, hash) VALUES (?, ?) RETURNING id`
	err = db.QueryRow(query, file.Name, file.Hash).Scan(&id)
	return id, err
}

func (rfs RepositoryFileSqlite) GetHashByName(name string) (string, error) {
	var fileName string
	db, err := getConectionSqlite()
	if err != nil {
		return fileName, err
	}
	defer db.Close()

	query := `SELECT hash FROM files WHERE name = ?`
	err = db.QueryRow(query, name).Scan(&fileName)
	if err == sql.ErrNoRows {
		return fileName, nil
	}

	return fileName, err
}

func (rfs RepositoryFileSqlite) IsExistsByHashFile(hashFile string) (bool, error) {
	db, err := getConectionSqlite()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	query := `SELECT COUNT(id) FROM files WHERE hash = ?`
	err = db.QueryRow(query, hashFile).Scan(&count)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (rfs RepositoryFileSqlite) RemoveByHashFile(hashFile string) error {
	db, err := getConectionSqlite()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM files WHERE hash = ?`, hashFile)
	return err
}

func (rfs RepositoryFileSqlite) RemoveAll() error {
	db, err := getConectionSqlite()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM files`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`DELETE FROM sqlite_sequence WHERE name = 'files'`)
	return err
}

func (rfs RepositoryFileSqlite) GetIdByHashFile(hashFile string) (int64, error) {
	var id int64
	db, err := getConectionSqlite()
	if err != nil {
		return id, err
	}
	defer db.Close()

	query := `SELECT id FROM files WHERE hash = ?`
	err = db.QueryRow(query, hashFile).Scan(&id)
	if err == sql.ErrNoRows {
		return id, nil
	}

	return id, err
}
