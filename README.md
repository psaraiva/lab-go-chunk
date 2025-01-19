# Lab: Go Chunck (CLI)

## Objetivo
Esse lab começou com um simples exemplo de uso armazenamento de partes de arquivo para otimizar o armazenamento. Mas como toda aplicação ao desevolver foi crescendo e criando mais funcionalides.

## Cenário
A aplicação funciona _"stand alone"_, tem interação com o disco e usa algumas pastas para armazenar os arquivos em suas etapas de processamento.

## Limitações
Seguindo a premissa de começar pequeno e expandir, esse lab começa pequeno e tem muitas limitações:
- Uso somente de disco para armazenar e processar os arquivos.
- Uso de arquivos **.json** para auxiliar o controle e _"integridade"_ dos processo.
- Uso síncrono e único, a aplicação não foi projetada para uso em paralelo.
- Arquivos auxiliares podem ficarem grandes, deteriorando o desempenho da aplicação conforme o uso.

## Documentação
A documentação foi gerada no [platUML](https://plantuml.com/) e estão no diretório `./doc`
[doc-comentada](./documentation.md)

## Como usar
clear: Restaura aplicação para estado inicial.
- `go run main.go -action clear`

upload: Envia o arquivo para o sistema.
- `go run main.go -action upload -file-target ./input_file_examples/day_book_17.07.txt`

dowload: Carrega o arquivo do sistema.
- `go run main.go -action download -file-target day_book_17.07.txt`

remove: Remove um arquivo do sistema.
- `go run main.go -action remove -file-target day_book_18.07.txt`

## Logs
A aplicação tem dois tipos de logs:
- Log de atividades: `./log/activity.log`
- Log de erros: `./log/error.log`

## Análise
Análise sobre a aplicação usando ponto vista do macro para micro, do ponto de vista negócio, arquitetura, boas práticas etc.
[Veja a análise](./analysis.md)

## Métricas
Foi utilizado como métrica a complexidade ciclomática. [Veja os resultados](./metrics.md)
