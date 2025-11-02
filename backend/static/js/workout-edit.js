

let workoutId = null;
let currentWorkout = null;

async function loadWorkout() {
    workoutId = window.location.pathname.split('/').slice(-2, -1)[0] || window.location.pathname.split('/').pop();
    
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

        currentWorkout = await response.json();
        await loadRoutines();
        populateForm(currentWorkout);
    } catch (error) {
        document.getElementById('loadingSpinner').classList.add('d-none');
        showError(error.message || 'Error al cargar el entrenamiento');
    }
}

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
        
        (data.routines || []).forEach(routine => {
            const option = document.createElement('option');
            option.value = routine.id;
            option.textContent = routine.name;
            select.appendChild(option);
        });
        
        
        if (currentWorkout && currentWorkout.routine_id) {
            select.value = currentWorkout.routine_id;
        }
    } catch (error) {
        console.error('Error loading routines:', error);
    }
}

function populateForm(workout) {
    
    if (workout.routine_id) {
        document.getElementById('routine_id').value = workout.routine_id;
    }
    
    
    if (workout.duration_minutes) {
        document.getElementById('duration_minutes').value = workout.duration_minutes;
    }
    if (workout.estimated_calories) {
        document.getElementById('estimated_calories').value = workout.estimated_calories;
    }
    
    
    if (workout.completed_at) {
        const completedAt = new Date(workout.completed_at);
        const dateStr = completedAt.toISOString().split('T')[0];
        const timeStr = completedAt.toTimeString().split(' ')[0].substring(0, 5);
        document.getElementById('completed_at_date').value = dateStr;
        document.getElementById('completed_at_time').value = timeStr;
    }
    
    
    if (workout.notes) {
        document.getElementById('notes').value = workout.notes;
    }
    
    
    document.getElementById('loadingSpinner').classList.add('d-none');
    document.getElementById('workoutForm').classList.remove('d-none');
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

document.getElementById('workoutForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    hideError();
    
    
    const dateStr = document.getElementById('completed_at_date').value;
    const timeStr = document.getElementById('completed_at_time').value;
    const completedAt = new Date(`${dateStr}T${timeStr}`);
    
    const workoutData = {
        routine_id: document.getElementById('routine_id').value || undefined,
        duration_minutes: parseInt(document.getElementById('duration_minutes').value) || undefined,
        estimated_calories: parseInt(document.getElementById('estimated_calories').value) || undefined,
        notes: document.getElementById('notes').value || undefined,
        completed_at: completedAt.toISOString()
    };

    
    Object.keys(workoutData).forEach(key => {
        if (workoutData[key] === undefined) {
            delete workoutData[key];
        }
    });

    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/workouts/${workoutId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
            body: JSON.stringify(workoutData)
        });

        if (response.ok) {
            window.location.href = `/workouts/${workoutId}`;
        } else {
            const error = await response.json();
            showError(error.error || 'Error al actualizar el entrenamiento');
        }
    } catch (error) {
        showError('Error al actualizar el entrenamiento');
        console.error('Error:', error);
    }
});

document.addEventListener('DOMContentLoaded', loadWorkout);

