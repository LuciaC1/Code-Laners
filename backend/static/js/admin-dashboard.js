// Admin dashboard functionality

document.addEventListener('DOMContentLoaded', function() {
    if (!requireAuth()) {
        return;
    }
    
    // Check if user is admin
    checkAdminAccess();
    loadDashboardData();
});

let usersChart = null;
let popularExercisesChart = null;

async function checkAdminAccess() {
    // El middleware del backend verifica el rol, así que solo verificamos autenticación
    // Si el backend rechaza, manejará el error automáticamente
    if (!requireAuth()) {
        return;
    }
}

async function loadDashboardData() {
    try {
        const headers = getApiHeaders();
        
        // Load all statistics
        const [users, exercises, routines, workouts] = await Promise.all([
            fetch('/api/admin/users', { headers }).then(r => r.json()),
            fetch('/api/exercises', { headers }).then(r => r.json()),
            fetch('/api/routines', { headers }).then(r => r.json()),
            fetch('/api/workouts', { headers }).then(r => r.json())
        ]);

        const usersData = users.users || [];
        const exercisesData = exercises.exercises || [];
        const routinesData = routines.routines || [];
        const workoutsData = workouts.workouts || [];

        // Update summary cards
        updateSummaryCards(usersData.length, exercisesData.length, routinesData.length, workoutsData.length);

        // Render charts
        if (typeof Chart !== 'undefined') {
            renderUsersChart(usersData);
            renderPopularExercisesChart(exercisesData, workoutsData);
        } else {
            loadChartLibrary().then(() => {
                renderUsersChart(usersData);
                renderPopularExercisesChart(exercisesData, workoutsData);
            });
        }

        // Render popular routines and recent activity
        renderPopularRoutines(routinesData, workoutsData);
        renderRecentActivity(workoutsData, usersData);
    } catch (error) {
        console.error('Error loading dashboard data:', error);
    }
}

function updateSummaryCards(users, exercises, routines, workouts) {
    document.getElementById('totalUsers').textContent = users;
    document.getElementById('totalExercises').textContent = exercises;
    document.getElementById('totalRoutines').textContent = routines;
    document.getElementById('totalWorkouts').textContent = workouts;
}

function renderUsersChart(users) {
    const ctx = document.getElementById('usersChart');
    if (!ctx) return;

    // Group users by month
    const monthlyData = {};
    users.forEach(user => {
        const date = new Date(user.created_at);
        const monthKey = date.toLocaleDateString('es-ES', { year: 'numeric', month: 'short' });
        monthlyData[monthKey] = (monthlyData[monthKey] || 0) + 1;
    });

    const labels = Object.keys(monthlyData).sort();
    const data = labels.map(key => monthlyData[key]);

    if (usersChart) {
        usersChart.destroy();
    }

    usersChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Usuarios Registrados',
                data: data,
                backgroundColor: 'rgba(13, 110, 253, 0.6)',
                borderColor: 'rgba(13, 110, 253, 1)',
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            scales: {
                y: {
                    beginAtZero: true,
                    ticks: {
                        stepSize: 1
                    }
                }
            }
        }
    });
}

function renderPopularExercisesChart(exercises, workouts) {
    const ctx = document.getElementById('popularExercisesChart');
    if (!ctx) return;

    // Count exercise usage (if workouts have exercise references)
    const exerciseCounts = {};
    workouts.forEach(workout => {
        // This would require workout data to include exercise IDs
        // For now, we'll use a placeholder approach
    });

    // Fallback: show exercises by category
    const categoryCounts = {};
    exercises.forEach(exercise => {
        const category = exercise.category || 'Otros';
        categoryCounts[category] = (categoryCounts[category] || 0) + 1;
    });

    const labels = Object.keys(categoryCounts);
    const data = Object.values(categoryCounts);
    const colors = generateColors(labels.length);

    if (popularExercisesChart) {
        popularExercisesChart.destroy();
    }

    popularExercisesChart = new Chart(ctx, {
        type: 'doughnut',
        data: {
            labels: labels,
            datasets: [{
                data: data,
                backgroundColor: colors
            }]
        },
        options: {
            responsive: true
        }
    });
}

