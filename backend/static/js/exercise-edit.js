let exerciseId = null;
let currentExercise = null;
function addStep(stepText = '') {
    const container = document.getElementById('stepsContainer');
    const stepCount = container.children.length + 1;
    const stepDiv = document.createElement('div');
    stepDiv.className = 'input-group mb-2';
    stepDiv.innerHTML = `
        <input type="text" class="form-control step-input" placeholder="Paso ${stepCount}" value="${stepText}">
        <button type="button" class="btn btn-outline-danger" onclick="removeStep(this)">
            <i class="bi bi-trash"></i>
        </button>
    `;
    container.appendChild(stepDiv);
}
function removeStep(button) {
    button.parentElement.remove();
}
function showError(message) {
    const errorDiv = document.getElementById('errorMessage');
    errorDiv.textContent = message;
    errorDiv.classList.remove('d-none');
}
function hideError() {
    const errorDiv = document.getElementById('errorMessage');
    errorDiv.classList.add('d-none');
}
async function loadExercise() {
    exerciseId = window.location.pathname.split('/').slice(-2, -1)[0] || window.location.pathname.split('/').pop();
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
        currentExercise = await response.json();
        populateForm(currentExercise);
    } catch (error) {
        document.getElementById('loadingSpinner').classList.add('d-none');
        showError(error.message || 'Error al cargar el ejercicio');
    }
}
function populateForm(exercise) {
    document.getElementById('name').value = exercise.name || '';
    document.getElementById('description').value = exercise.description || '';
    document.getElementById('category').value = exercise.category || '';
    document.getElementById('muscle_group').value = exercise.muscle_group || '';
    document.getElementById('difficulty').value = exercise.difficulty || '';
    document.getElementById('media_url').value = exercise.media_url || '';
    const stepsContainer = document.getElementById('stepsContainer');
    stepsContainer.innerHTML = '';
    if (exercise.steps && exercise.steps.length > 0) {
        exercise.steps.forEach(step => {
            addStep(step);
        });
    } else {
        addStep();
    }
    document.getElementById('loadingSpinner').classList.add('d-none');
    document.getElementById('exerciseForm').classList.remove('d-none');
}
document.getElementById('exerciseForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    hideError();
    const steps = Array.from(document.querySelectorAll('.step-input'))
        .map(input => input.value)
        .filter(value => value.trim() !== '');
    const exerciseData = {
        name: document.getElementById('name').value,
        description: document.getElementById('description').value,
        category: document.getElementById('category').value,
        muscle_group: document.getElementById('muscle_group').value,
        difficulty: document.getElementById('difficulty').value,
        media_url: document.getElementById('media_url').value,
        steps: steps
    };
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/exercises/${exerciseId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
            body: JSON.stringify(exerciseData)
        });
        if (response.ok) {
            window.location.href = `/exercises/${exerciseId}`;
        } else {
            const error = await response.json();
            showError(error.error || 'Error al actualizar el ejercicio');
        }
    } catch (error) {
        showError('Error al actualizar el ejercicio');
        console.error('Error:', error);
    }
});
document.addEventListener('DOMContentLoaded', loadExercise);
