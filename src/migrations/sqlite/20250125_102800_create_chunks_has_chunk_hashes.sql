CREATE TABLE chunks_has_chunk_hashes (
    chunk_id INTEGER NOT NULL,
    chunk_hash_id INTEGER NOT NULL,
    PRIMARY KEY (chunk_id, chunk_hash_id),
    FOREIGN KEY (chunk_id) REFERENCES chunks(id),
    FOREIGN KEY (chunk_hash_id) REFERENCES chunk_hashes(id)
);