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
7 main main src/cmd/main.go:13:1
2 main validArgAction src/cmd/main.go:56:1
1 main validArgFileTarget src/cmd/main.go:68:1

gocyclo ./src/internal/actions/
11 actions (*Action).processChunk src/internal/actions/upload.go:155:1
9 actions Execute src/internal/actions/action.go:39:1
8 actions (*Action).FeatureRemove src/internal/actions/remove.go:9:1
7 actions (*Action).addHashToCollection src/internal/actions/upload.go:100:1
7 actions (*Action).removeHashToCollection src/internal/actions/remove.go:115:1
7 actions (*Action).removeHashFileToChunkCollection src/internal/actions/remove.go:81:1
7 actions (*Action).FeatureClear src/internal/actions/clear.go:9:1
6 actions (*Action).isNewFile src/internal/actions/upload.go:53:1
6 actions (*Action).FeatureUpload src/internal/actions/upload.go:15:1
5 actions (*Action).generateFileByChunks src/internal/actions/download.go:26:1
5 actions (*Action).getChunksByHash src/internal/actions/action.go:109:1
5 actions (*Action).getHashByFileName src/internal/actions/action.go:86:1
4 actions (*Action).saveChunkBin src/internal/actions/upload.go:215:1
4 actions (*Action).sendFileToTmp src/internal/actions/upload.go:134:1
4 actions (*Action).isChunkCanBeRemoved src/internal/actions/remove.go:49:1
3 actions (*Action).GenerateHashFile src/internal/actions/upload.go:81:1
3 actions (*Action).getCountChunkMap src/internal/actions/remove.go:71:1
3 actions (*Action).FeatureDownload src/internal/actions/download.go:9:1
3 actions (*Action).cleanTmp src/internal/actions/clear.go:74:1
3 actions (*Action).cleanStorage src/internal/actions/clear.go:59:1
2 actions (*Action).resetHashCollection src/internal/actions/clear.go:51:1
2 actions (*Action).resetChunkCollection src/internal/actions/clear.go:43:1
1 actions (*Action).removeFileToTmp src/internal/actions/upload.go:232:1
1 actions (*Action).GetActionType src/internal/actions/action.go:82:1
1 actions MakeAction src/internal/actions/action.go:35:1

gocyclo ./src/logger/
3 logger (*Log).ClearLog src/logger/log.go:90:1
3 logger (*Log).WriteLog src/logger/log.go:31:1
2 logger (*Log).WriteLogMessageInfo src/logger/log.go:80:1
2 logger (*Log).WriteLogMessageError src/logger/log.go:70:1
2 logger (*Log).clearLogActivity src/logger/log.go:57:1
2 logger (*Log).clearLogError src/logger/log.go:45:1
1 logger GetLogError src/logger/log.go:27:1
1 logger GetLogActivity src/logger/log.go:23:1
1 logger LogSetConfig src/logger/log.go:19:1
```

[README](./README.md)
