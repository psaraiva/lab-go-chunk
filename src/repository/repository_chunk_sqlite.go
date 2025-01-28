package repository

import (
	"database/sql"
	"fmt"
	"lab/src/model"
)

type RepositoryChunkSqlite struct{}

func (rcs RepositoryChunkSqlite) Create(chunk model.Chunk) (int64, error) {
	var id int64
	db, err := getConectionSqlite()
	if err != nil {
		return id, err
	}
	defer db.Close()

	fileId, err := RepositoryFileSqlite{}.GetIdByHashFile(chunk.HashOriginalFile)
	if err != nil {
		return id, fmt.Errorf("falha ao buscar id do arquivo: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return id, err
	}
	defer tx.Rollback()

	id, err = rcs.create(chunk, fileId, tx)
	if err != nil {
		return id, fmt.Errorf("falha ao criar chunk: %v", err)
	}

	err = rcs.createChunkList(id, chunk.HashList, tx)
	if err != nil {
		return id, fmt.Errorf("falha ao criar lista de hash do chunk: %v", err)
	}

	return id, tx.Commit()
}

func (rcs RepositoryChunkSqlite) create(chunk model.Chunk, fileId int64, tx *sql.Tx) (int64, error) {
	var id int64
	query := `INSERT INTO chunks (size, file_id) VALUES (?, ?) RETURNING id`
	err := tx.QueryRow(query, chunk.Size, fileId).Scan(&id)
	return id, err
}

func (rcs RepositoryChunkSqlite) createChunkList(chunkid int64, hashList []string, tx *sql.Tx) error {
	var pks []int64
	for _, hash := range hashList {
		id, err := RepositoryChunkHashSqlite{}.Create(hash, tx)
		if err != nil {
			return fmt.Errorf("falha ao criar lista de hash do chunk: %v", err)
		}

		pks = append(pks, id)
	}

	for _, chunkHashId := range pks {
		err := RepositoryChunkHasChunkHashSqlite{}.Create(chunkid, chunkHashId, tx)
		if err != nil {
			return fmt.Errorf("falha ao criar relacionamento, chunk has chunk hash: %v", err)
		}
	}

	return nil
}

func (rcs RepositoryChunkSqlite) GetChunkHashListByHashOriginalFile(hashOriginalFile string) ([]string, error) {
	var hashList = []string{}
	db, err := getConectionSqlite()
	if err != nil {
		return hashList, err
	}
	defer db.Close()

	fileId, err := RepositoryFileSqlite{}.GetIdByHashFile(hashOriginalFile)
	if err != nil {
		return hashList, fmt.Errorf("falha ao buscar id do arquivo: %v", err)
	}

	id, err := rcs.getIdByFileId(fileId)
	if err != nil {
		return hashList, fmt.Errorf("falha ao buscar id do chunk pelo id do arquivo: %v", err)
	}

	query := `SELECT ch.hash
                FROM chunk_hashes AS ch
          INNER JOIN chunks_has_chunk_hashes AS chch
                  ON ch.id = chch.chunk_hash_id
               WHERE chch.chunk_id = ?`
	rows, err := db.Query(query, id)
	if err == sql.ErrNoRows {
		return hashList, nil
	}

	if err != nil {
		return hashList, err
	}

	defer rows.Close()
	for rows.Next() {
		var item string
		err := rows.Scan(&item)
		if err != nil {
			return nil, err
		}
		hashList = append(hashList, item)
	}

	return hashList, nil
}

func (rcs RepositoryChunkSqlite) CountChunkHash(hash string) (int64, error) {
	var count int64
	db, err := getConectionSqlite()
	if err != nil {
		return count, err
	}
	defer db.Close()

	query := `SELECT COUNT(chch.chunk_id)
	            FROM chunks_has_chunk_hashes AS chch
          INNER JOIN chunk_hashes AS ch
                  ON ch.id = chch.chunk_hash_id
               WHERE ch.hash = ?`

	err = db.QueryRow(query, hash).Scan(&count)
	if err == sql.ErrNoRows {
		return count, nil
	}

	return count, err
}

func (rcs RepositoryChunkSqlite) RemoveAll() error {
	db, err := getConectionSqlite()
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = RepositoryChunkHasChunkHashSqlite{}.RemoveAll(tx)
	if err != nil {
		return err
	}

	err = RepositoryChunkHashSqlite{}.RemoveAll(tx)
	if err != nil {
		return err
	}

	err = rcs.removeAll(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (rcs RepositoryChunkSqlite) removeAll(tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM chunks`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM sqlite_sequence WHERE name='chunks'`)
	return err
}

func (rcs RepositoryChunkSqlite) RemoveByHashOriginalFile(hashOriginalFile string) ([]string, error) {
	var hashList []string
	fileId, err := RepositoryFileSqlite{}.GetIdByHashFile(hashOriginalFile)
	if err != nil {
		return hashList, fmt.Errorf("falha ao buscar id do arquivo: %v", err)
	}

	id, err := rcs.getIdByFileId(fileId)
	if err != nil {
		return hashList, err
	}

	db, err := getConectionSqlite()
	if err != nil {
		return hashList, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return hashList, err
	}
	defer tx.Rollback()

	countItem, err := RepositoryChunkHasChunkHashSqlite{}.CountChunkHashByChunkId(id)
	if err != nil {
		return hashList, err
	}

	err = RepositoryChunkHasChunkHashSqlite{}.RemoveByChunkId(id, tx)
	if err != nil {
		return hashList, err
	}

	for _, item := range countItem {
		if item.Total > 1 {
			continue
		}

		chunkHash, err := RepositoryChunkHashSqlite{}.GetHashById(item.Id)
		if err != nil {
			return hashList, err
		}

		hashList = append(hashList, chunkHash)
		err = RepositoryChunkHashSqlite{}.RemoveByChunkHashId(item.Id, tx)
		if err != nil {
			return hashList, err
		}
	}

	err = rcs.removeByChunkId(id, tx)
	if err != nil {
		return hashList, err
	}

	return hashList, tx.Commit()
}

func (rcs RepositoryChunkSqlite) removeByChunkId(id int64, tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM chunks WHERE id = ?`, id)
	return err
}

func (rcs RepositoryChunkSqlite) getIdByFileId(fileId int64) (int64, error) {
	var id int64
	db, err := getConectionSqlite()
	if err != nil {
		return id, err
	}
	defer db.Close()

	query := `SELECT id FROM chunks WHERE file_id = ?`
	err = db.QueryRow(query, fileId).Scan(&id)
	if err == sql.ErrNoRows {
		return id, nil
	}

	return id, err
}
