package repository

import (
	"database/sql"
	"errors"
)

type RepositoryChunkHashSqlite struct{}

func (rchs RepositoryChunkHashSqlite) Create(hash string, tx *sql.Tx) (int64, error) {
	var id int64
	dml := `INSERT INTO chunk_hashes (hash) VALUES (?) RETURNING id`
	return id, tx.QueryRow(dml, hash).Scan(&id)
}

// @WARNING: ALTERAÇÃO DE COMPORTAMENTO: ERRO SE NÃO ENCONTRAR O REGISTRO
func (rchs RepositoryChunkHashSqlite) GetIdByHash(hash string) (int64, error) {
	var id int64
	db, err := getConectionSqlite()
	if err != nil {
		return id, err
	}
	defer db.Close()

	dml := `SELECT id FROM chunk_hashes WHERE hash = ?`
	err = db.QueryRow(dml, hash).Scan(&id)
	if err == sql.ErrNoRows {
		return id, errors.New("record not found")
	}

	return id, err
}

// @WARNING: ALTERAÇÃO DE COMPORTAMENTO: ERRO SE NÃO ENCONTRAR O REGISTRO
func (rchs RepositoryChunkHashSqlite) GetHashById(id int64) (string, error) {
	var hash string
	db, err := getConectionSqlite()
	if err != nil {
		return hash, err
	}
	defer db.Close()

	dml := `SELECT hash FROM chunk_hashes WHERE id = ?`
	err = db.QueryRow(dml, id).Scan(&hash)
	if err == sql.ErrNoRows {
		return hash, errors.New("record not found")
	}

	return hash, err
}

func (rchs RepositoryChunkHashSqlite) RemoveAllWithTransaction(tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM chunk_hashes`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM sqlite_sequence WHERE name='chunk_hashes'`)
	return err
}

func (rchs RepositoryChunkHashSqlite) RemoveByIdsWithTransaction(ids []int64, tx *sql.Tx) error {
	dml := `DELETE FROM chunk_hashes WHERE id = ?`
	stmt, err := tx.Prepare(dml)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, id := range ids {
		result, err := stmt.Exec(id)
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
	}

	return nil
}

func (rchs RepositoryChunkHashSqlite) RemoveByIdWithTransaction(id int64, tx *sql.Tx) error {
	result, err := tx.Exec(`DELETE FROM chunk_hashes WHERE id = ?`, id)
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
