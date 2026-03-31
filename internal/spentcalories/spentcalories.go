package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	LenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	// TODO: реализовать функцию

	// 1. Разделить строку на слайс строк по запятой
	parts := strings.Split(data, ",")

	// 2. Проверить, чтобы длина слайса была равна 3
	if len(parts) != 3 {
		return 0, "", 0, errors.New("неверный формат данных: ожидается 3 элемента")
	}

	// 3. Преобразовать первый элемент (шаги) в int
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		// При возникновении ошибки вернуть 0 шагов, 0 длительности и ошибку
		return 0, "", 0, err
	}

	// Вид активности берем как есть (второй элемент слайса)
	activityType := parts[1]

	if steps <= 0 {
		return 0, "", 0, errors.New("некорректное количество шагов")
	}

	// 4. Преобразовать третий элемент (длительность) в time.Duration
	duration, err := time.ParseDuration(parts[2])
	if duration <= 0 {
		return 0, "", 0, errors.New("некорректная продолжительность")
	}
	if err != nil {
		// При возникновении ошибки вернуть 0 шагов, 0 длительности и ошибку
		return 0, "", 0, err
	}

	// 5. Если всё прошло успешно, вернуть результат и nil для ошибки
	return steps, activityType, duration, nil
}

func distance(steps int, height float64) float64 {
	// TODO: реализовать функцию

	// 1. Рассчитываем длину шага.
	// Умножаем рост на коэффициент длины шага.
	stepLength := height * stepLengthCoefficient

	// 2. Умножаем пройденное количество шагов на длину шага.
	// Важно: приводим steps к типу float64 для корректного вычисления.
	distanceInMeters := float64(steps) * stepLength

	// 3. Разделяем полученное значение на число метров в километре.
	distanceInKm := distanceInMeters / mInKm

	return distanceInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	// TODO: реализовать функцию

	// 1. Проверить, что продолжительность duration больше 0.
	// В time.Duration 0 — это просто 0.
	if duration <= 0 {
		return 0
	}

	// 2. Вычислить дистанцию в километрах с помощью ранее созданной функции distance().
	dist := distance(steps, height)

	// 3. Вычислить и вернуть среднюю скорость.
	// Чтобы перевести продолжительность в часы (float64),
	// используем метод .Hours() из пакета time.
	speed := dist / duration.Hours()

	return speed
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	// TODO: реализовать функцию

	// 1. Получаем значения с помощью parseTraining()
	steps, activityType, duration, err := parseTraining(data)
	if err != nil {
		// Логируем ошибку и возвращаем её
		log.Println(err)
		return "", err
	}

	var calories float64
	var calcErr error

	// 2. Проверяем вид тренировки через switch
	switch activityType {
	case "Бег":
		calories, calcErr = RunningSpentCalories(steps, weight, height, duration)
	case "Ходьба":
		calories, calcErr = WalkingSpentCalories(steps, weight, height, duration)
	default:
		// 4. Если тип неизвестен, возвращаем ошибку
		return "", errors.New("неизвестный тип тренировки")
	}

	// Обработка ошибок, возникших при расчете калорий (например, некорректные параметры)
	if calcErr != nil {
		log.Println(calcErr)
		return "", calcErr
	}

	// 3. Формируем и возвращаем строку по образцу
	// Дистанция и скорость вычисляются через наши внутренние функции
	dist := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	report := fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		activityType,
		duration.Hours(),
		dist,
		speed,
		calories,
	)

	return report, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию

	// 1. Проверить входные параметры на корректность.
	// Вес, рост и время должны быть положительными числами.
	if weight <= 0 || height <= 0 || duration <= 0 || steps <= 0 {
		return 0, errors.New("некорректные входные параметры: значения должны быть больше нуля")
	}

	// 2. Рассчитать среднюю скорость с помощью meanSpeed().
	// Она вернет скорость в км/ч.
	speed := meanSpeed(steps, height, duration)

	// 3. Рассчитать количество калорий.

	// Переводим продолжительность в минуты.
	durationInMinutes := duration.Minutes()

	// Используем формулу из ТЗ: (weight * meanSpeed * durationInMinutes) / mInH
	// Важно: mInH — это константа (обычно 60).
	calories := (weight * speed * durationInMinutes) / minInH

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	// TODO: реализовать функцию

	// 1. Проверить входные параметры на корректность.
	// Физические параметры и время должны быть положительными.
	if weight <= 0 || height <= 0 || duration <= 0 || steps <= 0 {
		return 0, errors.New("некорректные параметры: проверьте вес, рост, шаги и длительность")
	}

	// 2. Рассчитать среднюю скорость с помощью meanSpeed().
	speed := meanSpeed(steps, height, duration)

	// 3. Рассчитать базовое количество калорий.

	// Переводим продолжительность в минуты.
	durationInMinutes := duration.Minutes()

	// Используем формулу: (weight * speed * durationInMinutes) / mInH
	// mInH — константа минут в часе (60).
	baseCalories := (weight * speed * durationInMinutes) / minInH

	// 4. Умножить полученное число калорий на корректирующий коэффициент.
	// walkingCaloriesCoefficient определена в пакете.
	finalCalories := baseCalories * walkingCaloriesCoefficient

	// Возвращаем итоговое значение и nil (отсутствие ошибки).
	return finalCalories, nil
}
