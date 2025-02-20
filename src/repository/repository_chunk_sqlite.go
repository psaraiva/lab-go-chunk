package repository

import (
	"database/sql"
	"fmt"
	"lab/src/internal/entity"
)

type RepositoryChunkSqlite struct{}

func (rcs RepositoryChunkSqlite) Create(chunk entity.Chunk) (int64, error) {
	var id int64
	db, err := getConectionSqlite()
	if err != nil {
		return id, err
	}
	defer db.Close()

	fileId, err := RepositoryFileSqlite{}.GetIdByHash(chunk.HashOriginalFile)
	if err != nil {
		return id, err
	}

	tx, err := db.Begin()
	if err != nil {
		return id, err
	}
	defer tx.Rollback()

	id, err = rcs.create(chunk, fileId, tx)
	if err != nil {
		return id, err
	}

	err = rcs.createChunkList(id, chunk.HashList, tx)
	if err != nil {
		return id, err
	}

	return id, tx.Commit()
}

func (rcs RepositoryChunkSqlite) create(chunk entity.Chunk, fileId int64, tx *sql.Tx) (int64, error) {
	var id int64
	dml := `INSERT INTO chunks (size, file_id) VALUES (?, ?) RETURNING id`
	return id, tx.QueryRow(dml, chunk.Size, fileId).Scan(&id)
}

func (rcs RepositoryChunkSqlite) createChunkList(chunkId int64, chunkHashList []string, tx *sql.Tx) error {
	var pks []int64
	for _, chunkHash := range chunkHashList {
		id, err := RepositoryChunkHashSqlite{}.GetIdByHash(chunkHash)
		if err == nil {
			pks = append(pks, id)
			continue
		}

		if err == ErrorRecordNotFound {
			id, err = RepositoryChunkHashSqlite{}.Create(chunkHash, tx)
			if err != nil {
				return err
			}

			pks = append(pks, id)
			continue
		}

		return err
	}

	for _, chunkHashId := range pks {
		err := RepositoryChunkHasChunkHashSqlite{}.Create(chunkId, chunkHashId, tx)
		if err != nil {
			return err
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

	fileId, err := RepositoryFileSqlite{}.GetIdByHash(hashOriginalFile)
	if err == sql.ErrNoRows {
		return hashList, ErrorRecordNotFound
	}

	if err != nil {
		return hashList, err
	}

	id, err := rcs.getIdByFileId(fileId)
	if err != nil {
		return hashList, fmt.Errorf("falha ao buscar id do chunk pelo id do arquivo: %v", err)
	}

	dml := `SELECT ch.hash
                FROM chunk_hashes AS ch
          INNER JOIN chunks_has_chunk_hashes AS chch
                  ON ch.id = chch.chunk_hash_id
               WHERE chch.chunk_id = ?`
	rows, err := db.Query(dml, id)
	if err == sql.ErrNoRows {
		return hashList, ErrorRecordNotFound
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

func (rcs RepositoryChunkSqlite) CountUsedChunkHash(hash string) (int64, error) {
	var count int64
	db, err := getConectionSqlite()
	if err != nil {
		return count, err
	}
	defer db.Close()

	dml := `SELECT COUNT(chch.chunk_id)
             FROM chunks_has_chunk_hashes AS chch
       INNER JOIN chunk_hashes AS ch
               ON ch.id = chch.chunk_hash_id
            WHERE ch.hash = ?`
	err = db.QueryRow(dml, hash).Scan(&count)
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

	err = RepositoryChunkHasChunkHashSqlite{}.RemoveAllWithTransaction(tx)
	if err != nil {
		return err
	}

	err = RepositoryChunkHashSqlite{}.RemoveAllWithTransaction(tx)
	if err != nil {
		return err
	}

	err = rcs.removeAllWithTransaction(tx)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (rcs RepositoryChunkSqlite) RemoveByHashOriginalFile(hashOriginalFile string) ([]string, error) {
	var hashList []string
	fileId, err := RepositoryFileSqlite{}.GetIdByHash(hashOriginalFile)
	if err != nil {
		return hashList, err
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
		err = RepositoryChunkHashSqlite{}.RemoveByIdWithTransaction(item.Id, tx)
		if err != nil {
			return hashList, err
		}
	}

	err = rcs.removeByIdWithTransaction(id, tx)
	if err != nil {
		return hashList, err
	}

	return hashList, tx.Commit()
}

func (rcs RepositoryChunkSqlite) removeAllWithTransaction(tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM chunks`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM sqlite_sequence WHERE name='chunks'`)
	return err
}

func (rcs RepositoryChunkSqlite) removeByIdWithTransaction(id int64, tx *sql.Tx) error {
	result, err := tx.Exec(`DELETE FROM chunks WHERE id = ?`, id)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return ErrorRecordNotFound
	}

	return nil
}

func (rcs RepositoryChunkSqlite) getIdByFileId(fileId int64) (int64, error) {
	var id int64
	db, err := getConectionSqlite()
	if err != nil {
		return id, err
	}
	defer db.Close()

	dml := `SELECT id FROM chunks WHERE file_id = ?`
	err = db.QueryRow(dml, fileId).Scan(&id)
	if err == sql.ErrNoRows {
		return id, ErrorRecordNotFound
	}

	return id, err
}
