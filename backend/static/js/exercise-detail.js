// Exercise detail page functionality

let currentExerciseId = null;
let currentExercise = null;
let isAdmin = false;

async function loadExercise() {
    const exerciseId = window.location.pathname.split('/').pop();
    currentExerciseId = exerciseId;
    
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
        currentExercise = exercise;
        
        // Check if user is admin (from token or user data)
        checkUserRole();
        
        renderExercise(exercise);
    } catch (error) {
        document.getElementById('loadingSpinner').classList.add('d-none');
        document.getElementById('errorMessage').textContent = error.message || 'Error al cargar el ejercicio';
        document.getElementById('errorMessage').classList.remove('d-none');
    }
}

function checkUserRole() {
    // Try to get user role from localStorage
    const userStr = localStorage.getItem('user');
    if (userStr) {
        try {
            const user = JSON.parse(userStr);
            isAdmin = user.role === 'admin';
            if (isAdmin) {
                showAdminControls();
            }
        } catch (e) {
            // Ignore parse errors
        }
    }
}

function showAdminControls() {
    const actionButtons = document.getElementById('exerciseActions');
    if (actionButtons) {
        actionButtons.innerHTML = `
            <button class="btn btn-warning" onclick="editExercise()">
                <i class="bi bi-pencil"></i> Editar Ejercicio
            </button>
            <button class="btn btn-outline-danger" onclick="deleteExercise()">
                <i class="bi bi-trash"></i> Eliminar Ejercicio
            </button>
        `;
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
    if (currentExerciseId) {
        window.location.href = `/routines/create?exercise_id=${currentExerciseId}`;
    }
}

async function editExercise() {
    if (currentExerciseId) {
        window.location.href = `/exercises/${currentExerciseId}/edit`;
    }
}

async function deleteExercise() {
    if (!confirm('¿Estás seguro de eliminar este ejercicio? Esta acción no se puede deshacer.')) {
        return;
    }
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/exercises/${currentExerciseId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });

        if (response.ok) {
            alert('Ejercicio eliminado correctamente');
            window.location.href = '/exercises';
        } else {
            const data = await response.json();
            alert(data.error || 'Error al eliminar el ejercicio');
        }
    } catch (error) {
        alert('Error al eliminar el ejercicio');
        console.error('Error:', error);
    }
}

document.addEventListener('DOMContentLoaded', loadExercise);

