package repository

import "database/sql"

type RepositoryChunkHasChunkHashSqlite struct{}

func (rchchs RepositoryChunkHasChunkHashSqlite) Create(chunkId int64, chunkHashId int64, tx *sql.Tx) error {
	dml := `INSERT INTO chunks_has_chunk_hashes (chunk_id, chunk_hash_id) VALUES (?, ?)`
	_, err := tx.Exec(dml, chunkId, chunkHashId)
	return err
}

func (rchchs RepositoryChunkHasChunkHashSqlite) RemoveAllWithTransaction(tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM chunks_has_chunk_hashes`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`DELETE FROM sqlite_sequence WHERE name='chunks_has_chunk_hashes'`)
	return err
}

func (rchchs RepositoryChunkHasChunkHashSqlite) RemoveByChunkId(id int64, tx *sql.Tx) error {
	_, err := tx.Exec(`DELETE FROM chunks_has_chunk_hashes WHERE chunk_id = ?`, id)
	return err
}

func (rchchs RepositoryChunkHasChunkHashSqlite) CountChunkHashByChunkId(chunkId int64) ([]CountItem, error) {
	var count []CountItem
	db, err := getConectionSqlite()
	if err != nil {
		return count, err
	}
	defer db.Close()

	chunkHashIds, err := rchchs.GetChunkHashIdsByChunkId(chunkId)
	if err != nil {
		return count, err
	}

	stmt, err := db.Prepare(`SELECT COUNT(chunk_hash_id) FROM chunks_has_chunk_hashes WHERE chunk_hash_id = ?`)
	if err != nil {
		return count, err
	}
	defer stmt.Close()

	for _, id := range chunkHashIds {
		var countId int64
		err = stmt.QueryRow(id).Scan(&countId)
		if err != nil && err != sql.ErrNoRows {
			return count, err
		}

		count = append(count, CountItem{
			Id:    id,
			Total: int(countId),
		})
	}

	return count, err
}

func (rchchs RepositoryChunkHasChunkHashSqlite) GetChunkHashIdsByChunkId(chunkId int64) ([]int64, error) {
	var ids []int64
	db, err := getConectionSqlite()
	if err != nil {
		return ids, err
	}
	defer db.Close()

	dml := `SELECT chunk_hash_id FROM chunks_has_chunk_hashes WHERE chunk_id = ?`
	rows, err := db.Query(dml, chunkId)
	if err == sql.ErrNoRows {
		return ids, nil
	}

	if err != nil {
		return ids, err
	}

	defer rows.Close()
	for rows.Next() {
		var item int64
		err := rows.Scan(&item)
		if err != nil {
			return nil, err
		}
		ids = append(ids, item)
	}

	return ids, nil
}
