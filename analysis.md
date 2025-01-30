# Análise

## Funcionalidade
As funcionalidades permite usar o sofware como planejado, mas falta uma funcionalidade de versionamento de arquivo, o que permitiria tirar mais proveito da tecnica de armazenamento usada.

## Arquitetura
A aplicação iniciou de uma forma simpes para atender uma "demanda" pessoal, mas com o amadurecimento, houve diversas alterações visando a flexibilidade e expansão.

- Coleções em arquivos JSON: Os arquivos JSON não são nada performáticos, conforme o uso da aplicação a performance piora devido a quantidade de dados no arquivo.
- Uso de banco de dados SQLite: A subistituição de uso de arquivos JSON para banco de dados SQLite é muito importante parar manter a performance com o uso.
- Repository Pattern: Quem disse que esse padrão é intimamente ligado ao banco de dados? Nessa vesão a aplicação funciona com as duas ENGINE_COLLECTION json|sqlite. Isso comprova que o pattern é eficiente permitindo a troca de ENGINE_COLLECTION de forma simples e transparente para a camada de serviço.
- Uso do Docker/Docker Compose: Permite criar outros recursos que interagem com a aplicação. Nessa versão mesmo com o uso do Docker a aplicação ainda usa o banco de dados SQLite dentro do mesmo container da aplicação. Mas isso vai mudar, novos containers serão usados para camada de repositório, cache e filas para o liberando a aplicação para uso simultâneo/paralelo.

A intenção é exatamente essa, pegar uma aplicação de uma "demanda" pessoal e transformar em algo profissional... Já viu algo semelhante?

## Issues
Para acompanhar a lista de tarefas aberto, consulte o issues desse projeto.

[README](./README.md)
