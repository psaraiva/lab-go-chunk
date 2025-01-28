package repository

import (
	"database/sql"
	"fmt"
)

type RepositoryChunkHashSqlite struct{}

func (rchs RepositoryChunkHashSqlite) Create(hash string, tx *sql.Tx) (int64, error) {
	id, err := rchs.GetIdByHash(hash)
	if err != nil {
		return 0, fmt.Errorf("falha ao buscar id de chunk hash por hash: %v", err)
	}

	if id > 0 {
		return id, nil
	}

	return rchs.create(hash, tx)
}

func (rchs RepositoryChunkHashSqlite) create(hash string, tx *sql.Tx) (int64, error) {
	var id int64
	query := `INSERT INTO chunk_hashes (hash) VALUES (?) RETURNING id`
	return id, tx.QueryRow(query, hash).Scan(&id)
}

func (rchs RepositoryChunkHashSqlite) GetIdByHash(hash string) (int64, error) {
	var id int64
	db, err := getConectionSqlite()
	if err != nil {
		return id, err
	}
	defer db.Close()

	query := `SELECT id FROM chunk_hashes WHERE hash = ?`
	err = db.QueryRow(query, hash).Scan(&id)
	if err == sql.ErrNoRows {
		return id, nil
	}

	return id, err
}

func (rchs RepositoryChunkHashSqlite) GetHashById(id int64) (string, error) {
	var hash string
	db, err := getConectionSqlite()
	if err != nil {
		return hash, err
	}
	defer db.Close()

	query := `SELECT hash FROM chunk_hashes WHERE id = ?`
	err = db.QueryRow(query, id).Scan(&hash)
	if err == sql.ErrNoRows {
		return hash, nil
	}

	return hash, err
}

func (rchs RepositoryChunkHashSqlite) RemoveAll(tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM chunk_hashes`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM sqlite_sequence WHERE name='chunk_hashes'`)
	return err
}

func (rchs RepositoryChunkHashSqlite) RemoveByChunkHashIds(ids []int64, tx *sql.Tx) error {
	query := `DELETE FROM chunk_hashes WHERE chunk_id = ?`
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, id := range ids {
		_, err = stmt.Exec(id)
		if err != nil {
			return err
		}
	}

	return nil
}

func (rchs RepositoryChunkHashSqlite) RemoveByChunkHashId(id int64, tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM chunk_hashes WHERE id = ?`, id)
	return err
}
