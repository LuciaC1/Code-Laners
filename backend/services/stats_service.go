package services

import (
	"backend/dto"
	"backend/repositories"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StatsServiceInterface interface {
	GetUserStats(ctx context.Context, userID string) (*dto.UserStatsResponse, error)
	GetAdminStats(ctx context.Context) (*dto.AdminStatsResponse, error)
}

type StatsService struct {
	statsRepo    *repositories.StatsRepository
	routineRepo  repositories.RoutineRepositoryInterface
	userRepo     repositories.UserRepositoryInterface
	exerciseRepo repositories.ExerciseRepositoryInterface
}

func NewStatsService(statsRepo *repositories.StatsRepository, routineRepo repositories.RoutineRepositoryInterface, userRepo repositories.UserRepositoryInterface, exerciseRepo repositories.ExerciseRepositoryInterface) *StatsService {
	return &StatsService{statsRepo: statsRepo, routineRepo: routineRepo, userRepo: userRepo, exerciseRepo: exerciseRepo}
}

func (s *StatsService) GetUserStats(ctx context.Context, userID string) (*dto.UserStatsResponse, error) {
	uid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	resp := &dto.UserStatsResponse{}
	// totals
	totalWorkouts, err := s.statsRepo.CountWorkoutsByUser(uid)
	if err != nil {
		return nil, err
	}
	// routines accessible to user (owner or public) reuse RoutineRepository
	routines, err := s.routineRepo.GetRoutines(uid, "")
	if err != nil {
		return nil, err
	}
	totalRoutines := int64(len(routines))
	sumMinutes, sumCalories, err := s.statsRepo.SumWorkoutFieldsByUser(uid)
	if err != nil {
		return nil, err
	}
	resp.Summary.TotalWorkouts = totalWorkouts
	resp.Summary.TotalRoutines = totalRoutines
	resp.Summary.TotalMinutes = sumMinutes
	resp.Summary.TotalCalories = sumCalories

	// frequency
	weeks, err := s.statsRepo.AggregateWorkoutsByWeek(uid)
	if err != nil {
		return nil, err
	}
	for _, w := range weeks {
		idm := w["_id"].(bson.M)
		year := int(idm["year"].(int32))
		week := int(idm["week"].(int32))
		count := int64(0)
		if c, ok := w["count"].(int32); ok {
			count = int64(c)
		} else if c64, ok := w["count"].(int64); ok {
			count = c64
		}
		label := fmt.Sprintf("Semana %d, %d", week, year)
		resp.Frequency = append(resp.Frequency, dto.WeeklyCount{Label: label, Count: count})
	}

	// routines distribution
	rd, err := s.statsRepo.AggregateWorkoutsByRoutine(uid)
	if err != nil {
		return nil, err
	}
	for _, r := range rd {
		var name string
		if rn, ok := r["routine_name"].(string); ok && rn != "" {
			name = rn
		} else {
			name = "Sin rutina"
		}
		var count int64
		if c, ok := r["count"].(int32); ok {
			count = int64(c)
		} else if c64, ok := r["count"].(int64); ok {
			count = c64
		}
		var rid string
		if idv, ok := r["routine_id"].(primitive.ObjectID); ok {
			rid = idv.Hex()
		}
		resp.RoutinesDistribution = append(resp.RoutinesDistribution, dto.RoutineCount{RoutineID: rid, RoutineName: name, Count: count})
	}

	// progress
	workoutsProgress, err := s.statsRepo.GetProgressPoints(uid, 0)
	if err != nil {
		return nil, err
	}
	for _, w := range workoutsProgress {
		resp.Progress = append(resp.Progress, dto.ProgressPoint{Date: w.CompletedAt, Calories: w.EstimatedCalories, Duration: w.DurationMinutes})
	}

	// recent workouts
	recent, err := s.statsRepo.GetRecentWorkouts(uid, 10)
	if err != nil {
		return nil, err
	}
	for _, w := range recent {
		var ridName string
		if !w.RoutineID.IsZero() {
			// try to load routine name
			r, err := s.routineRepo.GetRoutineByID(w.RoutineID.Hex())
			if err == nil {
				ridName = r.Name
			}
		}
		resp.RecentWorkouts = append(resp.RecentWorkouts, dto.RecentWorkout{ID: w.ID.Hex(), CompletedAt: w.CompletedAt, RoutineName: ridName, DurationMinutes: w.DurationMinutes, EstimatedCalories: w.EstimatedCalories})
	}

	return resp, nil
}

func (s *StatsService) GetAdminStats(ctx context.Context) (*dto.AdminStatsResponse, error) {
	resp := &dto.AdminStatsResponse{}
	totals := dto.AdminTotals{}
	var err error
	totals.TotalUsers, err = s.statsRepo.CountCollection("users")
	if err != nil {
		return nil, err
	}
	totals.TotalExercises, err = s.statsRepo.CountCollection("exercises")
	if err != nil {
		return nil, err
	}
	totals.TotalRoutines, err = s.statsRepo.CountCollection("routines")
	if err != nil {
		return nil, err
	}
	totals.TotalWorkouts, err = s.statsRepo.CountCollection("workouts")
	if err != nil {
		return nil, err
	}
	resp.Totals = totals

	// users by month
	usersMonths, err := s.statsRepo.UsersByMonth()
	if err != nil {
		return nil, err
	}
	for _, u := range usersMonths {
		idm := u["_id"].(bson.M)
		year := int(idm["year"].(int32))
		month := int(idm["month"].(int32))
		label := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).Format("Jan 2006")
		var count int64
		if c, ok := u["count"].(int32); ok {
			count = int64(c)
		} else if c64, ok := u["count"].(int64); ok {
			count = c64
		}
		resp.UsersByMonth = append(resp.UsersByMonth, dto.WeeklyCount{Label: label, Count: count})
	}

	// exercises by category
	exCats, err := s.statsRepo.ExercisesByCategory()
	if err != nil {
		return nil, err
	}
	for _, e := range exCats {
		var category string
		if idv, ok := e["_id"].(string); ok {
			category = idv
		} else {
			category = "Otros"
		}
		var count int64
		if c, ok := e["count"].(int32); ok {
			count = int64(c)
		} else if c64, ok := e["count"].(int64); ok {
			count = c64
		}
		resp.ExercisesByCategory = append(resp.ExercisesByCategory, dto.RoutineCount{RoutineName: category, Count: count})
	}

	// popular routines
	pr, err := s.statsRepo.PopularRoutines(5)
	if err != nil {
		return nil, err
	}
	for _, p := range pr {
		var name string
		if n, ok := p["name"].(string); ok && n != "" {
			name = n
		} else {
			name = "Sin rutina"
		}
		var usage int64
		if u32, ok := p["usage"].(int32); ok {
			usage = int64(u32)
		} else if u64, ok := p["usage"].(int64); ok {
			usage = u64
		}
		var rid string
		if idv, ok := p["routine_id"].(primitive.ObjectID); ok {
			rid = idv.Hex()
		}
		resp.PopularRoutines = append(resp.PopularRoutines, dto.RoutineCount{RoutineID: rid, RoutineName: name, Count: usage})
	}

	// recent activity: merge latest workouts and latest users
	recentWorkouts, _ := s.statsRepo.RecentWorkoutsGlobal(5)
	recentUsers, _ := s.userRepo.GetUser("")
	var activities []map[string]interface{}
	for _, w := range recentWorkouts {
		activities = append(activities, map[string]interface{}{"type": "workout", "date": w.CompletedAt, "text": fmt.Sprintf("Entrenamiento: %s", w.ID.Hex())})
	}
	// take newest 5 users by created_at
	for i := 0; i < len(recentUsers) && i < 5; i++ {
		u := recentUsers[i]
		activities = append(activities, map[string]interface{}{"type": "user", "date": u.CreatedAt, "text": fmt.Sprintf("Nuevo usuario: %s", u.Name)})
	}
	resp.RecentActivity = activities

	return resp, nil
}
