package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	// TODO: реализовать функцию

	// 1. Разделить строку на слайс строк по запятой
	parts := strings.Split(data, ",")

	// 2. Проверить, что длина слайса равна 2
	if len(parts) != 2 {
		return 0, 0, errors.New("неверный формат данных: ожидается 'шаги,длительность'")
	}

	// 3. Преобразовать первый элемент шаги в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга шагов: %w", err)
	}

	// 4. Проверить, что количество шагов больше 0
	if steps <= 0 {
		return 0, 0, errors.New("количество шагов должно быть больше 0")
	}

	// 5. Преобразовать второй элемент в time.Duration
	// time.ParseDuration поддерживает форматы типа "3h50m"
	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("ошибка парсинга длительности: %w", err)
	}

	// 6. Если всё успешно, возвращаем данные
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	// TODO: реализовать функцию

	// 1. Получаем данные с помощью parsePackage
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Println("Ошибка парсинга данных:", err)
		return ""
	}

	// 2. Проверяем, чтобы количество шагов было больше 0
	if steps <= 0 {
		return ""
	}

	// 3. Вычисляем дистанцию в метрах
	distanceInMeters := float64(steps) * stepLength

	// 4. Переводим дистанцию в километры
	distanceInKm := distanceInMeters / mInKm

	// 5. Вычисляем калории
	// Переводим duration в часы (float64), так как обычно формулы используют часы
	durationInHours := duration.Hours()
	calories := WalkingSpentCalories(steps, weight, height, durationInHours)

	// 6. Формируем и возвращаем результирующую строку
	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.",
		steps,
		distanceInKm,
		calories,
	)
}
