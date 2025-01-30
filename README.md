# Lab: Go Chunck (CLI)

## Objetivo
Esse lab começou com um simples exemplo de uso armazenamento de partes de arquivo para otimizar o armazenamento. Mas como toda aplicação ao desevolver foi crescendo e mudanças são necessárias para permitir a expansão.

## Cenário
A aplicação funciona _"stand alone"_, tem interação com o disco e usa algumas pastas para armazenar os arquivos em suas etapas de processamento.

## Limitações
Seguindo a premissa de começar pequeno e expandir, esse lab começa pequeno e tem muitas limitações:
- Uso somente de disco local para armazenar e processar os arquivos.
- ~~Uso de arquivos **.json** usado como coleção de dados, _File_ e _Chunk_.~~
- Uso síncrono e único, a aplicação não foi projetada para uso em paralelo.
- ~~As coleções JSON podem ficarem grandes, deteriorando o desempenho da aplicação.~~

## Documentação
A documentação foi gerada no [platUML](https://plantuml.com/) e estão no diretório [doc](./documentation.md)

## Usando no contexto **Develop**
- build: `docker-compose up --build -d`

Obs:
- Transferir arquivos para pasta **_out_application/input_file_examples_** para usar como arquivo de entrada, input.
- Minhas anotações sobre [Docker](https://gist.github.com/psaraiva/51467d6a49a46709e4c46006ee6015c1)

## Como usar: comandos
clear: Restaura aplicação para estado inicial.
- `docker exec -it lab-go-chunk-cli-dev go run src/cmd/main.go --action clear`

upload: Envia o arquivo para o sistema.
- `docker exec -it lab-go-chunk-cli-dev go run src/cmd/main.go -action upload -file-target out_application/input_file_examples/day_book_17.07.txt`
- `docker exec -it lab-go-chunk-cli-dev go run src/cmd/main.go -action upload -file-target out_application/input_file_examples/day_book_18.07.txt`
- `docker exec -it lab-go-chunk-cli-dev go run src/cmd/main.go -action upload -file-target out_application/input_file_examples/day_book_19.07.txt`
- `docker exec -it lab-go-chunk-cli-dev go run src/cmd/main.go -action upload -file-target out_application/input_file_examples/day_book_20.07.txt`

dowload: Carrega o arquivo do sistema.
- `docker exec -it lab-go-chunk-cli-dev go run src/cmd/main.go -action download -file-target day_book_17.07.txt`

remove: Remove um arquivo do sistema.
- `docker exec -it lab-go-chunk-cli-dev go run src/cmd/main.go -action remove -file-target day_book_18.07.txt`

Acessando o banco de dados sqlite:
- `sqlite3 out_application/collection/sqlite/data_base.db`

**Obs:**
- Minhas anotações sobre [SQLite](https://gist.github.com/psaraiva/23132a7b90dd629467b8efd13fbd1b25)

## Logs
A aplicação tem dois tipos de logs:
- Log de atividades: `out_application/log/activity.log`
- Log de erros: `out_application/log/error.log`

## Análise
Análise sobre a aplicação [link](./analysis.md)

## Métricas
Foi utilizado como métrica a complexidade ciclomática. [link](./metrics.md)
