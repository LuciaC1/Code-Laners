function updateAuthSection() {
    const authSection = document.getElementById('authSection');
    if (!authSection) return; 
    const token = localStorage.getItem('token');
    if (token) {
        const userStr = localStorage.getItem('user');
        let userName = 'Usuario';
        let user = null;
        if (userStr) {
            try {
                user = JSON.parse(userStr);
                userName = user.name || 'Usuario';
            } catch (e) {
            }
        }
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
        const serverMenu = document.getElementById('userMenu');
        if (serverMenu && isAdmin) {
            const adminLinks = serverMenu.querySelectorAll('#adminMenuLink');
            adminLinks.forEach(link => link.style.display = '');
        }
    } else {
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
function logout() {
    if (confirm('¿Estás seguro de que quieres cerrar sesión?')) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        window.location.href = '/login';
    }
}
function isAuthenticated() {
    const token = localStorage.getItem('token');
    return !!token;
}
function requireAuth() {
    if (!isAuthenticated()) {
        window.location.href = '/login';
        return false;
    }
    return true;
}
function getApiHeaders() {
    const token = localStorage.getItem('token');
    return {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`
    };
}
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
function updateNavbarAfterLogin(user) {
    const token = localStorage.getItem('token');
    if (token) {
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
        let userDropdown = document.querySelector('li.nav-item.dropdown');
        if (!userDropdown) {
            const navbarNav = document.querySelector('.navbar-nav:last-child');
            if (navbarNav) {
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
                const tempDiv = document.createElement('div');
                tempDiv.innerHTML = createUserDropdown(user);
                navbarNav.appendChild(tempDiv.firstElementChild);
                userDropdown = document.querySelector('li.nav-item.dropdown');
            }
        } else {
            userDropdown.style.display = '';
            if (user && user.name) {
                const dropdownToggle = userDropdown.querySelector('.dropdown-toggle');
                if (dropdownToggle) {
                    dropdownToggle.innerHTML = `<i class="bi bi-person-circle"></i> ${user.name}`;
                }
            }
        }
    }
}
function updateNavbarOnLoad() {
    const token = localStorage.getItem('token');
    if (token) {
        let userDropdown = document.querySelector('li.nav-item.dropdown');
        if (!userDropdown) {
            const navbarNav = document.querySelector('.navbar-nav:last-child');
            if (navbarNav) {
                const userStr = localStorage.getItem('user');
                let user = null;
                if (userStr) {
                    try {
                        user = JSON.parse(userStr);
                    } catch (e) {
                    }
                }
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
                const tempDiv = document.createElement('div');
                tempDiv.innerHTML = createUserDropdown(user);
                navbarNav.appendChild(tempDiv.firstElementChild);
            }
        } else {
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
document.addEventListener('DOMContentLoaded', function() {
    const isIndexPage = window.location.pathname === '/';
    updateNavbarOnLoad();
    if (!isIndexPage) {
        updateAuthSection();
        window.addEventListener('storage', function(e) {
            if (e.key === 'token') {
                updateAuthSection();
                updateNavbarOnLoad();
            }
        });
    } else {
        window.addEventListener('storage', function(e) {
            if (e.key === 'token') {
                updateNavbarOnLoad();
            }
        });
    }
});