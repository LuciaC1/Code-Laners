let allUsers = []; 
document.addEventListener('DOMContentLoaded', function() {
    if (!requireAuth()) {
        return;
    }
    checkAdminAccess();
    loadUsers();
    let nameFilterTimeout;
    const nameInput = document.getElementById('filterUserName');
    if (nameInput) {
        nameInput.addEventListener('input', function() {
            clearTimeout(nameFilterTimeout);
            nameFilterTimeout = setTimeout(() => {
                loadUsers(); 
            }, 500); 
        });
    }
    const roleSelect = document.getElementById('filterUserRole');
    if (roleSelect) {
        roleSelect.addEventListener('change', function() {
            applyFilters(); 
        });
    }
});
async function checkAdminAccess() {
    if (!requireAuth()) {
        return;
    }
}
async function loadUsers() {
    try {
        const headers = getApiHeaders();
        const nameFilter = document.getElementById('filterUserName')?.value || '';
        const url = nameFilter 
            ? `/api/admin/users?name=${encodeURIComponent(nameFilter)}`
            : '/api/admin/users';
        const response = await fetch(url, {
            headers: headers
        });
        if (!response.ok) {
            throw new Error('Error al cargar usuarios');
        }
        const data = await response.json();
        allUsers = data.users || [];
        applyFilters(); 
    } catch (error) {
        console.error('Error loading users:', error);
        const tbody = document.getElementById('usersTableBody');
        if (tbody) {
            tbody.innerHTML = '<tr><td colspan="7" class="text-center text-danger">Error al cargar usuarios</td></tr>';
        }
    }
}
function renderUsers(users) {
    const tbody = document.getElementById('usersTableBody');
    if (!tbody) return;
    if (users.length === 0) {
        tbody.innerHTML = '<tr><td colspan="7" class="text-center text-muted">No hay usuarios registrados</td></tr>';
        return;
    }
    tbody.innerHTML = users.map(user => {
        const date = new Date(user.created_at);
        const formattedDate = date.toLocaleDateString('es-ES', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        });
        const roleBadge = user.role === 'admin' 
            ? '<span class="badge bg-danger">Administrador</span>'
            : '<span class="badge bg-primary">Usuario</span>';
        let userId = 'N/A';
        if (user.id) {
            userId = typeof user.id === 'string' ? user.id : (user.id.$oid || user.id.toString());
        } else if (user._id) {
            userId = typeof user._id === 'string' ? user._id : (user._id.$oid || user._id.toString());
        }
        const shortId = userId !== 'N/A' && userId.length > 10 ? userId.substring(0, 10) + '...' : userId;
        return `
            <tr>
                <td><code>${shortId}</code></td>
                <td>${user.name || 'N/A'}</td>
                <td>${user.email || 'N/A'}</td>
                <td>${roleBadge}</td>
                <td>${formattedDate}</td>
                <td>
                    <span class="badge bg-success">Activo</span>
                </td>
                <td>
                    <button class="btn btn-sm btn-outline-primary" onclick="showUserDetails('${userId}')" title="Ver detalles">
                        <i class="bi bi-eye"></i>
                    </button>
                    <button class="btn btn-sm btn-outline-warning" onclick="toggleUserRole('${userId}', '${user.role}')" title="Cambiar rol">
                        <i class="bi bi-person-gear"></i>
                    </button>
                </td>
            </tr>
        `;
    }).join('');
}
async function showUserDetails(userId) {
    try {
        const headers = getApiHeaders();
        const response = await fetch(`/api/admin/users/${userId}`, {
            headers: headers
        });
        if (!response.ok) {
            alert('Error al cargar detalles del usuario');
            return;
        }
        const user = await response.json();
        const modalContent = document.getElementById('userDetailsContent');
        const dateOfBirth = user.date_of_birth 
            ? new Date(user.date_of_birth).toLocaleDateString('es-ES')
            : 'No especificada';
        const createdDate = new Date(user.created_at).toLocaleDateString('es-ES', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        });
        modalContent.innerHTML = `
            <div class="row">
                <div class="col-md-6">
                    <h6>Información Personal</h6>
                    <p><strong>Nombre:</strong> ${user.name || 'N/A'}</p>
                    <p><strong>Email:</strong> ${user.email || 'N/A'}</p>
                    <p><strong>Fecha de Nacimiento:</strong> ${dateOfBirth}</p>
                    <p><strong>Rol:</strong> <span class="badge ${user.role === 'admin' ? 'bg-danger' : 'bg-primary'}">${user.role || 'user'}</span></p>
                </div>
                <div class="col-md-6">
                    <h6>Datos Físicos</h6>
                    <p><strong>Peso:</strong> ${user.weight || 'No especificado'} kg</p>
                    <p><strong>Altura:</strong> ${user.height || 'No especificada'} cm</p>
                    <p><strong>Nivel:</strong> ${user.level || 'No especificado'}</p>
                    <p><strong>Objetivos:</strong> ${user.goals && user.goals.length > 0 ? user.goals.join(', ') : 'No especificados'}</p>
                </div>
            </div>
            <hr>
            <div class="row">
                <div class="col-12">
                    <h6>Información del Sistema</h6>
                    <p><strong>ID:</strong> <code>${user.id || user._id || 'N/A'}</code></p>
                    <p><strong>Fecha de Registro:</strong> ${createdDate}</p>
                </div>
            </div>
        `;
        const modal = new bootstrap.Modal(document.getElementById('userDetailsModal'));
        modal.show();
    } catch (error) {
        console.error('Error loading user details:', error);
        alert('Error al cargar detalles del usuario');
    }
}
async function toggleUserRole(userId, currentRole) {
    const newRole = currentRole === 'admin' ? 'user' : 'admin';
    const confirmMessage = `¿Estás seguro de cambiar el rol de este usuario a ${newRole === 'admin' ? 'Administrador' : 'Usuario'}?`;
    if (!confirm(confirmMessage)) {
        return;
    }
    try {
        const headers = getApiHeaders();
        const response = await fetch(`/api/admin/users/${userId}/role`, {
            method: 'PUT',
            headers: {
                ...headers,
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ role: newRole })
        });
        if (response.ok) {
            alert('Rol actualizado correctamente');
            loadUsers();
        } else {
            const error = await response.json();
            alert(error.error || 'Error al actualizar el rol');
        }
    } catch (error) {
        console.error('Error updating user role:', error);
        alert('Error al actualizar el rol');
    }
}
function applyFilters() {
    const role = document.getElementById('filterUserRole')?.value || '';
    let filteredUsers = [...allUsers];
    if (role) {
        filteredUsers = filteredUsers.filter(user => user.role === role);
    }
    renderUsers(filteredUsers);
}
