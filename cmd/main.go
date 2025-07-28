package main

import (
	"fmt"
	"strings"

	"quiz_go/internal/ui"
	"quiz_go/internal/quiz"

	"github.com/AlecAivazis/survey/v2"
)


func main() {
	quiz := quiz.NewQuiz()

	for {
		ui.MostrarTelaInicial()

		modo := quiz.SelecionarModoJogo()

		if strings.Contains(modo, "Sair") {
			break
		}

		if strings.Contains(modo, "estatísticas") {
			quiz.MostrarEstatisticas()

			var continuar bool
			prompt := &survey.Confirm{
				Message: "Deseja jogar agora?",
				Default: true,
			}
			survey.AskOne(prompt, &continuar)

			if !continuar {
				continue
			}

			modo = quiz.SelecionarModoJogo()
		}

		questoesSelecionadas := quiz.FiltrarQuestoes(modo)

		if len(questoesSelecionadas) == 0 {
			fmt.Println(ui.Red("❌ Nenhuma questão encontrada para este modo!"))
			continue
		}

		quiz.ExecutarQuiz(questoesSelecionadas)

		if !quiz.JogarNovamente() {
			break
		}
	}

	ui.MostrarDespedida()
}
