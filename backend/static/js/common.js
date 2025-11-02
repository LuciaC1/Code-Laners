// Common functionality for the application

// Check authentication status and update header
function updateAuthSection() {
    const authSection = document.getElementById('authSection');
    if (!authSection) return; // Exit if authSection doesn't exist
    
    const token = localStorage.getItem('token');
    
    if (token) {
        // User is logged in - get user info from token if needed
        // For now, just show a generic user menu
        // Get user info from localStorage if available
        const userStr = localStorage.getItem('user');
        let userName = 'Usuario';
        let user = null;
        if (userStr) {
            try {
                user = JSON.parse(userStr);
                userName = user.name || 'Usuario';
            } catch (e) {
                // Ignore parse errors
            }
        }

        // Check if user is admin
        const isAdmin = user && user.role === 'admin';
        const adminMenu = isAdmin ? `
            <li><hr class="dropdown-divider"></li>
            <li><a class="dropdown-item" href="/admin">
                <i class="bi bi-speedometer2 me-2"></i>Panel Admin
            </a></li>
        ` : '';

        authSection.innerHTML = `
            <li class="nav-item dropdown">
                <a class="nav-link dropdown-toggle" href="#" role="button" data-bs-toggle="dropdown">
                    <i class="bi bi-person-circle me-1"></i>
                    ${userName}
                </a>
                <ul class="dropdown-menu dropdown-menu-end">
                    <li><a class="dropdown-item" href="/profile">
                        <i class="bi bi-person me-2"></i>Perfil
                    </a></li>
                    <li><a class="dropdown-item" href="/stats">
                        <i class="bi bi-graph-up me-2"></i>Estadísticas
                    </a></li>
                    ${adminMenu}
                    <li><hr class="dropdown-divider"></li>
                    <li><a class="dropdown-item text-danger" href="#" onclick="logout()">
                        <i class="bi bi-box-arrow-right me-2"></i>Cerrar Sesión
                    </a></li>
                </ul>
            </li>
        `;
        
        // Also update server-side menu if it exists
        const serverMenu = document.getElementById('userMenu');
        if (serverMenu && isAdmin) {
            const adminLinks = serverMenu.querySelectorAll('#adminMenuLink');
            adminLinks.forEach(link => link.style.display = '');
        }
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

// Create user dropdown menu dynamically
function createUserDropdown(user) {
    const userStr = localStorage.getItem('user');
    let userName = 'Usuario';
    let isAdmin = false;
    
    if (user) {
        userName = user.name || 'Usuario';
        isAdmin = user.role === 'admin';
    } else if (userStr) {
        try {
            const parsedUser = JSON.parse(userStr);
            userName = parsedUser.name || 'Usuario';
            isAdmin = parsedUser.role === 'admin';
        } catch (e) {
            // Ignore parse errors
        }
    }
    
    const adminMenu = isAdmin ? `
        <li><hr class="dropdown-divider"></li>
        <li><a class="dropdown-item" href="/admin"><i class="bi bi-speedometer2"></i> Panel Admin</a></li>
    ` : '';
    
    return `
        <li class="nav-item dropdown">
            <a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-bs-toggle="dropdown">
                <i class="bi bi-person-circle"></i> ${userName}
            </a>
            <ul class="dropdown-menu dropdown-menu-end" id="userMenu">
                <li><a class="dropdown-item" href="/profile"><i class="bi bi-person"></i> Perfil</a></li>
                <li><a class="dropdown-item" href="/stats"><i class="bi bi-graph-up"></i> Estadísticas</a></li>
                ${adminMenu}
                <li><hr class="dropdown-divider"></li>
                <li><a class="dropdown-item" href="#" onclick="logout(); return false;"><i class="bi bi-box-arrow-right"></i> Cerrar Sesión</a></li>
            </ul>
        </li>
    `;
}

// Update navbar after login (hides login/register links, shows user menu)
function updateNavbarAfterLogin(user) {
    const token = localStorage.getItem('token');
    
    if (token) {
        // Hide login and register links
        const loginLinks = document.querySelectorAll('a[href="/login"]');
        const registerLinks = document.querySelectorAll('a[href="/register"]');
        
        loginLinks.forEach(link => {
            const navItem = link.closest('li.nav-item');
            if (navItem && !navItem.classList.contains('dropdown')) {
                navItem.style.display = 'none';
            }
        });
        
        registerLinks.forEach(link => {
            const navItem = link.closest('li.nav-item');
            if (navItem && !navItem.classList.contains('dropdown')) {
                navItem.style.display = 'none';
            }
        });
        
        // Check if user dropdown exists
        let userDropdown = document.querySelector('li.nav-item.dropdown');
        
        if (!userDropdown) {
            // Create user dropdown if it doesn't exist
            const navbarNav = document.querySelector('.navbar-nav:last-child');
            if (navbarNav) {
                // Remove login/register links first
                loginLinks.forEach(link => {
                    const navItem = link.closest('li.nav-item');
                    if (navItem) {
                        navItem.remove();
                    }
                });
                registerLinks.forEach(link => {
                    const navItem = link.closest('li.nav-item');
                    if (navItem) {
                        navItem.remove();
                    }
                });
                // Create and add user dropdown
                const tempDiv = document.createElement('div');
                tempDiv.innerHTML = createUserDropdown(user);
                navbarNav.appendChild(tempDiv.firstElementChild);
                userDropdown = document.querySelector('li.nav-item.dropdown');
            }
        } else {
            // Update existing dropdown
            userDropdown.style.display = '';
            
            // Update user name in dropdown if available
            if (user && user.name) {
                const dropdownToggle = userDropdown.querySelector('.dropdown-toggle');
                if (dropdownToggle) {
                    dropdownToggle.innerHTML = `<i class="bi bi-person-circle"></i> ${user.name}`;
                }
            }
        }
    }
}

// Update navbar on page load based on token
function updateNavbarOnLoad() {
    const token = localStorage.getItem('token');
    
    if (token) {
        // User is logged in, hide login/register, show user menu
        // Check if user dropdown exists first
        let userDropdown = document.querySelector('li.nav-item.dropdown');
        
        if (!userDropdown) {
            // Create user dropdown if it doesn't exist (e.g., on index page)
            const navbarNav = document.querySelector('.navbar-nav:last-child');
            if (navbarNav) {
                // Get user info from localStorage
                const userStr = localStorage.getItem('user');
                let user = null;
                if (userStr) {
                    try {
                        user = JSON.parse(userStr);
                    } catch (e) {
                        // Ignore parse errors
                    }
                }
                // Remove login/register links first
                const loginLinks = document.querySelectorAll('a[href="/login"]');
                const registerLinks = document.querySelectorAll('a[href="/register"]');
                
                loginLinks.forEach(link => {
                    const navItem = link.closest('li.nav-item');
                    if (navItem && !navItem.classList.contains('dropdown')) {
                        navItem.remove();
                    }
                });
                registerLinks.forEach(link => {
                    const navItem = link.closest('li.nav-item');
                    if (navItem && !navItem.classList.contains('dropdown')) {
                        navItem.remove();
                    }
                });
                // Create and add user dropdown
                const tempDiv = document.createElement('div');
                tempDiv.innerHTML = createUserDropdown(user);
                navbarNav.appendChild(tempDiv.firstElementChild);
            }
        } else {
            // Show existing dropdown and hide login/register links
            userDropdown.style.display = '';
            
            const loginLinks = document.querySelectorAll('a[href="/login"]');
            const registerLinks = document.querySelectorAll('a[href="/register"]');
            
            loginLinks.forEach(link => {
                const navItem = link.closest('li.nav-item');
                if (navItem && !navItem.classList.contains('dropdown')) {
                    navItem.style.display = 'none';
                }
            });
            
            registerLinks.forEach(link => {
                const navItem = link.closest('li.nav-item');
                if (navItem && !navItem.classList.contains('dropdown')) {
                    navItem.style.display = 'none';
                }
            });
        }
    } else {
        // User is not logged in, show login/register, hide user menu
        const loginLinks = document.querySelectorAll('a[href="/login"]');
        const registerLinks = document.querySelectorAll('a[href="/register"]');
        const userDropdown = document.querySelector('li.nav-item.dropdown');
        
        loginLinks.forEach(link => {
            const navItem = link.closest('li.nav-item');
            if (navItem && !navItem.classList.contains('dropdown')) {
                navItem.style.display = '';
            }
        });
        
        registerLinks.forEach(link => {
            const navItem = link.closest('li.nav-item');
            if (navItem && !navItem.classList.contains('dropdown')) {
                navItem.style.display = '';
            }
        });
        
        if (userDropdown) {
            userDropdown.style.display = 'none';
        }
    }
}

// Initialize common functionality on all pages
document.addEventListener('DOMContentLoaded', function() {
    // Check if we're on the index page
    const isIndexPage = window.location.pathname === '/';
    
    // Always update navbar on load
    updateNavbarOnLoad();
    
    if (!isIndexPage) {
        updateAuthSection();
        
        // Update auth section when storage changes
        window.addEventListener('storage', function(e) {
            if (e.key === 'token') {
                updateAuthSection();
                updateNavbarOnLoad();
            }
        });
    } else {
        // On index page, listen for storage changes to update navbar
        window.addEventListener('storage', function(e) {
            if (e.key === 'token') {
                updateNavbarOnLoad();
            }
        });
    }
});