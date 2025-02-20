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
11 service (*Action).FeatureUpload src/internal/service/upload.go:16:1
10 service (*Action).FeatureRemove src/internal/service/remove.go:9:1
9 service Execute src/internal/service/action.go:51:1
8 service (*Action).FeatureClear src/internal/service/clear.go:7:1
7 service (Action).GenerateChunkByHashFile src/internal/service/upload.go:159:1
7 service (Action).GenerateChunksToStorage src/internal/service/upload.go:130:1
5 service (*Action).generateFileByChunkHashList src/internal/service/download.go:48:1
5 service (*Action).FeatureDownload src/internal/service/download.go:11:1
4 service (*Action).GenerateHashFile src/internal/service/upload.go:196:1
4 service (Storage).CreateFile src/internal/service/storage.go:16:1
3 service (TemporaryArea).CreateFileByFileSource src/internal/service/temporary_area.go:35:1
3 service (TemporaryArea).Clear src/internal/service/temporary_area.go:16:1
3 service (Storage).Clear src/internal/service/storage.go:33:1
2 service (*Action).SendFileToTemporaryArea src/internal/service/upload.go:216:1
1 service (Action).isNewFile src/internal/service/upload.go:225:1
1 service (TemporaryArea).RemoveFile src/internal/service/temporary_area.go:50:1
1 service (TemporaryArea).GetFile src/internal/service/temporary_area.go:31:1
1 service MakeServiceTemporaryArea src/internal/service/temporary_area.go:12:1
1 service (Storage).RemoveFile src/internal/service/storage.go:52:1
1 service (Storage).GetFile src/internal/service/storage.go:48:1
1 service MakeServiceStorage src/internal/service/storage.go:12:1
1 service (Action).LogErrorWrite src/internal/service/action.go:94:1
1 service (Action).GetActionType src/internal/service/action.go:90:1
1 service MakeAction src/internal/service/action.go:38:1

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
12 repository TestRepositoryChunkSqliteRemoveAll src/repository/repository_chunk_sqlite_test.go:264:1
12 repository (RepositoryChunkSqlite).RemoveByHashOriginalFile src/repository/repository_chunk_sqlite.go:182:1
11 repository TestRepositoryChunkJsonCreate src/repository/repository_chunk_json_test.go:34:1
10 repository (RepositoryChunkJson).RemoveByHashOriginalFile src/repository/repository_chunk_json.go:100:1
9 repository TestRepositoryFileJsonRemoveAll src/repository/repository_file_json_test.go:206:1
9 repository (RepositoryChunkSqlite).GetChunkHashListByHashOriginalFile src/repository/repository_chunk_sqlite.go:81:1
9 repository TestRepostoryChunkHashSqliteRemoveAllWithTransaction src/repository/repository_chunk_hash_sqlite_test.go:168:1
8 repository TestRepositoryFileJsonCreate src/repository/repository_file_json_test.go:34:1
8 repository (RepositoryFileJson).RemoveByHash src/repository/repository_file_json.go:93:1
8 repository TestRepositoryChunkSqliteRemoveAllSequence src/repository/repository_chunk_sqlite_test.go:331:1
8 repository TestRepositoryChunkJsonRemoveAll src/repository/repository_chunk_json_test.go:213:1
8 repository TestRepostoryChunkHashSqliteRemoveByChunkHashWithTransactionIds src/repository/repository_chunk_hash_sqlite_test.go:224:1
7 repository TestRepositoryFileSqliteRemoveAll src/repository/repository_file_sqlite_test.go:257:1
7 repository TestRepositoryFileJsonRemove src/repository/repository_file_json_test.go:155:1
7 repository (RepositoryFileJson).Create src/repository/repository_file_json.go:12:1
7 repository TestRepositoryChunkSqliteCountChunkHashTwo src/repository/repository_chunk_sqlite_test.go:220:1
7 repository setUpRepositoryChunkSqliteTest src/repository/repository_chunk_sqlite_test.go:13:1
7 repository (RepositoryChunkSqlite).createChunkList src/repository/repository_chunk_sqlite.go:49:1
7 repository (RepositoryChunkHasChunkHashSqlite).CountChunkHashByChunkId src/repository/repository_chunk_has_chunk_hash_sqlite.go:48:1
6 repository TestRepositoryFileJsonRemoveAllEmptyList src/repository/repository_file_json_test.go:266:1
6 repository TestRepositoryChunkSqliteCountChunkHashOne src/repository/repository_chunk_sqlite_test.go:186:1
6 repository (RepositoryChunkSqlite).RemoveAll src/repository/repository_chunk_sqlite.go:151:1
6 repository (RepositoryChunkSqlite).Create src/repository/repository_chunk_sqlite.go:11:1
6 repository TestRepositoryChunkJsonRemoveByHashOriginalFileComplex src/repository/repository_chunk_json_test.go:289:1
6 repository TestRepostoryChunkHashSqliteRemoveByChunkHashIdWithTransaction src/repository/repository_chunk_hash_sqlite_test.go:289:1
6 repository TestRepostoryChunkHashSqliteGetIdByHash src/repository/repository_chunk_hash_sqlite_test.go:87:1
6 repository (RepositoryChunkHashSqlite).RemoveByIdsWithTransaction src/repository/repository_chunk_hash_sqlite.go:59:1
6 repository (RepositoryChunkHasChunkHashSqlite).GetChunkHashIdsByChunkId src/repository/repository_chunk_has_chunk_hash_sqlite.go:83:1
5 repository TestRepositoryFileSqliteRemoveByHash src/repository/repository_file_sqlite_test.go:208:1
5 repository (RepositoryFileSqlite).resetAutoIncrement src/repository/repository_file_sqlite.go:140:1
5 repository (RepositoryFileSqlite).RemoveByHash src/repository/repository_file_sqlite.go:63:1
5 repository (RepositoryFileJson).IsExistsByHash src/repository/repository_file_json.go:70:1
5 repository (RepositoryFileJson).GetHashByName src/repository/repository_file_json.go:47:1
5 repository TestRepositoryChunkSqliteRemoveByHashOriginalFile src/repository/repository_chunk_sqlite_test.go:386:1
5 repository TestRepositoryChunkSqliteGetChunkHashListByHashOriginalFile src/repository/repository_chunk_sqlite_test.go:125:1
5 repository TestRepositoryChunkJsonCountUsedChunkHashTwo src/repository/repository_chunk_json_test.go:176:1
5 repository TestRepositoryChunkJsonGetChunkHashListByHashOriginalFile src/repository/repository_chunk_json_test.go:94:1
5 repository (RepositoryChunkJson).GetChunkHashListByHashOriginalFile src/repository/repository_chunk_json.go:41:1
5 repository (RepositoryChunkJson).Create src/repository/repository_chunk_json.go:12:1
5 repository TestRepostoryChunkHashSqliteGetHashById src/repository/repository_chunk_hash_sqlite_test.go:129:1
4 repository TestRepositoryMakeRepositoryChunk src/repository/repository_test.go:28:1
4 repository TestRepositoryMakeRepositoryFile src/repository/repository_test.go:8:1
4 repository TestRepositoryFileSqliteGetIdByHash src/repository/repository_file_sqlite_test.go:303:1
4 repository TestRepositoryFileSqliteIsExistsByHashNotFound src/repository/repository_file_sqlite_test.go:184:1
4 repository TestRepositoryFileSqliteIsExistsByHash src/repository/repository_file_sqlite_test.go:160:1
4 repository TestRepositoryFileSqliteGetHashByName src/repository/repository_file_sqlite_test.go:113:1
4 repository setUpRepositoryFileSqliteTest src/repository/repository_file_sqlite_test.go:14:1
4 repository (RepositoryFileSqlite).resetTable src/repository/repository_file_sqlite.go:118:1
4 repository (RepositoryFileSqlite).IsExistsByHash src/repository/repository_file_sqlite.go:42:1
4 repository TestRepositoryFileJsonIsExistsByHash src/repository/repository_file_json_test.go:115:1
4 repository TestRepositoryFileJsonGetHashByName src/repository/repository_file_json_test.go:79:1
4 repository TestRepositoryChunkSqliteCountChunkHashZero src/repository/repository_chunk_sqlite_test.go:167:1
4 repository TestRepositoryChunkSqliteCreate src/repository/repository_chunk_sqlite_test.go:84:1
4 repository (RepositoryChunkSqlite).removeByIdWithTransaction src/repository/repository_chunk_sqlite.go:251:1
4 repository TestRepositoryChunkJsonRemoveByHashOriginalFileSimple src/repository/repository_chunk_json_test.go:264:1
4 repository TestRepositoryChunkJsonCountUsedChunkHashOne src/repository/repository_chunk_json_test.go:150:1
4 repository TestRepostoryChunkHashSqliteCreateUniqueConstraintHash src/repository/repository_chunk_hash_sqlite_test.go:64:1
4 repository TestRepostoryChunkHashSqliteCreate src/repository/repository_chunk_hash_sqlite_test.go:44:1
4 repository setUpRepostoryChunkHashSqliteTest src/repository/repository_chunk_hash_sqlite_test.go:12:1
4 repository (RepositoryChunkHashSqlite).RemoveByIdWithTransaction src/repository/repository_chunk_hash_sqlite.go:86:1
4 repository (RepositoryChunkHasChunkHashSqlite).RemoveByChunkId src/repository/repository_chunk_has_chunk_hash_sqlite.go:30:1
3 repository TestRepositoryFileSqliteRemoveByHashNotFound src/repository/repository_file_sqlite_test.go:237:1
3 repository TestRepositoryFileSqliteGetHashByNameNotFound src/repository/repository_file_sqlite_test.go:139:1
3 repository TestRepositoryFileSqliteCreateUniqueConstraintHash src/repository/repository_file_sqlite_test.go:89:1
3 repository TestRepositoryFileSqliteCreateUniqueConstraintName src/repository/repository_file_sqlite_test.go:65:1
3 repository TestRepositoryFileSqliteCreate src/repository/repository_file_sqlite_test.go:47:1
3 repository (RepositoryFileSqlite).GetIdByHash src/repository/repository_file_sqlite.go:101:1
3 repository (RepositoryFileSqlite).RemoveAll src/repository/repository_file_sqlite.go:87:1
3 repository (RepositoryFileSqlite).GetHashByName src/repository/repository_file_sqlite.go:25:1
3 repository TestRepositoryFileJsonIsNotExistsByHash src/repository/repository_file_json_test.go:140:1
3 repository setUpRepositoryFileJsonTest src/repository/repository_file_json_test.go:13:1
3 repository (RepositoryChunkSqlite).getIdByFileId src/repository/repository_chunk_sqlite.go:269:1
3 repository (RepositoryChunkSqlite).CountUsedChunkHash src/repository/repository_chunk_sqlite.go:130:1
3 repository TestRepositoryChunkJsonCountUsedChunkHashZero src/repository/repository_chunk_json_test.go:136:1
3 repository setUpRepositoryChunkJsonTest src/repository/repository_chunk_json_test.go:13:1
3 repository (RepositoryChunkJson).CountUsedChunkHash src/repository/repository_chunk_json.go:74:1
3 repository (RepositoryChunkJson).getCountChunkMap src/repository/repository_chunk_json.go:64:1
3 repository TestRepostoryChunkHashSqliteRemoveByChunkHashIdWithTransactionNotFound src/repository/repository_chunk_hash_sqlite_test.go:324:1
3 repository TestRepostoryChunkHashSqliteRemoveByChunkHashIdsWithTransactionNotFound src/repository/repository_chunk_hash_sqlite_test.go:271:1
3 repository (RepositoryChunkHashSqlite).GetHashById src/repository/repository_chunk_hash_sqlite.go:32:1
3 repository (RepositoryChunkHashSqlite).GetIdByHash src/repository/repository_chunk_hash_sqlite.go:15:1
3 repository (RepositoryChunkHasChunkHashSqlite).RemoveAllWithTransaction src/repository/repository_chunk_has_chunk_hash_sqlite.go:16:1
3 repository MakeRepositoryChunk src/repository/repository.go:49:1
3 repository MakeRepositoryFile src/repository/repository.go:38:1
2 repository TestRepositoryFileSqliteGetIdByHashNotFound src/repository/repository_file_sqlite_test.go:327:1
2 repository setDownRepositoryFileSqliteTest src/repository/repository_file_sqlite_test.go:40:1
2 repository (RepositoryFileSqlite).Create src/repository/repository_file_sqlite.go:12:1
2 repository TestRepositoryFileJsonRemoveNotFound src/repository/repository_file_json_test.go:194:1
2 repository TestRepositoryFileJsonGetHashByNameNotFound src/repository/repository_file_json_test.go:104:1
2 repository setDownRepositoryFileJsonTest src/repository/repository_file_json_test.go:27:1
2 repository TestRepositoryChunkSqliteRemoveByHashOriginalFileNotFound src/repository/repository_chunk_sqlite_test.go:420:1
2 repository TestRepositoryChunkSqliteGetChunkHashListByHashOriginalFileNotFound src/repository/repository_chunk_sqlite_test.go:156:1
2 repository TestRepositoryChunkSqliteCreateNotFoundFileHash src/repository/repository_chunk_sqlite_test.go:109:1
2 repository setDownRepositoryChunkSqliteTest src/repository/repository_chunk_sqlite_test.go:77:1
2 repository (RepositoryChunkSqlite).removeAllWithTransaction src/repository/repository_chunk_sqlite.go:241:1
2 repository TestRepositoryChunkJsonGetChunkHashListByHashOriginalFileNotFound src/repository/repository_chunk_json_test.go:125:1
2 repository setDownRepositoryChunkJsonTest src/repository/repository_chunk_json_test.go:27:1
2 repository (RepositoryChunkJson).isChunkHashCanBeRemoved src/repository/repository_chunk_json.go:150:1
2 repository (RepositoryChunkJson).RemoveAll src/repository/repository_chunk_json.go:92:1
2 repository TestRepostoryChunkHashSqliteGetHashByIdNotFound src/repository/repository_chunk_hash_sqlite_test.go:157:1
2 repository TestRepostoryChunkHashSqliteGetIdByHashNotFound src/repository/repository_chunk_hash_sqlite_test.go:118:1
2 repository setDownRepostoryChunkHashSqliteTest src/repository/repository_chunk_hash_sqlite_test.go:37:1
2 repository (RepositoryChunkHashSqlite).RemoveAllWithTransaction src/repository/repository_chunk_hash_sqlite.go:49:1
2 repository setUpRepostoryChunkAsChunkHashSqliteTest src/repository/repository_chunk_as_chunk_hash_sqlite_test.go:12:1
2 repository ping src/repository/repository.go:72:1
2 repository getConectionSqlite src/repository/repository.go:60:1
1 repository (RepositoryFileJson).RemoveAll src/repository/repository_file_json.go:133:1
1 repository (RepositoryChunkSqlite).create src/repository/repository_chunk_sqlite.go:43:1
1 repository (RepositoryChunkHashSqlite).Create src/repository/repository_chunk_hash_sqlite.go:9:1
1 repository (RepositoryChunkHasChunkHashSqlite).Create src/repository/repository_chunk_has_chunk_hash_sqlite.go:10:1
1 repository TestRepositoryChunkHasChunkHashSqliteGetChunkHashIdsByChunkIdNotFound src/repository/repository_chunk_as_chunk_hash_sqlite_test.go:189:1
1 repository TestRepositoryChunkHasChunkHashSqliteGetChunkHashIdsByChunkId src/repository/repository_chunk_as_chunk_hash_sqlite_test.go:168:1
1 repository TestRepositoryChunkHasChunkHashSqliteCountChunkHashByChunkIdNotFound src/repository/repository_chunk_as_chunk_hash_sqlite_test.go:159:1
1 repository TestRepositoryChunkHasChunkHashSqliteCountChunkHashByChunkId src/repository/repository_chunk_as_chunk_hash_sqlite_test.go:134:1
1 repository TestRepositoryChunkHasChunkHashSqliteRemoveByChunkIdNotFound src/repository/repository_chunk_as_chunk_hash_sqlite_test.go:121:1
1 repository TestRepositoryChunkHasChunkHashSqliteRemoveByChunkId src/repository/repository_chunk_as_chunk_hash_sqlite_test.go:98:1
1 repository TestRepositoryChunkHasChunkHashSqliteRemoveAllWithTransaction src/repository/repository_chunk_as_chunk_hash_sqlite_test.go:76:1
1 repository TestRepositoryChunkHasChunkHashSqliteCreate src/repository/repository_chunk_as_chunk_hash_sqlite_test.go:55:1
1 repository setDownRepostoryChunkAsChunkHashSqliteTest src/repository/repository_chunk_as_chunk_hash_sqlite_test.go:51:1
```

## Cobertura de testes
pasta: */src/repository*
- `go test -cover` coverage: 80.3% of statements
- `go test -v -cover`
- `go test -coverprofile=coverage.out`
- `go tool cover -html=coverage.out`

[README](./README.md)
