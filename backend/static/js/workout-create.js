async function loadRoutines() {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/routines', {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });
        const data = await response.json();
        const select = document.getElementById('routine_id');
        const urlParams = new URLSearchParams(window.location.search);
        const routineId = urlParams.get('routine_id');
        (data.routines || []).forEach(routine => {
            const option = document.createElement('option');
            option.value = routine.id;
            option.textContent = routine.name;
            if (routineId && routine.id === routineId) {
                option.selected = true;
            }
            select.appendChild(option);
        });
    } catch (error) {
        console.error('Error loading routines:', error);
    }
}
function showError(message) {
    const errorDiv = document.getElementById('errorMessage');
    errorDiv.textContent = message;
    errorDiv.classList.remove('d-none');
}
document.getElementById('workoutForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const workoutData = {
        routine_id: document.getElementById('routine_id').value || undefined,
        duration_minutes: parseInt(document.getElementById('duration_minutes').value) || undefined,
        estimated_calories: parseInt(document.getElementById('estimated_calories').value) || undefined,
        notes: document.getElementById('notes').value || undefined
    };
    Object.keys(workoutData).forEach(key => {
        if (workoutData[key] === undefined) {
            delete workoutData[key];
        }
    });
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/workouts', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
            body: JSON.stringify(workoutData)
        });
        if (response.ok) {
            const data = await response.json();
            window.location.href = `/workouts/${data.id}`;
        } else {
            const error = await response.json();
            showError(error.error || 'Error al registrar el entrenamiento');
        }
    } catch (error) {
        showError('Error al registrar el entrenamiento');
    }
});
document.addEventListener('DOMContentLoaded', loadRoutines);