function renderPopularRoutines(routines, workouts) {
    const container = document.getElementById('popularRoutinesList');
    if (!container) return;

    // Count workouts by routine
    const routineCounts = {};
    workouts.forEach(workout => {
        const routineId = workout.routine_id;
        if (routineId) {
            routineCounts[routineId] = (routineCounts[routineId] || 0) + 1;
        }
    });

    // Sort routines by usage
    const sortedRoutines = routines.map(r => ({
        ...r,
        usage: routineCounts[r.id] || 0
    })).sort((a, b) => b.usage - a.usage).slice(0, 5);

    if (sortedRoutines.length === 0) {
        container.innerHTML = '<p class="text-muted text-center">No hay rutinas aún</p>';
        return;
    }

    container.innerHTML = sortedRoutines.map(routine => `
        <div class="d-flex justify-content-between align-items-center mb-3 p-2 border rounded">
            <div>
                <strong>${routine.name}</strong>
                <br>
                <small class="text-muted">${routine.usage} entrenamientos</small>
            </div>
            <span class="badge bg-primary">${routine.exercises?.length || 0} ejercicios</span>
        </div>
    `).join('');
}

function renderRecentActivity(workouts, users) {
    const container = document.getElementById('recentActivityList');
    if (!container) return;

    // Combine recent workouts and new users
    const activities = [];

    // Recent workouts (last 5)
    workouts.slice(0, 5).forEach(workout => {
        activities.push({
            type: 'workout',
            date: new Date(workout.completed_at),
            text: `Entrenamiento realizado: ${workout.routine_name || 'Sin rutina'}`
        });
    });

    // Recent users (last 5)
    users.slice(0, 5).forEach(user => {
        activities.push({
            type: 'user',
            date: new Date(user.created_at),
            text: `Nuevo usuario: ${user.name}`
        });
    });

    // Sort by date (most recent first)
    activities.sort((a, b) => b.date - a.date).slice(0, 10);

    if (activities.length === 0) {
        container.innerHTML = '<p class="text-muted text-center">No hay actividad reciente</p>';
        return;
    }

    container.innerHTML = activities.map(activity => {
        const icon = activity.type === 'workout' ? 'bi-calendar-check' : 'bi-person-plus';
        const color = activity.type === 'workout' ? 'text-success' : 'text-primary';
        const dateStr = activity.date.toLocaleDateString('es-ES', {
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });

        return `
            <div class="d-flex align-items-start mb-2">
                <i class="bi ${icon} ${color} me-2 mt-1"></i>
                <div class="flex-grow-1">
                    <small class="text-muted">${dateStr}</small>
                    <p class="mb-0">${activity.text}</p>
                </div>
            </div>
        `;
    }).join('');
}

function generateColors(count) {
    const colors = [
        'rgba(13, 110, 253, 0.6)',
        'rgba(25, 135, 84, 0.6)',
        'rgba(255, 193, 7, 0.6)',
        'rgba(220, 53, 69, 0.6)',
        'rgba(13, 202, 240, 0.6)',
        'rgba(111, 66, 193, 0.6)',
        'rgba(214, 51, 132, 0.6)'
    ];
    
    const result = [];
    for (let i = 0; i < count; i++) {
        result.push(colors[i % colors.length]);
    }
    return result;
}

function loadChartLibrary() {
    return new Promise((resolve, reject) => {
        if (typeof Chart !== 'undefined') {
            resolve();
            return;
        }

        const script = document.createElement('script');
        script.src = 'https://cdn.jsdelivr.net/npm/chart.js@4.4.0/dist/chart.umd.min.js';
        script.onload = resolve;
        script.onerror = reject;
        document.head.appendChild(script);
    });
}

