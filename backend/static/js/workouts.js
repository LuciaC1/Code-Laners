// Workouts list page functionality

async function loadWorkouts() {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/workouts', {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });
        const data = await response.json();
        renderWorkouts(data || []);
    } catch (error) {
        console.error('Error loading workouts:', error);
        document.getElementById('workoutsContainer').innerHTML = 
            '<div class="col-12"><div class="alert alert-danger">Error al cargar entrenamientos</div></div>';
    }
}

function renderWorkouts(workouts) {
    const container = document.getElementById('workoutsContainer');
    container.innerHTML = '';
    
    if (workouts.length === 0) {
        container.innerHTML = `
            <div class="col-12">
                <div class="alert alert-info">
                    <i class="bi bi-info-circle"></i> No tienes entrenamientos registrados. 
                    <a href="/workouts/create">Registra tu primer entrenamiento</a>
                </div>
            </div>
        `;
        return;
    }

    workouts.forEach(workout => {
        const card = document.createElement('div');
        card.className = 'col-md-6 col-lg-4';
        const date = new Date(workout.completed_at);
        card.innerHTML = `
            <div class="card h-100 shadow-sm">
                <div class="card-body">
                    <h5 class="card-title">Entrenamiento</h5>
                    <p class="text-muted mb-2">
                        <i class="bi bi-calendar"></i> ${date.toLocaleDateString()}
                    </p>
                    ${workout.duration_minutes ? `
                        <p class="mb-1"><i class="bi bi-clock"></i> ${workout.duration_minutes} minutos</p>
                    ` : ''}
                    ${workout.estimated_calories ? `
                        <p class="mb-1"><i class="bi bi-fire"></i> ${workout.estimated_calories} calorías</p>
                    ` : ''}
                    ${workout.notes ? `
                        <p class="card-text">${workout.notes}</p>
                    ` : ''}
                </div>
                <div class="card-footer bg-transparent">
                    <a href="/workouts/${workout.id}" class="btn btn-primary btn-sm">
                        <i class="bi bi-eye"></i> Ver Detalles
                    </a>
                    <button class="btn btn-outline-danger btn-sm" onclick="deleteWorkout('${workout.id}')">
                        <i class="bi bi-trash"></i>
                    </button>
                </div>
            </div>
        `;
        container.appendChild(card);
    });
}

async function deleteWorkout(id) {
    if (!confirm('¿Estás seguro de eliminar este entrenamiento?')) return;

    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/workouts/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });

        if (response.ok) {
            loadWorkouts();
        } else {
            alert('Error al eliminar el entrenamiento');
        }
    } catch (error) {
        alert('Error al eliminar el entrenamiento');
    }
}

document.addEventListener('DOMContentLoaded', loadWorkouts);

