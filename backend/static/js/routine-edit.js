

let routineId = null;
let selectedExercises = [];
let exerciseOrder = 0;
let currentRoutine = null;

async function loadRoutine() {
    routineId = window.location.pathname.split('/').slice(-2, -1)[0] || window.location.pathname.split('/').pop();
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/routines/${routineId}`, {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });

        if (!response.ok) {
            throw new Error('Rutina no encontrada');
        }

        currentRoutine = await response.json();
        populateForm(currentRoutine);
    } catch (error) {
        document.getElementById('loadingSpinner').classList.add('d-none');
        showError(error.message || 'Error al cargar la rutina');
    }
}

function populateForm(routine) {
    document.getElementById('name').value = routine.name || '';
    document.getElementById('description').value = routine.description || '';
    document.getElementById('is_public').checked = routine.is_public || false;

    
    selectedExercises = [];
    exerciseOrder = 0;
    
    if (routine.exercises && routine.exercises.length > 0) {
        routine.exercises.forEach(ex => {
            exerciseOrder++;
            selectedExercises.push({
                exercise_id: ex.exercise_id,
                exercise_name: ex.exercise_name || `Ejercicio ${exerciseOrder}`,
                order: ex.order || exerciseOrder,
                sets: ex.sets || 3,
                reps: ex.reps || 10,
                weight: ex.weight || 0
            });
        });
        exerciseOrder = selectedExercises.length;
    }
    
    renderRoutineExercises();
    
    
    document.getElementById('loadingSpinner').classList.add('d-none');
    document.getElementById('routineForm').classList.remove('d-none');
}

async function loadExerciseSelector() {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/exercises', {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });
        const data = await response.json();
        renderAvailableExercises(data.exercises || []);
        
        const modal = new bootstrap.Modal(document.getElementById('exerciseSelectorModal'));
        modal.show();
    } catch (error) {
        console.error('Error loading exercises:', error);
    }
}

function renderAvailableExercises(exercises) {
    const container = document.getElementById('availableExercises');
    container.innerHTML = '';
    
    exercises.forEach(exercise => {
        const item = document.createElement('a');
        item.href = '#';
        item.className = 'list-group-item list-group-item-action';
        item.innerHTML = `
            <div class="d-flex w-100 justify-content-between">
                <h6 class="mb-1">${exercise.name}</h6>
                <small>${exercise.category} - ${exercise.muscle_group}</small>
            </div>
            <p class="mb-1">${exercise.description || ''}</p>
        `;
        item.onclick = (e) => {
            e.preventDefault();
            addExerciseToRoutine(exercise);
            bootstrap.Modal.getInstance(document.getElementById('exerciseSelectorModal')).hide();
        };
        container.appendChild(item);
    });
}

function addExerciseToRoutine(exercise) {
    exerciseOrder++;
    const exerciseItem = {
        exercise_id: exercise.id,
        exercise_name: exercise.name,
        order: exerciseOrder,
        sets: 3,
        reps: 10,
        weight: 0
    };
    selectedExercises.push(exerciseItem);
    renderRoutineExercises();
}

function renderRoutineExercises() {
    const container = document.getElementById('routineExercises');
    container.innerHTML = '<h6 class="mt-3">Ejercicios Seleccionados:</h6>';
    
    if (selectedExercises.length === 0) {
        container.innerHTML += '<p class="text-muted">No hay ejercicios agregados</p>';
        return;
    }

    const table = document.createElement('table');
    table.className = 'table table-bordered';
    table.innerHTML = `
        <thead>
            <tr>
                <th>Orden</th>
                <th>Ejercicio</th>
                <th>Series</th>
                <th>Reps</th>
                <th>Peso (kg)</th>
                <th>Acciones</th>
            </tr>
        </thead>
        <tbody id="routineExercisesBody"></tbody>
    `;
    container.appendChild(table);

    const tbody = document.getElementById('routineExercisesBody');
    selectedExercises.forEach((item, index) => {
        const row = document.createElement('tr');
        row.innerHTML = `
            <td>${item.order}</td>
            <td>${item.exercise_name || `Ejercicio ${item.order}`}</td>
            <td><input type="number" class="form-control form-control-sm" value="${item.sets}" 
                onchange="updateExercise(${index}, 'sets', this.value)"></td>
            <td><input type="number" class="form-control form-control-sm" value="${item.reps}" 
                onchange="updateExercise(${index}, 'reps', this.value)"></td>
            <td><input type="number" class="form-control form-control-sm" value="${item.weight}" step="0.5"
                onchange="updateExercise(${index}, 'weight', this.value)"></td>
            <td><button class="btn btn-sm btn-outline-danger" onclick="removeExercise(${index})">
                <i class="bi bi-trash"></i></button></td>
        `;
        tbody.appendChild(row);
    });
}

function updateExercise(index, field, value) {
    selectedExercises[index][field] = field === 'weight' ? parseFloat(value) : parseInt(value);
}

function removeExercise(index) {
    selectedExercises.splice(index, 1);
    
    selectedExercises.forEach((item, idx) => {
        item.order = idx + 1;
    });
    exerciseOrder = selectedExercises.length;
    renderRoutineExercises();
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

document.getElementById('routineForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    hideError();
    
    
    const exercisesToSend = selectedExercises.map(ex => ({
        exercise_id: ex.exercise_id,
        order: ex.order,
        sets: ex.sets,
        reps: ex.reps,
        weight: ex.weight || 0
    }));
    
    const routineData = {
        name: document.getElementById('name').value,
        description: document.getElementById('description').value,
        is_public: document.getElementById('is_public').checked,
        exercises: exercisesToSend
    };

    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/routines/${routineId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
            body: JSON.stringify(routineData)
        });

        if (response.ok) {
            window.location.href = `/routines/${routineId}`;
        } else {
            const error = await response.json();
            showError(error.error || 'Error al actualizar la rutina');
        }
    } catch (error) {
        showError('Error al actualizar la rutina');
        console.error('Error:', error);
    }
});

document.addEventListener('DOMContentLoaded', loadRoutine);

