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
            const tbody = document.getElementById('logsTableBody');
            if (tbody) {
                tbody.innerHTML = '<tr><td colspan="6" class="text-center text-danger">Error al cargar los logs</td></tr>';
            }
            return;
        }
        
        const data = await response.json();
        const logs = data.logs || [];
        renderLogs(logs);
    } catch (error) {
        console.error('Error loading logs:', error);
        const tbody = document.getElementById('logsTableBody');
        if (tbody) {
            tbody.innerHTML = '<tr><td colspan="6" class="text-center text-danger">Error al cargar los logs</td></tr>';
        }
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
        'auth': '<span class="badge bg-danger">Autenticación</span>',
        'general': '<span class="badge bg-secondary">General</span>'
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
            const errorData = await response.json().catch(() => ({}));
            alert('Error al limpiar logs: ' + (errorData.error || 'Error desconocido'));
        }
    } catch (error) {
        console.error('Error clearing logs:', error);
        alert('Error al limpiar logs');
    }
}
