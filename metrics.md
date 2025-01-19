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
gocyclo ./main.go
5 main main ./main.go:10:1
2 main validArgAction ./main.go:45:1
1 main validArgFileTarget ./main.go:57:1

gocyclo ./feature/log.go
3 feature (*Log).ClearLog ./feature/log.go:67:1
3 feature (*Log).WriteLog ./feature/log.go:33:1
2 feature (*Log).clearLogActivity ./feature/log.go:93:1
2 feature (*Log).clearLogError ./feature/log.go:81:1
2 feature (*Log).writeLogActivity ./feature/log.go:57:1
2 feature (*Log).writeLogError ./feature/log.go:47:1
1 feature GetLogError ./feature/log.go:29:1
1 feature GetLogActivity ./feature/log.go:25:1
1 feature LogSetConfig ./feature/log.go:21:1

gocyclo ./feature/action.go
11 feature (*Action).processChunk ./feature/action.go:505:1
9 feature (*Action).Execute ./feature/action.go:54:1
8 feature (*Action).actionRemove ./feature/action.go:97:1
7 feature (*Action).addHashToFile ./feature/action.go:450:1
7 feature (*Action).actionClear ./feature/action.go:237:1
7 feature (*Action).removeHashToChunkFile ./feature/action.go:171:1
7 feature (*Action).removeHashToHashFile ./feature/action.go:137:1
5 feature (*Action).getChunksByHash ./feature/action.go:409:1
5 feature (*Action).getHashByFileName ./feature/action.go:386:1
5 feature (*Action).generateFileByChunks ./feature/action.go:363:1
5 feature (*Action).actionUpload ./feature/action.go:317:1
4 feature (*Action).saveChunckBin ./feature/action.go:565:1
4 feature (*Action).sendFileToTmp ./feature/action.go:484:1
4 feature (*Action).isChunkCanBeRemoved ./feature/action.go:205:1
3 feature (*Action).processHash ./feature/action.go:432:1
3 feature (*Action).actionDownload ./feature/action.go:349:1
3 feature (*Action).cleanStorage ./feature/action.go:302:1
3 feature (*Action).cleanTmp ./feature/action.go:287:1
3 feature (*Action).getCountChunkMap ./feature/action.go:227:1
2 feature (*Action).restoreFileConfigHash ./feature/action.go:279:1
2 feature (*Action).restoreFileConfigChunk ./feature/action.go:271:1
1 feature (*Action).removeFileToTmp ./feature/action.go:582:1
1 feature MakeAction ./feature/action.go:50:1
```

[README](./README.md)
