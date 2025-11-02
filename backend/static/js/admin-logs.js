document.addEventListener('DOMContentLoaded', function() {
    if (!requireAuth()) {
        return;
    }
    checkAdminAccess();
    loadLogs();
});
async function checkAdminAccess() {
    if (!requireAuth()) {
        return;
    }
}
async function loadLogs() {
    try {
        const headers = getApiHeaders();
        const response = await fetch('/api/admin/logs', {
            headers: headers
        });
        if (!response.ok) {
            renderMockLogs();
            return;
        }
        const data = await response.json();
        const logs = data.logs || [];
        renderLogs(logs);
    } catch (error) {
        console.error('Error loading logs:', error);
        renderMockLogs();
    }
}
function renderLogs(logs) {
    const tbody = document.getElementById('logsTableBody');
    if (!tbody) return;
    if (logs.length === 0) {
        tbody.innerHTML = '<tr><td colspan="6" class="text-center text-muted">No hay logs registrados</td></tr>';
        return;
    }
    tbody.innerHTML = logs.map(log => {
        const date = new Date(log.timestamp || log.created_at);
        const formattedDate = date.toLocaleDateString('es-ES', {
            year: 'numeric',
            month: 'short',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
            second: '2-digit'
        });
        const levelBadge = getLevelBadge(log.level || 'info');
        const typeBadge = getTypeBadge(log.type || 'general');
        return `
            <tr>
                <td>${formattedDate}</td>
                <td>${levelBadge}</td>
                <td>${typeBadge}</td>
                <td>${log.action || 'N/A'}</td>
                <td>${log.user_name || log.user_id || 'Sistema'}</td>
                <td>
                    <button class="btn btn-sm btn-outline-info" onclick="showLogDetails(${JSON.stringify(log).replace(/"/g, '&quot;')})" title="Ver detalles">
                        <i class="bi bi-info-circle"></i>
                    </button>
                </td>
            </tr>
        `;
    }).join('');
}
function renderMockLogs() {
    const mockLogs = [
        {
            timestamp: new Date().toISOString(),
            level: 'success',
            type: 'auth',
            action: 'Usuario inició sesión',
            user_name: 'Usuario',
            details: 'Login exitoso'
        },
        {
            timestamp: new Date(Date.now() - 3600000).toISOString(),
            level: 'info',
            type: 'exercise',
            action: 'Ejercicio creado',
            user_name: 'Admin',
            details: 'Nuevo ejercicio agregado al catálogo'
        },
        {
            timestamp: new Date(Date.now() - 7200000).toISOString(),
            level: 'info',
            type: 'routine',
            action: 'Rutina creada',
            user_name: 'Usuario',
            details: 'Nueva rutina de entrenamiento'
        },
        {
            timestamp: new Date(Date.now() - 10800000).toISOString(),
            level: 'success',
            type: 'workout',
            action: 'Entrenamiento completado',
            user_name: 'Usuario',
            details: 'Entrenamiento registrado exitosamente'
        }
    ];
    renderLogs(mockLogs);
}
function getLevelBadge(level) {
    const badges = {
        'info': '<span class="badge bg-info">Info</span>',
        'warning': '<span class="badge bg-warning text-dark">Advertencia</span>',
        'error': '<span class="badge bg-danger">Error</span>',
        'success': '<span class="badge bg-success">Éxito</span>'
    };
    return badges[level] || badges['info'];
}
function getTypeBadge(type) {
    const badges = {
        'user': '<span class="badge bg-primary">Usuario</span>',
        'exercise': '<span class="badge bg-success">Ejercicio</span>',
        'routine': '<span class="badge bg-info">Rutina</span>',
        'workout': '<span class="badge bg-warning text-dark">Entrenamiento</span>',
        'auth': '<span class="badge bg-danger">Autenticación</span>'
    };
    return badges[type] || '<span class="badge bg-secondary">General</span>';
}
function showLogDetails(log) {
    const details = `
        <div class="card">
            <div class="card-body">
                <h6>Detalles del Log</h6>
                <pre class="bg-light p-3 rounded">${JSON.stringify(log, null, 2)}</pre>
            </div>
        </div>
    `;
    alert(`Detalles:\n${JSON.stringify(log, null, 2)}`);
}
function applyLogFilters() {
    loadLogs();
}
async function clearLogs() {
    if (!confirm('¿Estás seguro de limpiar todos los logs? Esta acción no se puede deshacer.')) {
        return;
    }
    try {
        const headers = getApiHeaders();
        const response = await fetch('/api/admin/logs', {
            method: 'DELETE',
            headers: headers
        });
        if (response.ok) {
            alert('Logs limpiados correctamente');
            loadLogs();
        } else {
            alert('Error al limpiar logs');
        }
    } catch (error) {
        console.error('Error clearing logs:', error);
        alert('Error al limpiar logs');
    }
}
