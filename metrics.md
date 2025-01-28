## Métricas

## Objetivo
É sempre importante acompanhar as métricas ao longo da vida de um sofware, isso garante que as coisas irão se manter estavel, sem surpresas ao com o passar do tempo.

## gocyclo
Para gerar métrica de complexidade ciclomática foi utilizado o [gocyclo](https://pkg.go.dev/github.com/fzipp/gocyclo).

## Escala usada
1. **Complexidade 1-10**: Funções com uma complexidade ciclomática nesta faixa são geralmente consideradas simples e fáceis de entender. Este é o alvo ideal para a maioria das funções.
2. **Complexidade 11-20**: Funções nesta faixa começam a se tornar mais complexas. Pode ser necessário refatorar essas funções para melhorar a legibilidade e a manutenibilidade.
3. **Complexidade 21-50**: Funções com essa pontuação são bastante complexas e podem ser difíceis de entender e manter. Refatoração é altamente recomendada.
4. **Complexidade > 50**: Funções com uma complexidade ciclomática acima de 50 são extremamente complexas e provavelmente precisam ser divididas em funções menores e mais gerenciáveis.

## Overview
```bash
gocyclo ./src/cmd
8 main main src/cmd/main.go:18:1
2 main isValidConfigRepositorty src/cmd/main.go:91:1
2 main isValidArgAction src/cmd/main.go:75:1
1 main isValidArgFileTarget src/cmd/main.go:87:1

gocyclo ./src/internal/service/
9 service (*Action).FeatureUpload src/internal/service/upload.go:15:1
9 service Execute src/internal/service/action.go:41:1
7 service (Action).GenerateChunkByHashFile src/internal/service/upload.go:105:1
7 service (Action).GenerateChunksToStorage src/internal/service/upload.go:76:1
7 service (*Action).FeatureRemove src/internal/service/remove.go:8:1
7 service (*Action).FeatureClear src/internal/service/clear.go:7:1
5 service (*Action).generateFileByChunkHashList src/internal/service/download.go:26:1
4 service (Storage).CreateFile src/internal/service/storage.go:16:1
3 service (*Action).GenerateHashFile src/internal/service/upload.go:142:1
3 service (TemporaryArea).CreateFileByFileSource src/internal/service/temporary_area.go:35:1
3 service (TemporaryArea).Clear src/internal/service/temporary_area.go:16:1
3 service (Storage).Clear src/internal/service/storage.go:33:1
3 service (*Action).FeatureDownload src/internal/service/download.go:9:1
2 service (Action).isNewFile src/internal/service/upload.go:167:1
2 service (*Action).SendFileToTmp src/internal/service/upload.go:158:1
1 service (TemporaryArea).RemoveFile src/internal/service/temporary_area.go:50:1
1 service (TemporaryArea).GetFile src/internal/service/temporary_area.go:31:1
1 service MakeServiceTemporaryArea src/internal/service/temporary_area.go:12:1
1 service (Storage).RemoveFile src/internal/service/storage.go:52:1
1 service (Storage).GetFile src/internal/service/storage.go:48:1
1 service MakeServiceStorage src/internal/service/storage.go:12:1
1 service (Action).GetActionType src/internal/service/action.go:84:1
1 service MakeAction src/internal/service/action.go:28:1

gocyclo ./src/logger/
3 logger (*Log).Clear src/logger/log.go:90:1
3 logger (*Log).WriteLog src/logger/log.go:31:1
2 logger (*Log).WriteLogMessageInfo src/logger/log.go:80:1
2 logger (*Log).WriteLogMessageError src/logger/log.go:70:1
2 logger (*Log).clearLogActivity src/logger/log.go:57:1
2 logger (*Log).clearLogError src/logger/log.go:45:1
1 logger GetLogError src/logger/log.go:27:1
1 logger GetLogActivity src/logger/log.go:23:1
1 logger LogSetConfig src/logger/log.go:19:1

gocyclo ./src/repository/
12 repository (RepositoryChunkSqlite).RemoveByHashOriginalFile src/repository/repository_chunk_sqlite.go:179:1
10 repository (RepositoryChunkJson).RemoveByHashOriginalFile src/repository/repository_chunk_json.go:100:1
8 repository (RepositoryChunkSqlite).GetChunkHashListByHashOriginalFile src/repository/repository_chunk_sqlite.go:71:1
7 repository (RepositoryFileJson).RemoveByHashFile src/repository/repository_file_json.go:93:1
7 repository (RepositoryFileJson).Create src/repository/repository_file_json.go:12:1
7 repository (RepositoryChunkHasChunkHashSqlite).CountChunkHashByChunkId src/repository/repository_chunk_has_chunk_hash_sqlite.go:28:1
6 repository (RepositoryChunkSqlite).RemoveAll src/repository/repository_chunk_sqlite.go:138:1
6 repository (RepositoryChunkSqlite).Create src/repository/repository_chunk_sqlite.go:11:1
6 repository (RepositoryChunkHasChunkHashSqlite).GetChunkHashIdsByChunkId src/repository/repository_chunk_has_chunk_hash_sqlite.go:63:1
5 repository (RepositoryFileJson).IsExistsByHashFile src/repository/repository_file_json.go:70:1
5 repository (RepositoryFileJson).GetHashByName src/repository/repository_file_json.go:47:1
5 repository (RepositoryChunkSqlite).createChunkList src/repository/repository_chunk_sqlite.go:50:1
5 repository (RepositoryChunkJson).GetChunkHashListByHashOriginalFile src/repository/repository_chunk_json.go:41:1
5 repository (RepositoryChunkJson).Create src/repository/repository_chunk_json.go:12:1
4 repository (RepositoryFileSqlite).IsExistsByHashFile src/repository/repository_file_sqlite.go:42:1
4 repository (RepositoryChunkHashSqlite).RemoveByChunkHashIds src/repository/repository_chunk_hash_sqlite.go:73:1
3 repository (RepositoryFileSqlite).GetIdByHashFile src/repository/repository_file_sqlite.go:90:1
3 repository (RepositoryFileSqlite).RemoveAll src/repository/repository_file_sqlite.go:74:1
3 repository (RepositoryFileSqlite).GetHashByName src/repository/repository_file_sqlite.go:25:1
3 repository (RepositoryChunkSqlite).getIdByFileId src/repository/repository_chunk_sqlite.go:243:1
3 repository (RepositoryChunkSqlite).CountChunkHash src/repository/repository_chunk_sqlite.go:116:1
3 repository (RepositoryChunkJson).CountChunkHash src/repository/repository_chunk_json.go:74:1
3 repository (RepositoryChunkJson).getCountChunkMap src/repository/repository_chunk_json.go:64:1
3 repository (RepositoryChunkHashSqlite).GetHashById src/repository/repository_chunk_hash_sqlite.go:46:1
3 repository (RepositoryChunkHashSqlite).GetIdByHash src/repository/repository_chunk_hash_sqlite.go:29:1
3 repository (RepositoryChunkHashSqlite).Create src/repository/repository_chunk_hash_sqlite.go:10:1
3 repository MakeRepositoryChunk src/repository/repository.go:46:1
3 repository MakeRepositoryFile src/repository/repository.go:35:1
2 repository (RepositoryFileSqlite).RemoveByHashFile src/repository/repository_file_sqlite.go:63:1
2 repository (RepositoryFileSqlite).Create src/repository/repository_file_sqlite.go:12:1
2 repository (RepositoryFileJson).RemoveAll src/repository/repository_file_json.go:127:1
2 repository (RepositoryChunkSqlite).removeAll src/repository/repository_chunk_sqlite.go:169:1
2 repository (RepositoryChunkJson).isChunkHashCanBeRemoved src/repository/repository_chunk_json.go:150:1
2 repository (RepositoryChunkJson).RemoveAll src/repository/repository_chunk_json.go:92:1
2 repository (RepositoryChunkHashSqlite).RemoveAll src/repository/repository_chunk_hash_sqlite.go:63:1
2 repository (RepositoryChunkHasChunkHashSqlite).RemoveAll src/repository/repository_chunk_has_chunk_hash_sqlite.go:13:1
2 repository ping src/repository/repository.go:66:1
2 repository getConectionSqlite src/repository/repository.go:57:1
1 repository (RepositoryChunkSqlite).removeByChunkId src/repository/repository_chunk_sqlite.go:238:1
1 repository (RepositoryChunkSqlite).create src/repository/repository_chunk_sqlite.go:43:1
1 repository (RepositoryChunkHashSqlite).RemoveByChunkHashId src/repository/repository_chunk_hash_sqlite.go:91:1
1 repository (RepositoryChunkHashSqlite).create src/repository/repository_chunk_hash_sqlite.go:23:1
1 repository (RepositoryChunkHasChunkHashSqlite).RemoveByChunkId src/repository/repository_chunk_has_chunk_hash_sqlite.go:23:1
1 repository (RepositoryChunkHasChunkHashSqlite).Create src/repository/repository_chunk_has_chunk_hash_sqlite.go:7:1
```

[README](./README.md)
