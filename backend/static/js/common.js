// Common functionality for the application

// Check authentication status and update header
function updateAuthSection() {
    const authSection = document.getElementById('authSection');
    if (!authSection) return; // Exit if authSection doesn't exist
    
    const token = localStorage.getItem('token');
    
    if (token) {
        // User is logged in - get user info from token if needed
        // For now, just show a generic user menu
        authSection.innerHTML = `
            <li class="nav-item dropdown">
                <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
                    <i class="bi bi-person-circle me-1"></i>
                    Usuario
                </a>
                <ul class="dropdown-menu dropdown-menu-end">
                    <li><a class="dropdown-item" href="/products">
                        <i class="bi bi-list-ul me-2"></i>Mis Productos
                    </a></li>
                    <li><hr class="dropdown-divider"></li>
                    <li><a class="dropdown-item text-danger" href="#" onclick="logout()">
                        <i class="bi bi-box-arrow-right me-2"></i>Cerrar Sesión
                    </a></li>
                </ul>
            </li>
        `;
    } else {
        // User is not logged in
        authSection.innerHTML = `
            <li class="nav-item">
                <a class="nav-link" href="/login">
                    <i class="bi bi-box-arrow-in-right me-1"></i>Iniciar Sesión
                </a>
            </li>
            <li class="nav-item">
                <a class="nav-link" href="/register">
                    <i class="bi bi-person-plus me-1"></i>Registrarse
                </a>
            </li>
        `;
    }
}

// Logout function
function logout() {
    if (confirm('¿Estás seguro de que quieres cerrar sesión?')) {
        localStorage.removeItem('token');
        window.location.href = '/login';
    }
}

// Check if user is authenticated
function isAuthenticated() {
    const token = localStorage.getItem('token');
    return !!token;
}

// Redirect to login if not authenticated
function requireAuth() {
    if (!isAuthenticated()) {
        window.location.href = '/login';
        return false;
    }
    return true;
}

// Get API headers with authentication
function getApiHeaders() {
    const token = localStorage.getItem('token');
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    };
}

// Initialize common functionality - only if not on index page
document.addEventListener('DOMContentLoaded', function() {
    // Check if we're on the index page
    const isIndexPage = window.location.pathname === '/';
    
    if (!isIndexPage) {
        updateAuthSection();
        
        // Update auth section when storage changes
        window.addEventListener('storage', function(e) {
            if (e.key === 'token') {
                updateAuthSection();
            }
        });
    }
});