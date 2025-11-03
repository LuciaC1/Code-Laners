async function loadExercises(params = null) {
    try {
        const token = localStorage.getItem('token');
        let url = '/api/exercises';
        if (params) {
            const queryString = params.toString();
            if (queryString) {
                url += '?' + queryString;
            }
        }
        const response = await fetch(url, {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });
        const data = await response.json();
        renderExercises(data.exercises || []);
    } catch (error) {
        console.error('Error loading exercises:', error);
        const container = document.getElementById('exercisesContainer');
        container.innerHTML = '<div class="col-12"><div class="alert alert-danger">Error al cargar ejercicios</div></div>';
    }
}
function renderExercises(exercises) {
    const container = document.getElementById('exercisesContainer');
    container.innerHTML = '';
    if (exercises.length === 0) {
        container.innerHTML = '<div class="col-12"><div class="alert alert-info">No se encontraron ejercicios</div></div>';
        return;
    }
    exercises.forEach(exercise => {
        const card = document.createElement('div');
        card.className = 'col-md-4';
        card.innerHTML = `
            <div class="card h-100 shadow-sm">
                <div class="card-body">
                    <h5 class="card-title">${exercise.name}</h5>
                    <p class="card-text text-muted">${exercise.description || 'Sin descripción'}</p>
                    <div class="mb-2">
                        <span class="badge bg-primary">${exercise.category}</span>
                        <span class="badge bg-success">${exercise.muscle_group}</span>
                        <span class="badge bg-warning text-dark">${exercise.difficulty}</span>
                    </div>
                    <a href="/exercises/${exercise.id}" class="btn btn-primary btn-sm">
                        <i class="bi bi-eye"></i> Ver Detalles
                    </a>
                </div>
            </div>
        `;
        container.appendChild(card);
    });
}
function applyFilters() {
    const name = document.getElementById('filterName').value;
    const category = document.getElementById('filterCategory').value;
    const muscleGroup = document.getElementById('filterMuscleGroup').value;
    const params = new URLSearchParams();
    if (name) params.append('name', name);
    if (category) params.append('category', category);
    if (muscleGroup) params.append('muscle_group', muscleGroup);
    loadExercises(params);
}
document.addEventListener('DOMContentLoaded', loadExercises);
