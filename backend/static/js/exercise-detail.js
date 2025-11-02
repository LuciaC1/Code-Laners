// Exercise detail page functionality

async function loadExercise() {
    const exerciseId = window.location.pathname.split('/').pop();
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/exercises/${exerciseId}`, {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });

        if (!response.ok) {
            throw new Error('Ejercicio no encontrado');
        }

        const exercise = await response.json();
        renderExercise(exercise);
    } catch (error) {
        document.getElementById('loadingSpinner').classList.add('d-none');
        document.getElementById('errorMessage').textContent = error.message || 'Error al cargar el ejercicio';
        document.getElementById('errorMessage').classList.remove('d-none');
    }
}

function renderExercise(exercise) {
    document.getElementById('loadingSpinner').classList.add('d-none');
    document.getElementById('exerciseContent').classList.remove('d-none');

    document.getElementById('exerciseName').textContent = exercise.name || 'Sin nombre';
    document.getElementById('exerciseDescription').textContent = exercise.description || 'Sin descripción';

    const badges = document.getElementById('exerciseBadges');
    badges.innerHTML = `
        <span class="badge bg-primary fs-6">${exercise.category || 'N/A'}</span>
        <span class="badge bg-success fs-6">${exercise.muscle_group || 'N/A'}</span>
        <span class="badge bg-warning text-dark fs-6">${exercise.difficulty || 'N/A'}</span>
    `;

    const info = document.getElementById('exerciseInfo');
    info.innerHTML = `
        <li class="mb-2"><strong>Categoría:</strong> ${exercise.category || 'N/A'}</li>
        <li class="mb-2"><strong>Grupo Muscular:</strong> ${exercise.muscle_group || 'N/A'}</li>
        <li class="mb-2"><strong>Dificultad:</strong> ${exercise.difficulty || 'N/A'}</li>
    `;

    if (exercise.steps && exercise.steps.length > 0) {
        const stepsList = document.getElementById('stepsList');
        stepsList.innerHTML = exercise.steps.map(step => `<li>${step}</li>`).join('');
        document.getElementById('exerciseSteps').classList.remove('d-none');
    }

    if (exercise.media_url) {
        document.getElementById('exerciseImage').src = exercise.media_url;
        document.getElementById('exerciseImage').alt = exercise.name;
        document.getElementById('exerciseMedia').classList.remove('d-none');
    }
}

function addToRoutine() {
    alert('Funcionalidad de agregar a rutina próximamente');
}

document.addEventListener('DOMContentLoaded', loadExercise);

