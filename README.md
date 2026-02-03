# DIY-VectorDB

Um projeto para aprendizado prático com **Docker**, **bancos de dados vetoriais(Muito utililzados em aplicoes de IA)** e **modelos de linguagem da Ollama**. O foco é explorar a capacidade de buscar informações similares usando embeddings e construir uma API REST simples em **Golang**.

---

## Tecnologias Utilizadas

- **Docker**: Para containerização do ambiente e facilidade de execução.
- **Ollama**: Para geração de embeddings através do modelo `embeddinggemma`.
- **Golang**: Linguagem principal para a API REST.
- **Chi Router**: Para gerenciamento de rotas na API.

---

## Funcionalidades

Atualmente o projeto possui algumas funcionalidades implementadas e outras planejadas:

- ✅ Gerar embeddings a partir de textos.
- ✅ Receber informações via requisições HTTP.
- ⚪ Api rest fazendo aceitando requisções básicas
- ⚪ Armazenar informações e embeddings como chave (em desenvolvimento).
- ⚪ Pesquisar informações similares usando embeddings (em desenvolvimento).

---

## Como Usar

1. **Clone o projeto**

```bash
git clone https://github.com/Thiago-C-Lessa/DIY-VectorDB
cd DIY-VectorDB
```

2. **Inicie o projeto com Docker Compose**  

```bash
docker compose up --build
```
3. **Baixe o modelo embeddinggemma**  
Abra um novo terminal e execute:

```bash
docker exec -it diy-vectordb-ollama-1 ollama pull embeddinggemma:300m
```
Isso é necessário pois a imagem do modelo ainda não está disponível para o conteiner ollama

5. **Reinicie o projeto com Docker Compose**

```bash
docker compose up --build