package repository

import (
	"database/sql"
	"errors"
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

// @WARNING: ALTEREÇÃO DE COMPORTAMENTO: ERRO SE NÃO ENCONTRAR O REGISTRO
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
		return fileName, errors.New("record not found")
	}

	return fileName, err
}

func (rfs RepositoryFileSqlite) IsExistsByHash(hash string) (bool, error) {
	db, err := getConectionSqlite()
	if err != nil {
		return false, err
	}
	defer db.Close()

	var count int
	query := `SELECT COUNT(id) FROM files WHERE hash = ?`
	err = db.QueryRow(query, hash).Scan(&count)
	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// @WARNING: ALTEREÇÃO DE COMPORTAMENTO: ERRO SE NÃO ENCONTRAR O REGISTRO
func (rfs RepositoryFileSqlite) RemoveByHash(hash string) error {
	db, err := getConectionSqlite()
	if err != nil {
		return err
	}
	defer db.Close()

	result, err := db.Exec(`DELETE FROM files WHERE hash = ?`, hash)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("record not found")
	}

	return err
}

func (rfs RepositoryFileSqlite) RemoveAll() error {
	err := rfs.resetTable()
	if err != nil {
		return err
	}

	err = rfs.resetAutoIncrement()
	if err != nil {
		return err
	}

	return nil
}

// @WARNING: ALTEREÇÃO DE COMPORTAMENTO: ERRO SE NÃO ENCONTRAR O REGISTRO
func (rfs RepositoryFileSqlite) GetIdByHash(hash string) (int64, error) {
	var id int64
	db, err := getConectionSqlite()
	if err != nil {
		return id, err
	}
	defer db.Close()

	query := `SELECT id FROM files WHERE hash = ?`
	err = db.QueryRow(query, hash).Scan(&id)
	if err == sql.ErrNoRows {
		return id, errors.New("record not found")
	}

	return id, err
}

func (rfs RepositoryFileSqlite) resetTable() error {
	db, err := getConectionSqlite()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM files`)
	if err != nil {
		return err
	}

	var count int
	query := `SELECT COUNT(id) FROM files`
	err = db.QueryRow(query).Scan(&count)
	if err == sql.ErrNoRows {
		return nil
	}

	return err
}

func (rfs RepositoryFileSqlite) resetAutoIncrement() error {
	db, err := getConectionSqlite()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(`DELETE FROM sqlite_sequence WHERE name = 'files'`)
	if err != nil {
		return err
	}

	var seq int
	query := `SELECT seq FROM sqlite_sequence WHERE name='files'`
	err = db.QueryRow(query).Scan(&seq)
	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		return err
	}

	return nil
}
