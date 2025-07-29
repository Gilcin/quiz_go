# Go Quiz CLI 🚀

Um quiz interativo de linha de comando sobre a linguagem de programação Go. Teste seus conhecimentos com questões que vão do básico ao avançado, geradas dinamicamente por uma IA local com Ollama ou usando um conjunto de questões pré-definidas.

 <!-- Substitua por um GIF de demonstração do seu app -->

## ✨ Funcionalidades

- **Geração Dinâmica de Questões**: Integração com [Ollama](https://ollama.com/) para criar questões novas e desafiadoras a cada quiz, sobre diversas categorias de Go.
- **Modo Offline**: Funciona perfeitamente com questões pré-definidas caso o Ollama não esteja disponível ou desativado.
- **Estatísticas de Desempenho**: Acompanhe seu progresso com estatísticas detalhadas, como total de acertos, melhor pontuação e média de acertos.
- **Interface de Terminal Rica**: Experiência de usuário aprimorada com [pterm](https://github.com/pterm/pterm) e [survey](https://github.com/AlecAivazis/survey) para uma navegação colorida e interativa.
- **Feedback Instantâneo**: Receba a resposta correta e uma explicação detalhada após cada pergunta para aprimorar seu aprendizado.

---

## 📋 Pré-requisitos

Para executar este projeto, você precisará de:

- **Go**: Versão 1.18 ou superior.
- **(Opcional) Ollama**: Para a funcionalidade de geração de questões por IA.
  - Siga as instruções de instalação em [ollama.com](https://ollama.com).
  - Após instalar, baixe o modelo usado no projeto (ou outro de sua preferência):
    ```bash
    ollama pull llama3:8b
    ```

## ⚙️ Instalação e Execução

1. **Clone o repositório:**


2. **Instale as dependências:**
   O Go Modules cuidará disso automaticamente. Para garantir que tudo está correto, você pode executar:
   ```bash
   go mod tidy
   ```

3. **Execute a aplicação:**
   ```bash
   go run ./cmd/main.go
   ```

A aplicação irá verificar a conexão com o Ollama e iniciar o menu principal.


---

## 🔧 Configuração

As configurações do Ollama (URL do servidor e modelo a ser usado) podem ser ajustadas diretamente no arquivo `internal/quiz/quiz.go`, na função `NewQuiz`.

```go
func NewQuiz() *Quiz {
	q := &Quiz{
		statsFile:   "quiz_stats.json",
		ollamaURL:   "http://localhost:11434/api/generate", // Altere se seu Ollama estiver em outro endereço
		ollamaModel: "llama3:8b",                           // Altere para outro modelo se desejar
		usarOllama:  true,
        // ...
    }
    // ...
}
```

---

## 📂 Estrutura do Projeto

O projeto está organizado da seguinte forma para manter o código limpo e modular:

```
.
├── cmd/
│   └── main.go         # Ponto de entrada da aplicação, lida com o loop principal
├── internal/
│   ├── quiz/
│   │   └── quiz.go     # Lógica principal do quiz, geração de questões, estatísticas
│   ├── stats/
│   │   └── stats.go    # Estruturas e funções para carregar/salvar estatísticas
│   └── ui/
│       └── ui.go       # Funções de ajuda para a interface do usuário (cores, telas)
├── go.mod
├── go.sum
└── quiz_stats.json     # Arquivo de estatísticas (gerado após o primeiro quiz)
```

---


## 📄 Licença

Distribuído sob a licença MIT. Veja o arquivo `LICENSE` para mais informações.