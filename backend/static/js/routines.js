// Routines list page functionality

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
    
    if (routines.length === 0) {
        container.innerHTML = `
            <div class="col-12">
                <div class="alert alert-info">
                    <i class="bi bi-info-circle"></i> No tienes rutinas creadas. 
                    <a href="/routines/create">Crea tu primera rutina</a>
                </div>
            </div>
        `;
        return;
    }

    routines.forEach(routine => {
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
                </div>
                <div class="card-footer bg-transparent">
                    <a href="/routines/${routine.id}" class="btn btn-primary btn-sm">
                        <i class="bi bi-eye"></i> Ver Detalles
                    </a>
                    <button class="btn btn-outline-danger btn-sm" onclick="deleteRoutine('${routine.id}')">
                        <i class="bi bi-trash"></i>
                    </button>
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

document.addEventListener('DOMContentLoaded', loadRoutines);

