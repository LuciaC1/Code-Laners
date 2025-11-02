// Routine detail page functionality

let currentRoutineId = null;

async function loadRoutine() {
    const routineId = window.location.pathname.split('/').pop();
    currentRoutineId = routineId;
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

        const routine = await response.json();
        renderRoutine(routine);
    } catch (error) {
        document.getElementById('loadingSpinner').classList.add('d-none');
        document.getElementById('errorMessage').textContent = error.message || 'Error al cargar la rutina';
        document.getElementById('errorMessage').classList.remove('d-none');
    }
}

function renderRoutine(routine) {
    document.getElementById('loadingSpinner').classList.add('d-none');
    document.getElementById('routineContent').classList.remove('d-none');

    document.getElementById('routineName').textContent = routine.name || 'Sin nombre';
    document.getElementById('routineDescription').textContent = routine.description || 'Sin descripción';

    // Get current user ID
    const userStr = localStorage.getItem('user');
    let currentUserId = null;
    if (userStr) {
        try {
            const user = JSON.parse(userStr);
            currentUserId = user.id || user._id;
        } catch (e) {
            console.error('Error parsing user:', e);
        }
    }
    
    const isOwner = currentUserId && routine.user_id === currentUserId;

    const badges = document.getElementById('routineBadges');
    let badgeHTML = routine.is_public 
        ? '<span class="badge bg-success fs-6">Rutina Pública</span>'
        : '<span class="badge bg-secondary fs-6">Rutina Privada</span>';
    
    // Show owner name if it's a public routine and user is not the owner
    if (routine.owner_name && !isOwner) {
        badgeHTML += ` <small class="text-muted ms-2"><i class="bi bi-person"></i> Creada por: ${routine.owner_name}</small>`;
    }
    
    badges.innerHTML = badgeHTML;

    const tbody = document.getElementById('exercisesTableBody');
    tbody.innerHTML = '';
    
    if (routine.exercises && routine.exercises.length > 0) {
        routine.exercises.forEach(exercise => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${exercise.order || '-'}</td>
                <td>${exercise.exercise_name || 'Ejercicio sin nombre'}</td>
                <td>${exercise.sets || '-'}</td>
                <td>${exercise.reps || '-'}</td>
                <td>${exercise.weight ? exercise.weight : '-'}</td>
            `;
            tbody.appendChild(row);
        });
    } else {
        tbody.innerHTML = '<tr><td colspan="5" class="text-center text-muted">No hay ejercicios en esta rutina</td></tr>';
    }

    const totalSets = routine.exercises ? routine.exercises.reduce((sum, ex) => sum + (ex.sets || 0), 0) : 0;
    const estimatedTime = Math.ceil(totalSets * 2); // Estimate 2 minutes per set

    const summary = document.getElementById('routineSummary');
    summary.innerHTML = `
        <li class="mb-2"><strong>Total de Ejercicios:</strong> ${routine.exercises?.length || 0}</li>
        <li class="mb-2"><strong>Total de Series:</strong> ${totalSets}</li>
        <li class="mb-2"><strong>Tiempo Estimado:</strong> ${estimatedTime} min</li>
    `;
    
    // Show/hide edit and delete buttons based on ownership
    const editBtn = document.querySelector('button[onclick*="editRoutine"]');
    const deleteBtn = document.querySelector('button[onclick*="deleteRoutine"]');
    if (editBtn) {
        editBtn.style.display = isOwner ? '' : 'none';
    }
    if (deleteBtn) {
        deleteBtn.style.display = isOwner ? '' : 'none';
    }
}

function startWorkout() {
    if (currentRoutineId) {
        window.location.href = `/workouts/create?routine_id=${currentRoutineId}`;
    }
}

function editRoutine() {
    if (currentRoutineId) {
        window.location.href = `/routines/${currentRoutineId}/edit`;
    }
}

async function deleteRoutine() {
    if (!confirm('¿Estás seguro de eliminar esta rutina? Esta acción no se puede deshacer.')) {
        return;
    }
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/routines/${currentRoutineId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });

        if (response.ok) {
            alert('Rutina eliminada correctamente');
            window.location.href = '/routines';
        } else {
            const data = await response.json();
            alert(data.error || 'Error al eliminar la rutina');
        }
    } catch (error) {
        alert('Error al eliminar la rutina');
        console.error('Error:', error);
    }
}

document.addEventListener('DOMContentLoaded', loadRoutine);

