# Análise

## Funcionalidade
As funcionalidades permite usar o sofware como planejado, mas falta uma funcionalidade de versionamento de arquivo, o que permitiria tirar mais proveito da tecnica de armazenamento usada.

## Arquitetura
O uso de arquivos auxiliares (config) *.json*, pode ocorrer locks e ficar com conteúdo defasado devido a falha em rotina.

## Código
A lista abaixo são pontos de atenção que serão gerados _issues_:
- Algumas constantes devem ser migradas para arquivo .env para organizar os dados de configuração.
- As estruturas `hashFile`, `hashFileItem`, `chunkFile` e `chunkFileItem` não devem estar relacionados a implementação, `File`, essas estruturas devem atender somente ao modelo de aplicação não fazendo referência a origem de dados.
- As actions `upload`, `donwload`, `clear` e `remove` devem seguir uma máquina de estado para permitir seu processamento assíncrono.
- Deve ser desenvolvido a funcionalidade de *consulta de status* para o usuário saber qual andamento da transação.
- Deve ser desenvolvido a funcionalidade de *versionamento de arquivos*.
- O arquivo action.go faz muita coisa, é necessário quebrar em outros arquivos conforme as responsabilidades.
- A actions devem ter uma interface em comum para melhor flexibilidade.
- Uma fábrica de action deve ser desenvolvida.
- O uso de Logs deve ser embutido na action para uso interno.
- Cada action deve ter seu próprio arquivo `.go`
- Cada etapa de uma action deve seguir ACID.
- A aplicação deve usar *SQLite* para armazenas os arquivos.
- Usar padrão repository junto ao *SQLite*.
- Desenvolver testes unitários.
- Uso de Logs deve ser migrado de arquivo texto para banco de dados.
- Atualizar a documentação conforme as altrações.

## Observação
Existe alguns outros pontos de serão explorados após a conclusão da lista acima.
Mesmo que a aplicação funcione bem, o código precisa de muita alteração para poder ser flexível a mudanças.

[README](./README.md)
