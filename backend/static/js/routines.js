async function loadRoutines() {
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/routines', {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });
        const data = await response.json();
        renderRoutines(data.routines || []);
    } catch (error) {
        console.error('Error loading routines:', error);
        document.getElementById('routinesContainer').innerHTML = 
            '<div class="col-12"><div class="alert alert-danger">Error al cargar rutinas</div></div>';
    }
}
function renderRoutines(routines) {
    const container = document.getElementById('routinesContainer');
    container.innerHTML = '';
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
    if (routines.length === 0) {
        container.innerHTML = `
            <div class="col-12">
                <div class="alert alert-info">
                    <i class="bi bi-info-circle"></i> No hay rutinas disponibles. 
                    <a href="/routines/create">Crea tu primera rutina</a>
                </div>
            </div>
        `;
        return;
    }
    routines.forEach(routine => {
        const isOwner = currentUserId && routine.user_id === currentUserId;
        const card = document.createElement('div');
        card.className = 'col-md-6 col-lg-4';
        card.innerHTML = `
            <div class="card h-100 shadow-sm">
                <div class="card-body">
                    <h5 class="card-title">${routine.name}</h5>
                    <p class="card-text text-muted">${routine.description || 'Sin descripción'}</p>
                    <div class="mb-2">
                        <span class="badge bg-primary">${routine.exercises?.length || 0} ejercicios</span>
                        ${routine.is_public ? '<span class="badge bg-success">Pública</span>' : '<span class="badge bg-secondary">Privada</span>'}
                    </div>
                    ${routine.owner_name && !isOwner ? `<small class="text-muted"><i class="bi bi-person"></i> Creada por: ${routine.owner_name}</small>` : ''}
                </div>
                <div class="card-footer bg-transparent">
                    <a href="/routines/${routine.id}" class="btn btn-primary btn-sm">
                        <i class="bi bi-eye"></i> Ver
                    </a>
                    ${isOwner ? `
                        <button class="btn btn-outline-success btn-sm" onclick="duplicateRoutine('${routine.id}')" title="Duplicar rutina">
                            <i class="bi bi-files"></i>
                        </button>
                        <button class="btn btn-outline-danger btn-sm" onclick="deleteRoutine('${routine.id}')" title="Eliminar rutina">
                            <i class="bi bi-trash"></i>
                        </button>
                    ` : `
                        <button class="btn btn-outline-success btn-sm" onclick="duplicateRoutine('${routine.id}')" title="Duplicar rutina">
                            <i class="bi bi-files"></i> Duplicar
                        </button>
                    `}
                </div>
            </div>
        `;
        container.appendChild(card);
    });
}
async function deleteRoutine(id) {
    if (!confirm('¿Estás seguro de eliminar esta rutina?')) return;
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/routines/${id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });
        if (response.ok) {
            loadRoutines();
        } else {
            alert('Error al eliminar la rutina');
        }
    } catch (error) {
        alert('Error al eliminar la rutina');
    }
}
async function duplicateRoutine(id) {
    if (!confirm('¿Deseas duplicar esta rutina?')) return;
    try {
        const token = localStorage.getItem('token');
        const getResponse = await fetch(`/api/routines/${id}`, {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });
        if (!getResponse.ok) {
            alert('Error al obtener la rutina');
            return;
        }
        const routine = await getResponse.json();
        const newRoutine = {
            name: `${routine.name} (Copia)`,
            description: routine.description || '',
            exercises: routine.exercises || [],
            is_public: false
        };
        const createResponse = await fetch('/api/routines', {
            method: 'POST',
            headers: {
                'Authorization': 'Bearer ' + token,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(newRoutine)
        });
        if (createResponse.ok) {
            const newRoutineData = await createResponse.json();
            alert('Rutina duplicada correctamente');
            window.location.href = `/routines/${newRoutineData.id || newRoutineData.routine?.id}`;
        } else {
            const error = await createResponse.json();
            alert(error.error || 'Error al duplicar la rutina');
        }
    } catch (error) {
        alert('Error al duplicar la rutina');
        console.error('Error:', error);
    }
}
document.addEventListener('DOMContentLoaded', loadRoutines);
