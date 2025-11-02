let currentWorkoutId = null;
async function loadWorkout() {
    const workoutId = window.location.pathname.split('/').pop();
    currentWorkoutId = workoutId;
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/workouts/${workoutId}`, {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });
        if (!response.ok) {
            throw new Error('Entrenamiento no encontrado');
        }
        const workout = await response.json();
        renderWorkout(workout);
    } catch (error) {
        document.getElementById('loadingSpinner').classList.add('d-none');
        document.getElementById('errorMessage').textContent = error.message || 'Error al cargar el entrenamiento';
        document.getElementById('errorMessage').classList.remove('d-none');
    }
}
function renderWorkout(workout) {
    document.getElementById('loadingSpinner').classList.add('d-none');
    document.getElementById('workoutContent').classList.remove('d-none');
    const completedAt = new Date(workout.completed_at);
    document.getElementById('workoutDate').innerHTML = 
        `<i class="bi bi-calendar"></i> ${completedAt.toLocaleDateString('es-ES')} ${completedAt.toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit' })}`;
    if (workout.routine_id) {
        const routineDiv = document.getElementById('workoutRoutine');
        routineDiv.innerHTML = `
            <strong>Rutina:</strong> ${workout.routine_name || 'Rutina asociada'}
            <a href="/routines/${workout.routine_id}" class="btn btn-sm btn-outline-primary ms-2">
                Ver Rutina
            </a>
        `;
        routineDiv.classList.remove('d-none');
    }
    if (workout.duration_minutes) {
        const durationDiv = document.getElementById('workoutDuration');
        durationDiv.innerHTML = `<i class="bi bi-clock"></i> <strong>Duración:</strong> ${workout.duration_minutes} minutos`;
        durationDiv.classList.remove('d-none');
    }
    if (workout.estimated_calories) {
        const caloriesDiv = document.getElementById('workoutCalories');
        caloriesDiv.innerHTML = `<i class="bi bi-fire"></i> <strong>Calorías estimadas:</strong> ${workout.estimated_calories} kcal`;
        caloriesDiv.classList.remove('d-none');
    }
    if (workout.notes) {
        document.getElementById('notesText').textContent = workout.notes;
        document.getElementById('workoutNotes').classList.remove('d-none');
    }
    const updatedAt = workout.updated_at ? new Date(workout.updated_at) : null;
    const info = document.getElementById('workoutInfo');
    let infoHTML = `
        <li class="mb-2">
            <strong>Fecha:</strong><br>
            ${completedAt.toLocaleDateString('es-ES')}
        </li>
        <li class="mb-2">
            <strong>Hora:</strong><br>
            ${completedAt.toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit' })}
        </li>
    `;
    if (updatedAt) {
        infoHTML += `
            <li class="mb-2">
                <strong>Última actualización:</strong><br>
                ${updatedAt.toLocaleDateString('es-ES')} ${updatedAt.toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit' })}
            </li>
        `;
    }
    info.innerHTML = infoHTML;
}
function editWorkout() {
    if (currentWorkoutId) {
        window.location.href = `/workouts/${currentWorkoutId}/edit`;
    }
}
async function deleteWorkout() {
    if (!confirm('¿Estás seguro de eliminar este entrenamiento? Esta acción no se puede deshacer.')) {
        return;
    }
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/workouts/${currentWorkoutId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });
        if (response.ok) {
            alert('Entrenamiento eliminado correctamente');
            window.location.href = '/workouts';
        } else {
            const data = await response.json();
            alert(data.error || 'Error al eliminar el entrenamiento');
        }
    } catch (error) {
        alert('Error al eliminar el entrenamiento');
        console.error('Error:', error);
    }
}
document.addEventListener('DOMContentLoaded', loadWorkout);
