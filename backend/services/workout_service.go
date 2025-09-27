package services

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"backend/models"
	"backend/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WorkoutService struct {
	repo repositories.WorkoutRepositoryInterface
}

type RoutineUsage struct {
	RoutineID string `json:"routine_id"`
	Count     int    `json:"count"`
}

type DataPoint struct {
	Timestamp time.Time `json:"timestamp"`
	Value     float64   `json:"value"`
}

func NewWorkoutService(repo repositories.WorkoutRepositoryInterface) *WorkoutService {
	return &WorkoutService{repo: repo}
}

func (s *WorkoutService) GetWorkouts(userHex string) ([]models.Workout, error) {
	if userHex == "" {
		return nil, errors.New("user id required")
	}
	uid, err := primitive.ObjectIDFromHex(userHex)
	if err != nil {
		return nil, err
	}
	return s.repo.GetWorkouts(uid)
}

func (s *WorkoutService) GetWorkoutByID(idHex, userHex string) (models.Workout, error) {
	if idHex == "" {
		return models.Workout{}, errors.New("id required")
	}
	w, err := s.repo.GetWorkoutByID(idHex)
	if err != nil {
		return models.Workout{}, err
	}
	if userHex != "" {
		uid, err := primitive.ObjectIDFromHex(userHex)
		if err != nil {
			return models.Workout{}, err
		}
		if w.UserID != uid {
			return models.Workout{}, errors.New("forbidden: not the owner")
		}
	}
	return w, nil
}

func (s *WorkoutService) CreateWorkout(w models.Workout, userHex string) (*mongo.InsertOneResult, error) {
	if userHex == "" {
		return nil, errors.New("user id required")
	}
	uid, err := primitive.ObjectIDFromHex(userHex)
	if err != nil {
		return nil, err
	}
	if w.RoutineID.IsZero() {
		return nil, errors.New("routine id required")
	}
	if w.CompletedAt.IsZero() {
		w.CompletedAt = time.Now()
	}
	w.UserID = uid
	if w.PerformedExercises == nil {
		w.PerformedExercises = []models.ExercisePerformance{}
	}
	return s.repo.CreateWorkout(w)
}

func (s *WorkoutService) UpdateWorkout(idHex string, payload models.Workout, userHex string) (*mongo.UpdateResult, error) {
	if idHex == "" {
		return nil, errors.New("id required")
	}
	existing, err := s.repo.GetWorkoutByID(idHex)
	if err != nil {
		return nil, err
	}
	if userHex == "" {
		return nil, errors.New("user id required")
	}
	uid, err := primitive.ObjectIDFromHex(userHex)
	if err != nil {
		return nil, err
	}
	if existing.UserID != uid {
		return nil, errors.New("forbidden: not the owner")
	}

	if !payload.CompletedAt.IsZero() {
		existing.CompletedAt = payload.CompletedAt
	}
	if payload.DurationMinutes != 0 {
		existing.DurationMinutes = payload.DurationMinutes
	}
	if payload.Notes != "" {
		existing.Notes = payload.Notes
	}
	if payload.RoutineID != &primitive.NilObjectID {
		existing.RoutineID = payload.RoutineID
	}
	if len(payload.PerformedExercises) > 0 {
		existing.PerformedExercises = payload.PerformedExercises
	}

	existing.UpdatedAt = time.Now()
	return s.repo.UpdateWorkout(existing)
}

func (s *WorkoutService) DeleteWorkout(idHex, userHex string) (*mongo.DeleteResult, error) {
	if idHex == "" {
		return nil, errors.New("id required")
	}
	existing, err := s.repo.GetWorkoutByID(idHex)
	if err != nil {
		return nil, err
	}
	if userHex == "" {
		return nil, errors.New("user id required")
	}
	uid, err := primitive.ObjectIDFromHex(userHex)
	if err != nil {
		return nil, err
	}
	if existing.UserID != uid {
		return nil, errors.New("forbidden: not the owner")
	}
	oid, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, err
	}
	return s.repo.DeleteWorkout(oid)
}

func (s *WorkoutService) GetWorkoutFrequency(userHex, period string, from, to time.Time) (map[string]int, error) {
	workouts, err := s.GetWorkouts(userHex)
	if err != nil {
		return nil, err
	}
	result := map[string]int{}
	for _, w := range workouts {
		if !from.IsZero() && w.CompletedAt.Before(from) {
			continue
		}
		if !to.IsZero() && w.CompletedAt.After(to) {
			continue
		}
		var key string
		switch period {
		case "weekly":
			year, week := w.CompletedAt.ISOWeek()
			key = formatWeekKey(year, week)
		case "monthly":
			key = w.CompletedAt.Format("2006-01")
		default:
			key = w.CompletedAt.Format("2006-01-02")
		}
		result[key]++
	}
	return result, nil
}

func (s *WorkoutService) GetTopRoutines(userHex string, limit int) ([]RoutineUsage, error) {
	workouts, err := s.GetWorkouts(userHex)
	if err != nil {
		return nil, err
	}
	counts := map[string]int{}
	for _, w := range workouts {
		if w.RoutineID.IsZero() {
			continue
		}
		counts[w.RoutineID.Hex()]++
	}
	usages := make([]RoutineUsage, 0, len(counts))
	for rid, c := range counts {
		usages = append(usages, RoutineUsage{RoutineID: rid, Count: c})
	}
	sort.Slice(usages, func(i, j int) bool { return usages[i].Count > usages[j].Count })
	if limit > 0 && len(usages) > limit {
		usages = usages[:limit]
	}
	return usages, nil
}

func (s *WorkoutService) GetProgressOverTime(userHex, metric string, from, to time.Time) ([]DataPoint, error) {
	workouts, err := s.GetWorkouts(userHex)
	if err != nil {
		return nil, err
	}
	byDay := map[string]float64{}
	loc := time.UTC
	for _, w := range workouts {
		if !from.IsZero() && w.CompletedAt.Before(from) {
			continue
		}
		if !to.IsZero() && w.CompletedAt.After(to) {
			continue
		}
		dayKey := w.CompletedAt.In(loc).Format("2006-01-02")
		switch metric {
		case "duration":
			byDay[dayKey] += float64(w.DurationMinutes)
		case "volume":
			var total float64
			for _, pe := range w.PerformedExercises {
				total += float64(pe.Sets) * float64(pe.Reps) * pe.Weight
			}
			byDay[dayKey] += total
		default: 
			byDay[dayKey] += 1
		}
	}

	keys := make([]string, 0, len(byDay))
	for k := range byDay {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	points := make([]DataPoint, 0, len(keys))
	for _, k := range keys {
		tm, _ := time.ParseInLocation("2006-01-02", k, loc)
		points = append(points, DataPoint{Timestamp: tm, Value: byDay[k]})
	}
	return points, nil
}

func formatWeekKey(year, week int) string {
	return strconv.Itoa(year) + "-W" + twoDigit(week)
}

func twoDigit(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}
