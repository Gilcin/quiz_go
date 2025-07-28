package stats

import (
	"encoding/json"
	"os"
)

type Estatisticas struct {
	TotalQuizzes    int     `json:"total_quizzes"`
	TotalAcertos    int     `json:"total_acertos"`
	TotalQuestoes   int     `json:"total_questoes"`
	MelhorScore     int     `json:"melhor_score"`
	MediaPercentual float64 `json:"media_percentual"`
	UltimoQuiz      string  `json:"ultimo_quiz"`
}

func CarregarEstatisticas(statsFile string) (Estatisticas, error) {
	var stats Estatisticas
	data, err := os.ReadFile(statsFile)
	if err != nil {
		return stats, err
	}
	json.Unmarshal(data, &stats)
	return stats, nil
}

func SalvarEstatisticas(statsFile string, stats Estatisticas) error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(statsFile, data, 0644)
}
