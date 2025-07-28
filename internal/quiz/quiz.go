package quiz

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"quiz_go/internal/ui"
	"quiz_go/internal/stats"
	"github.com/pterm/pterm"

	"github.com/AlecAivazis/survey/v2"
)

type Questao struct {
	ID          int      `json:"id"`
	Questao     string   `json:"questao"`
	Opcoes      []string `json:"opcoes"`
	Resposta    string   `json:"resposta"`
	Explicacao  string   `json:"explicacao"`
	Dificuldade string   `json:"dificuldade"` // "facil", "medio", "dificil"
	Categoria   string   `json:"categoria"`   // "sintaxe", "tipos", "concorrencia", etc.
}

type Quiz struct {
	questoes     []Questao
	stats        stats.Estatisticas
	statsFile    string
	ollamaURL    string
	ollamaModel  string
	usarOllama   bool
}

// Estrutura para requisição ao Ollama
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Estrutura esperada da resposta da IA para questões
type QuestaoGerada struct {
	Questao     string   `json:"questao"`
	Opcoes      []string `json:"opcoes"`
	Resposta    string   `json:"resposta"`
	Explicacao  string   `json:"explicacao"`
	Dificuldade string   `json:"dificuldade"`
	Categoria   string   `json:"categoria"`
}

func NewQuiz() *Quiz {
	q := &Quiz{
		statsFile:   "quiz_stats.json",
		ollamaURL:   "http://localhost:11434/api/generate",
		ollamaModel: "llama3.2", // Pode ser alterado conforme o modelo disponível
		usarOllama:  true,
		questoes: []Questao{
			// Questões de fallback caso o Ollama não esteja disponível
			{
				ID:          1,
				Questao:     "Qual palavra-chave define uma função em Go?",
				Opcoes:      []string{"func", "function", "def", "lambda"},
				Resposta:    "func",
				Explicacao:  "Em Go, usamos a palavra-chave 'func' para definir funções. Exemplo: func minhaFuncao() {}",
				Dificuldade: "facil",
				Categoria:   "sintaxe",
			},
			{
				ID:          2,
				Questao:     "Como declarar uma variável em Go?",
				Opcoes:      []string{"let x = 10", "var x int = 10", "int x = 10", "x := int(10)"},
				Resposta:    "var x int = 10",
				Explicacao:  "Go usa 'var' para declaração explícita de variáveis. Também podemos usar := para declaração curta.",
				Dificuldade: "facil",
				Categoria:   "tipos",
			},
		},
	}

	// Verificar se o Ollama está disponível
	if !q.testarConexaoOllama() {
		fmt.Println(ui.Yellow("⚠️  Ollama não está disponível. Usando questões pré-definidas."))
		q.usarOllama = false
	} else {
		fmt.Println(ui.Green("✅ Ollama conectado! Questões serão geradas dinamicamente."))
	}

	loadedStats, err := stats.CarregarEstatisticas(q.statsFile)
	if err != nil {
		fmt.Printf("Erro ao carregar estatísticas: %v. Iniciando com estatísticas zeradas.\n", err)
	} else {
		q.stats = loadedStats
	}
	return q
}

