package daysteps

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"my-app/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// 1. Разделить строку на слайс строк по запятой
	parts := strings.Split(data, ",")

	// 2. Проверить, что длина слайса равна 2
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных")
	}

	// 3. Преобразовать первый элемент шаги в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.New("неверный формат данных")
	}

	// 4. Проверить, что количество шагов больше 0
	if steps <= 0 {
		return 0, 0, errors.New("неверный формат данных")
	}

	// 5. Преобразовать второй элемент в time.Duration
	duration, err := time.ParseDuration(parts[1])
	if duration <= 0 {
		return 0, 0, errors.New("неверный формат данных")
	}

	if err != nil {
		return 0, 0, errors.New("неверный формат данных")
	}

	// 6. Если всё успешно, возвращаем данные
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// 1. Получаем данные с помощью parsePackage
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}

	// 2. Проверяем, чтобы количество шагов было больше 0
	if steps <= 0 || duration <= 0 {
		return ""
	}

	distanceInMeters := float64(steps) * stepLength

	distanceInKm := distanceInMeters / mInKm

	durationInHours := duration.Hours()

	d := time.Duration(durationInHours * float64(time.Hour))

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, d)

	if err != nil {
		// Если функция вернула ошибку, нужно её обработать
		log.Println(err)
		return ""
	}

	// 6. Формируем и возвращаем результирующую строку
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps,
		distanceInKm,
		calories,
	)
}
