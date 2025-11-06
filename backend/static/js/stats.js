document.addEventListener('DOMContentLoaded', function() {
    if (!requireAuth()) {
        return;
    }
    loadStatistics();
});
let frequencyChart = null;
let routinesChart = null;
let progressChart = null;
async function loadStatistics() {
    try {
        const headers = getApiHeaders();
        const statsResponse = await fetch('/api/stats', { headers });
        if (!statsResponse.ok) {
            if (statsResponse.status === 401) {
                window.location.href = '/login';
                return;
            }
            throw new Error(`Error al cargar estadísticas: ${statsResponse.status}`);
        }
        const statsData = await statsResponse.json();
        
        const summary = statsData.summary || {};
        document.getElementById('totalWorkouts').textContent = summary.totalWorkouts || 0;
        document.getElementById('totalRoutines').textContent = summary.totalRoutines || 0;
        document.getElementById('totalMinutes').textContent = summary.totalMinutes || 0;
        document.getElementById('totalCalories').textContent = summary.totalCalories || 0;
        const workouts = statsData.recentWorkouts || [];
        
        if (typeof Chart === 'undefined') {
            loadChartLibrary().then(() => {
                renderFrequencyChartFromData(statsData.frequency || []);
                renderRoutinesChartFromData(statsData.routinesDistribution || []);
                renderProgressChartFromData(statsData.progress || []);
            });
        } else {
            renderFrequencyChartFromData(statsData.frequency || []);
            renderRoutinesChartFromData(statsData.routinesDistribution || []);
            renderProgressChartFromData(statsData.progress || []);
        }
        renderRecentWorkouts(workouts);
    } catch (error) {
        console.error('Error loading statistics:', error);
        const tbody = document.getElementById('recentWorkoutsTable');
        if (tbody) {
            tbody.innerHTML = `<tr><td colspan="5" class="text-center text-danger">Error al cargar estadísticas: ${error.message}</td></tr>`;
        }
    }
}
function calculateSummaryStats(workouts, routines) {
    const totalWorkouts = workouts.length;
    const totalRoutines = routines.length;
    const totalMinutes = workouts.reduce((sum, w) => sum + (w.duration_minutes || 0), 0);
    const totalCalories = workouts.reduce((sum, w) => sum + (w.estimated_calories || 0), 0);
    document.getElementById('totalWorkouts').textContent = totalWorkouts;
    document.getElementById('totalRoutines').textContent = totalRoutines;
    document.getElementById('totalMinutes').textContent = totalMinutes;
    document.getElementById('totalCalories').textContent = totalCalories;
}
function renderCharts(workouts, routines) {
    renderFrequencyChart(workouts);
    renderRoutinesChart(workouts);
    renderProgressChart(workouts);
}
function renderFrequencyChart(workouts) {
    const ctx = document.getElementById('frequencyChart');
    if (!ctx) return;
    const weeklyData = {};
    workouts.forEach(workout => {
        const date = new Date(workout.completed_at);
        const weekKey = getWeekKey(date);
        weeklyData[weekKey] = (weeklyData[weekKey] || 0) + 1;
    });
    const labels = Object.keys(weeklyData).sort();
    const data = labels.map(key => weeklyData[key]);
    if (frequencyChart) {
        frequencyChart.destroy();
    }
    frequencyChart = new Chart(ctx, {
        type: 'bar',
        data: {
            labels: labels,
            datasets: [{
                label: 'Entrenamientos por Semana',
                data: data,
                backgroundColor: 'rgba(13, 110, 253, 0.6)',
                borderColor: 'rgba(13, 110, 253, 1)',
                borderWidth: 1
            }]
        },
        options: {
            responsive: true,
            maintainAspectRatio: false,
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

function renderFrequencyChartFromData(frequency) {
    const ctx = document.getElementById('frequencyChart');
    if (!ctx) return;
    const labels = frequency.map(f => f.label);
    const data = frequency.map(f => f.count);
    if (frequencyChart) { frequencyChart.destroy(); }
    frequencyChart = new Chart(ctx, {
        type: 'bar',
        data: { labels: labels, datasets: [{ label: 'Entrenamientos por Semana', data: data, backgroundColor: 'rgba(13, 110, 253, 0.6)', borderColor: 'rgba(13, 110, 253, 1)', borderWidth: 1 }] },
        options: { responsive: true, maintainAspectRatio: false, scales: { y: { beginAtZero: true, ticks: { stepSize: 1 } } } }
    });
}
function renderRoutinesChart(workouts) {

}

function renderRoutinesChartFromData(routinesDistribution) {
    const ctx = document.getElementById('routinesChart');
    if (!ctx) return;
    const labels = routinesDistribution.map(r => r.routine_name || r.routineName || 'Sin rutina');
    const data = routinesDistribution.map(r => r.count || 0);
    const colors = generateColors(labels.length);
    if (routinesChart) { routinesChart.destroy(); }
    routinesChart = new Chart(ctx, { type: 'doughnut', data: { labels: labels, datasets: [{ data: data, backgroundColor: colors }] }, options: { responsive: true, maintainAspectRatio: false } });

}
function renderProgressChart(workouts) {
}

function renderProgressChartFromData(points) {
    const ctx = document.getElementById('progressChart');
    if (!ctx) return;
    const sorted = [...points].sort((a,b) => new Date(a.date) - new Date(b.date));
    const labels = sorted.map(p => new Date(p.date).toLocaleDateString('es-ES', { month: 'short', day: 'numeric' }));
    const caloriesData = sorted.map(p => p.calories || 0);
    const durationData = sorted.map(p => p.duration || 0);
    if (progressChart) { progressChart.destroy(); }
    progressChart = new Chart(ctx, {
        type: 'line',
        data: { labels: labels, datasets: [ { label: 'Calorías Quemadas', data: caloriesData, borderColor: 'rgba(220, 53, 69, 1)', backgroundColor: 'rgba(220, 53, 69, 0.1)', yAxisID: 'y' }, { label: 'Duración (min)', data: durationData, borderColor: 'rgba(13, 110, 253, 1)', backgroundColor: 'rgba(13, 110, 253, 0.1)', yAxisID: 'y1' } ] },
        options: { responsive: true, interaction: { mode: 'index', intersect: false }, scales: { y: { type: 'linear', display: true, position: 'left' }, y1: { type: 'linear', display: true, position: 'right', grid: { drawOnChartArea: false } } } }
    });
}
function renderRecentWorkouts(workouts) {
    const tbody = document.getElementById('recentWorkoutsTable');
    if (!tbody) return;
    if (workouts.length === 0) {
        tbody.innerHTML = '<tr><td colspan="5" class="text-center text-muted">No hay entrenamientos registrados</td></tr>';
        return;
    }
    tbody.innerHTML = workouts.map(workout => {
        const date = new Date(workout.completed_at);
        const formattedDate = date.toLocaleDateString('es-ES', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
        return `
            <tr>
                <td>${formattedDate}</td>
                <td>${workout.routine_name || 'Sin rutina'}</td>
                <td>${workout.duration_minutes || 0} min</td>
                <td>${workout.estimated_calories || 0} cal</td>
                <td>
                    <a href="/workouts/${workout.id}" class="btn btn-sm btn-outline-primary">
                        <i class="bi bi-eye"></i> Ver
                    </a>
                </td>
            </tr>
        `;
    }).join('');
}
function getWeekKey(date) {
    const year = date.getFullYear();
    const week = getWeekNumber(date);
    return `Semana ${week}, ${year}`;
}
function getWeekNumber(date) {
    const d = new Date(Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()));
    const dayNum = d.getUTCDay() || 7;
    d.setUTCDate(d.getUTCDate() + 4 - dayNum);
    const yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1));
    return Math.ceil((((d - yearStart) / 86400000) + 1) / 7);
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