func (q *Quiz) testarConexaoOllama() bool {
	client := &http.Client{Timeout: 5 * time.Second}
	
	reqBody := OllamaRequest{
		Model:  q.ollamaModel,
		Prompt: "test",
		Stream: false,
	}
	
	jsonData, _ := json.Marshal(reqBody)
	resp, err := client.Post(q.ollamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	
	return resp.StatusCode == 200
}

func (q *Quiz) gerarQuestaoComOllama(dificuldade, categoria string) (*Questao, error) {
	prompt := fmt.Sprintf(`Gere uma questão de múltipla escolha sobre programação Go com as seguintes especificações:

Dificuldade: %s
Categoria: %s

Retorne APENAS um JSON válido no seguinte formato:
{
  "questao": "Texto da pergunta aqui",
  "opcoes": ["opção 1", "opção 2", "opção 3", "opção 4"],
  "resposta": "resposta correta exata (deve ser uma das opções)",
  "explicacao": "Explicação detalhada da resposta",
  "dificuldade": "%s",
  "categoria": "%s"
}

Requisitos:
- A questão deve ser sobre Go/Golang
- Deve ter exatamente 4 opções
- Uma resposta deve estar correta
- A explicação deve ser educativa
- Use português brasileiro
- Não inclua texto adicional, apenas o JSON`, dificuldade, categoria, dificuldade, categoria)

	reqBody := OllamaRequest{
		Model:  q.ollamaModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("erro ao serializar requisição: %v", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Post(q.ollamaURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar com Ollama: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler resposta: %v", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("erro ao decodificar resposta do Ollama: %v", err)
	}

	// Tentar extrair JSON da resposta
	response := strings.TrimSpace(ollamaResp.Response)
	
	// Encontrar o JSON na resposta (às vezes a IA adiciona texto extra)
	startIdx := strings.Index(response, "{")
	endIdx := strings.LastIndex(response, "}")
	
	if startIdx == -1 || endIdx == -1 {
		return nil, fmt.Errorf("JSON não encontrado na resposta")
	}
	
	jsonStr := response[startIdx : endIdx+1]

	var questaoGerada QuestaoGerada
	if err := json.Unmarshal([]byte(jsonStr), &questaoGerada); err != nil {
		return nil, fmt.Errorf("erro ao decodificar questão gerada: %v", err)
	}

	// Validar a questão gerada
	if err := q.validarQuestao(&questaoGerada); err != nil {
		return nil, fmt.Errorf("questão inválida: %v", err)
	}

	questao := &Questao{
		ID:          rand.Intn(10000) + 1000, // ID aleatório
		Questao:     questaoGerada.Questao,
		Opcoes:      questaoGerada.Opcoes,
		Resposta:    questaoGerada.Resposta,
		Explicacao:  questaoGerada.Explicacao,
		Dificuldade: questaoGerada.Dificuldade,
		Categoria:   questaoGerada.Categoria,
	}

	return questao, nil
}

func (q *Quiz) validarQuestao(questao *QuestaoGerada) error {
	if questao.Questao == "" {
		return fmt.Errorf("questão vazia")
	}
	
	if len(questao.Opcoes) != 4 {
		return fmt.Errorf("deve ter exatamente 4 opções, encontradas: %d", len(questao.Opcoes))
	}
	
	// Verificar se a resposta está entre as opções
	respostaEncontrada := false
	for _, opcao := range questao.Opcoes {
		if strings.TrimSpace(opcao) == strings.TrimSpace(questao.Resposta) {
			respostaEncontrada = true
			break
		}
	}
	
	if !respostaEncontrada {
		return fmt.Errorf("resposta '%s' não encontrada nas opções", questao.Resposta)
	}
	
	return nil
}

func (q *Quiz) gerarQuestoes(quantidade int, dificuldade string) []Questao {
	if !q.usarOllama {
		return q.questoes
	}

	categorias := []string{"sintaxe", "tipos", "concorrencia", "bibliotecas", "interfaces", "erros", "estruturas"}
	questoes := make([]Questao, 0, quantidade)

	fmt.Printf("%s Gerando %d questões com IA...\n", ui.Magenta("🤖"), quantidade)
	
	// Barra de progresso
	spinner, _ := pterm.DefaultSpinner.Start(ui.Cyan("Conectando com a IA..."))

	for i := 0; i < quantidade; i++ {
		categoria := categorias[rand.Intn(len(categorias))]
		dif := dificuldade
		
		// Se não especificou dificuldade, escolher aleatoriamente
		if dif == "" {
			dificuldades := []string{"facil", "medio", "dificil"}
			dif = dificuldades[rand.Intn(len(dificuldades))]
		}

		spinner.UpdateText(fmt.Sprintf("Gerando questão %d/%d - %s (%s)", i+1, quantidade, categoria, dif))

		questao, err := q.gerarQuestaoComOllama(dif, categoria)
		if err != nil {
			fmt.Printf("\n%s Erro ao gerar questão %d: %v\n", ui.Red("❌"), i+1, err)
			fmt.Printf("%s Usando questão pré-definida como fallback.\n", ui.Yellow("⚠️"))
			
			// Usar questão de fallback
			if i < len(q.questoes) {
				questoes = append(questoes, q.questoes[i])
			}
			continue
		}

		questoes = append(questoes, *questao)
		time.Sleep(1 * time.Second) // Evitar sobrecarregar a API
	}

	if len(questoes) > 0 {
		spinner.Success(fmt.Sprintf("✅ %d questões geradas pela IA!", len(questoes)))
	} else {
		spinner.Fail("❌ Falha ao gerar questões. Usando questões pré-definidas.")
		return q.questoes
	}

	return questoes
}

func (q *Quiz) MostrarEstatisticas() {
	if q.stats.TotalQuizzes == 0 {
		fmt.Println(ui.Yellow("📊 Nenhuma estatística disponível ainda."))
		return
	}

	fmt.Println(ui.Cyan("╔══════════════════════════════════════════════════════════╗"))
	fmt.Println(ui.Cyan("║") + "                    " + ui.Bold("📊 SUAS ESTATÍSTICAS 📊") + "                    " + ui.Cyan("║"))
	fmt.Println(ui.Cyan("╚══════════════════════════════════════════════════════════╝"))
	fmt.Println()

	fmt.Printf("%s Total de quizzes realizados: %s\n",
		ui.Magenta("🎯"),
		ui.Bold(fmt.Sprintf("%d", q.stats.TotalQuizzes)))

	fmt.Printf("%s Total de acertos: %s de %s\n",
		ui.Green("✅"),
		ui.Bold(fmt.Sprintf("%d", q.stats.TotalAcertos)),
		ui.Bold(fmt.Sprintf("%d", q.stats.TotalQuestoes)))

	fmt.Printf("%s Melhor score: %s questões\n",
		ui.Yellow("🏆"),
		ui.Bold(fmt.Sprintf("%d", q.stats.MelhorScore)))

	fmt.Printf("%s Média de acertos: %s\n",
		ui.Blue("📈"),
		ui.Bold(fmt.Sprintf("%.1f%%", q.stats.MediaPercentual)))

	if q.stats.UltimoQuiz != "" {
		fmt.Printf("%s Último quiz: %s\n",
			ui.Cyan("📅"),
			ui.Bold(q.stats.UltimoQuiz))
	}

	if q.usarOllama {
		fmt.Printf("%s Modo IA: %s (Modelo: %s)\n",
			ui.Green("🤖"),
			ui.Bold("ATIVO"),
			ui.Bold(q.ollamaModel))
	} else {
		fmt.Printf("%s Modo IA: %s\n",
			ui.Red("🤖"),
			ui.Bold("DESATIVADO"))
	}

	fmt.Println()
}

func (q *Quiz) SelecionarModoJogo() string {
	options := []string{
		"🎯 Todas as questões (10 questões)",
		"⚡ Quiz rápido (5 questões aleatórias)",
		"🧠 Apenas questões difíceis",
		"📊 Ver estatísticas",
	}

	// Adicionar opções específicas para IA se disponível
	if q.usarOllama {
		options = append([]string{
			"🤖 IA: Quiz personalizado (5 questões geradas)",
			"🎓 IA: Questões avançadas (3 questões difíceis)",
			"🚀 IA: Desafio extremo (10 questões mistas)",
		}, options...)
	}

	options = append(options, "❌ Sair")

	var modo string
	prompt := &survey.Select{
		Message: "Escolha o modo de jogo:",
		Options: options,
	}

	survey.AskOne(prompt, &modo)
	return modo
}

func (q *Quiz) FiltrarQuestoes(modo string) []Questao {
	switch {
	case strings.Contains(modo, "IA: Quiz personalizado"):
		return q.gerarQuestoes(5, "")
	case strings.Contains(modo, "IA: Questões avançadas"):
		return q.gerarQuestoes(3, "dificil")
	case strings.Contains(modo, "IA: Desafio extremo"):
		return q.gerarQuestoes(10, "")
	case strings.Contains(modo, "Todas as questões"):
		if q.usarOllama {
			return q.gerarQuestoes(10, "")
		}
		return q.questoes
	case strings.Contains(modo, "Quiz rápido"):
		if q.usarOllama {
			return q.gerarQuestoes(5, "")
		}
		// Fallback para questões pré-definidas
		questoesAleatorias := make([]Questao, len(q.questoes))
		copy(questoesAleatorias, q.questoes)

		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(questoesAleatorias), func(i, j int) {
			questoesAleatorias[i], questoesAleatorias[j] = questoesAleatorias[j], questoesAleatorias[i]
		})

		if len(questoesAleatorias) > 5 {
			return questoesAleatorias[:5]
		}
		return questoesAleatorias
	case strings.Contains(modo, "difíceis"):
		if q.usarOllama {
			return q.gerarQuestoes(5, "dificil")
		}
		// Fallback para questões pré-definidas
		var dificeis []Questao
		for _, q := range q.questoes {
			if q.Dificuldade == "dificil" {
				dificeis = append(dificeis, q)
			}
		}
		return dificeis
	default:
		return q.questoes
	}
}

// O resto dos métodos permanecem iguais...
func (q *Quiz) ExecutarQuiz(questoesSelecionadas []Questao) {
	fmt.Println()
	fmt.Printf("%s Você terá %s questões para responder!\n",
		ui.Magenta("📚"),
		ui.Bold(fmt.Sprintf("%d", len(questoesSelecionadas))))
	fmt.Println()

	score := 0
	respostasCorretas := []bool{}
	tempoInicio := time.Now()

	for i, questao := range questoesSelecionadas {
		ui.LimparTela()
		fmt.Println(ui.Cyan("═══════════════════════════════════════════════════════════"))
		fmt.Printf("%s Questão %d de %d | %s | %s\n",
			ui.Yellow("📝"),
			i+1,
			len(questoesSelecionadas),
			ui.Blue(fmt.Sprintf("Categoria: %s", questao.Categoria)),
			q.getDificuldadeIcon(questao.Dificuldade))
		fmt.Println(ui.Cyan("═══════════════════════════════════════════════════════════"))
		fmt.Println()

		var resposta string
		prompt := &survey.Select{
			Message: ui.Bold(questao.Questao),
			Options: questao.Opcoes,
		}

		err := survey.AskOne(prompt, &resposta)
		if err != nil {
			fmt.Printf(ui.Red("Erro ao ler resposta: %v\n"), err)
			continue
		}

		fmt.Println()

		if strings.TrimSpace(resposta) == questao.Resposta {
			fmt.Println(ui.Green("✅ Resposta correta! Parabéns!"))
			score++
			respostasCorretas = append(respostasCorretas, true)
		} else {
			fmt.Printf(ui.Red("❌ Resposta incorreta! A resposta correta é: %s\n"),
				ui.Bold(questao.Resposta))
			respostasCorretas = append(respostasCorretas, false)
		}

		fmt.Printf("%s %s\n", ui.Blue("💡 Explicação:"), questao.Explicacao)
		fmt.Println()

		if i < len(questoesSelecionadas)-1 {
			fmt.Printf(ui.Magenta("📊 Progresso: %d/%d questões | Acertos: %d\n"),
				i+1, len(questoesSelecionadas), score)
			fmt.Println()

			var continuar bool
			continuePrompt := &survey.Confirm{
				Message: "Continuar para a próxima questão?",
				Default: true,
			}
			survey.AskOne(continuePrompt, &continuar)

			if !continuar {
				fmt.Println(ui.Yellow("Quiz interrompido pelo usuário."))
				return
			}
			fmt.Println()
		}
	}

	tempoTotal := time.Since(tempoInicio)
	q.MostrarResultados(score, len(questoesSelecionadas), respostasCorretas, tempoTotal)
	q.AtualizarEstatisticas(score, len(questoesSelecionadas))
}

func (q *Quiz) getDificuldadeIcon(dificuldade string) string {
	switch dificuldade {
	case "facil":
		return ui.Green("🟢 Fácil")
	case "medio":
		return ui.Yellow("🟡 Médio")
	case "dificil":
		return ui.Red("🔴 Difícil")
	default:
		return ui.Blue("🔵 Normal")
	}
}

func (q *Quiz) MostrarResultados(score, total int, respostasCorretas []bool, tempo time.Duration) {
	fmt.Println()
	fmt.Println(ui.Cyan("╔══════════════════════════════════════════════════════════╗"))
	fmt.Println(ui.Cyan("║") + "                    " + ui.Bold("🏆 RESULTADOS FINAIS 🏆") + "                    " + ui.Cyan("║"))
	fmt.Println(ui.Cyan("╚══════════════════════════════════════════════════════════╝"))
	fmt.Println()

	spinner, _ := pterm.DefaultSpinner.Start(ui.Magenta("Calculando resultados..."))
	time.Sleep(2 * time.Second)
	spinner.Success(pterm.Green("Cálculos finalizados!"))
	fmt.Println()

	percentual := float64(score) / float64(total) * 100

	fmt.Printf("%s Você acertou %s de %s questões\n",
		ui.Magenta("📊"),
		ui.Bold(ui.Green(fmt.Sprintf("%d", score))),
		ui.Bold(fmt.Sprintf("%d", total)))

	fmt.Printf("%s Percentual de acertos: %s\n",
		ui.Magenta("📈"),
		ui.Bold(fmt.Sprintf("%.1f%%", percentual)))

	fmt.Printf("%s Tempo total: %s\n",
		ui.Blue("⏱️"),
		ui.Bold(fmt.Sprintf("%.1f segundos", tempo.Seconds())))

	fmt.Printf("%s Tempo médio por questão: %s\n",
		ui.Blue("⚡"),
		ui.Bold(fmt.Sprintf("%.1f segundos", tempo.Seconds()/float64(total))))

	fmt.Println()

	fmt.Println(ui.Cyan("📋 Resumo das suas respostas:"))
	for i, correto := range respostasCorretas {
		status := ui.Red("❌")
		if correto {
			status = ui.Green("✅")
		}
		fmt.Printf("   Questão %d: %s\n", i+1, status)
	}
	fmt.Println()

	q.MostrarMensagemFinal(score, total, percentual)
}

func (q *Quiz) MostrarMensagemFinal(score, total int, percentual float64) {
	switch {
	case score == total:
		fmt.Println(ui.Green("🎉 PERFEITO! Você acertou todas as questões!"))
		fmt.Println(ui.Green("🏆 Você é um verdadeiro expert em Go!"))
		fmt.Println(ui.Green("🌟 Considerado um GoGuru!"))
	case percentual >= 80:
		fmt.Println(ui.Green("🌟 Excelente! Você tem um ótimo conhecimento em Go!"))
		fmt.Println(ui.Green("👏 Continue assim!"))
		fmt.Println(ui.Blue("🚀 Próximo nível: tente as questões difíceis!"))
	case percentual >= 60:
		fmt.Println(ui.Yellow("👍 Muito bem! Você está no caminho certo!"))
		fmt.Println(ui.Yellow("📚 Continue estudando para melhorar ainda mais!"))
		fmt.Println(ui.Blue("💡 Dica: revise os conceitos que errou!"))
	case percentual >= 40:
		fmt.Println(ui.Yellow("😊 Bom começo! Você já sabe algumas coisas sobre Go!"))
		fmt.Println(ui.Yellow("💪 Com mais estudo você chegará lá!"))
		fmt.Println(ui.Blue("📖 Recomendo focar nos fundamentos primeiro!"))
	default:
		fmt.Println(ui.Red("📖 Você precisa estudar mais sobre Go!"))
		fmt.Println(ui.Red("💡 Que tal revisar a documentação oficial?"))
		fmt.Println(ui.Yellow("🔗 Recursos recomendados:"))
		fmt.Println(ui.Cyan("   • https://golang.org/doc/"))
		fmt.Println(ui.Cyan("   • https://tour.golang.org/"))
		fmt.Println(ui.Cyan("   • https://gobyexample.com/"))
	}
	fmt.Println()
}

func (q *Quiz) AtualizarEstatisticas(score, total int) {
	q.stats.TotalQuizzes++
	q.stats.TotalAcertos += score
	q.stats.TotalQuestoes += total

	if score > q.stats.MelhorScore {
		q.stats.MelhorScore = score
	}

	q.stats.MediaPercentual = float64(q.stats.TotalAcertos) / float64(q.stats.TotalQuestoes) * 100
	q.stats.UltimoQuiz = time.Now().Format("02/01/2006 15:04")

	_ = stats.SalvarEstatisticas(q.statsFile, q.stats)
}

func (q *Quiz) JogarNovamente() bool {
	var jogarNovamente bool
	playAgainPrompt := &survey.Confirm{
		Message: "Gostaria de jogar novamente?",
		Default: false,
	}
	survey.AskOne(playAgainPrompt, &jogarNovamente)
	return jogarNovamente
}