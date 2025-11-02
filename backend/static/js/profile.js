// Profile page handler
let userData = null;

document.addEventListener('DOMContentLoaded', function() {
    // Check authentication
    if (!requireAuth()) {
        return;
    }

    loadUserProfile();
    setupForms();
});

async function loadUserProfile() {
    try {
        const headers = getApiHeaders();
        const response = await fetch('/api/me', {
            headers: headers
        });

        if (response.ok) {
            userData = await response.json();
            populateProfileForms(userData);
        } else if (response.status === 401) {
            window.location.href = '/login';
        } else {
            showError('personalInfoError', 'Error al cargar el perfil');
        }
    } catch (error) {
        console.error('Error loading profile:', error);
        showError('personalInfoError', 'Error de conexión');
    }
}

function populateProfileForms(user) {
    // Personal Information
    if (document.getElementById('profileName')) {
        document.getElementById('profileName').value = user.name || '';
        document.getElementById('profileEmail').value = user.email || '';
        
        // Format date for input[type="date"]
        if (user.date_of_birth) {
            const date = new Date(user.date_of_birth);
            const formattedDate = date.toISOString().split('T')[0];
            document.getElementById('profileDateOfBirth').value = formattedDate;
        }
    }

    // Physical Data
    if (document.getElementById('profileWeight')) {
        document.getElementById('profileWeight').value = user.weight || '';
        document.getElementById('profileHeight').value = user.height || '';
        document.getElementById('profileLevel').value = user.level || '';
    }

    // Goals
    if (user.goals && Array.isArray(user.goals)) {
        user.goals.forEach(goal => {
            const checkbox = document.querySelector(`input[name="goals"][value="${goal}"]`);
            if (checkbox) {
                checkbox.checked = true;
            }
        });
    }
}

function setupForms() {
    // Personal Information Form
    const personalInfoForm = document.getElementById('personalInfoForm');
    if (personalInfoForm) {
        personalInfoForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            await updatePersonalInfo();
        });
    }

    // Physical Data Form
    const physicalDataForm = document.getElementById('physicalDataForm');
    if (physicalDataForm) {
        physicalDataForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            await updatePhysicalData();
        });
    }

    // Goals Form
    const goalsForm = document.getElementById('goalsForm');
    if (goalsForm) {
        goalsForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            await updateGoals();
        });
    }

    // Change Password Form
    const changePasswordForm = document.getElementById('changePasswordForm');
    if (changePasswordForm) {
        changePasswordForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            await changePassword();
        });
    }
}

async function updatePersonalInfo() {
    const name = document.getElementById('profileName').value;
    const email = document.getElementById('profileEmail').value;

    const requestData = {
        name: name,
        email: email
    };

    try {
        const response = await fetch('/api/me', {
            method: 'PUT',
            headers: getApiHeaders(),
            body: JSON.stringify(requestData)
        });

        const data = await response.json();

        if (response.ok) {
            showSuccess('personalInfoSuccess', 'Información personal actualizada correctamente');
            hideError('personalInfoError');
            // Reload profile to get updated data
            setTimeout(() => loadUserProfile(), 1000);
        } else {
            showError('personalInfoError', data.error || 'Error al actualizar la información');
            hideSuccess('personalInfoSuccess');
        }
    } catch (error) {
        console.error('Error:', error);
        showError('personalInfoError', 'Error de conexión');
        hideSuccess('personalInfoSuccess');
    }
}

async function updatePhysicalData() {
    const weight = document.getElementById('profileWeight').value;
    const height = document.getElementById('profileHeight').value;
    const level = document.getElementById('profileLevel').value;

    const requestData = {};

    if (weight && weight.trim() !== '') {
        requestData.weight = parseFloat(weight);
    }
    if (height && height.trim() !== '') {
        requestData.height = parseFloat(height);
    }
    if (level && level.trim() !== '') {
        requestData.level = level;
    }

    try {
        const response = await fetch('/api/me', {
            method: 'PUT',
            headers: getApiHeaders(),
            body: JSON.stringify(requestData)
        });

        const data = await response.json();

        if (response.ok) {
            showSuccess('physicalDataSuccess', 'Datos físicos actualizados correctamente');
            hideError('physicalDataError');
            setTimeout(() => loadUserProfile(), 1000);
        } else {
            showError('physicalDataError', data.error || 'Error al actualizar los datos físicos');
            hideSuccess('physicalDataSuccess');
        }
    } catch (error) {
        console.error('Error:', error);
        showError('physicalDataError', 'Error de conexión');
        hideSuccess('physicalDataSuccess');
    }
}

async function updateGoals() {
    const goalCheckboxes = document.querySelectorAll('#goalsForm input[name="goals"]:checked');
    const goals = Array.from(goalCheckboxes).map(cb => cb.value);

    const requestData = {
        goals: goals
    };

    try {
        const response = await fetch('/api/me', {
            method: 'PUT',
            headers: getApiHeaders(),
            body: JSON.stringify(requestData)
        });

        const data = await response.json();

        if (response.ok) {
            showSuccess('goalsSuccess', 'Objetivos actualizados correctamente');
            hideError('goalsError');
            setTimeout(() => loadUserProfile(), 1000);
        } else {
            showError('goalsError', data.error || 'Error al actualizar los objetivos');
            hideSuccess('goalsSuccess');
        }
    } catch (error) {
        console.error('Error:', error);
        showError('goalsError', 'Error de conexión');
        hideSuccess('goalsSuccess');
    }
}

async function changePassword() {
    const oldPassword = document.getElementById('oldPassword').value;
    const newPassword = document.getElementById('newPassword').value;
    const confirmPassword = document.getElementById('confirmPassword').value;

    // Validate passwords match
    if (newPassword !== confirmPassword) {
        showError('passwordError', 'Las contraseñas no coinciden');
        hideSuccess('passwordSuccess');
        return;
    }

    const requestData = {
        old_password: oldPassword,
        new_password: newPassword
    };

    try {
        const response = await fetch('/api/me/password', {
            method: 'PUT',
            headers: getApiHeaders(),
            body: JSON.stringify(requestData)
        });

        const data = await response.json();

        if (response.ok) {
            showSuccess('passwordSuccess', 'Contraseña actualizada correctamente. Serás redirigido al login...');
            hideError('passwordError');
            
            // Clear form
            document.getElementById('changePasswordForm').reset();
            
            // Redirect to login after 2 seconds (user needs to re-login)
            setTimeout(() => {
                localStorage.removeItem('token');
                window.location.href = '/login';
            }, 2000);
        } else {
            showError('passwordError', data.error || 'Error al cambiar la contraseña');
            hideSuccess('passwordSuccess');
        }
    } catch (error) {
        console.error('Error:', error);
        showError('passwordError', 'Error de conexión');
        hideSuccess('passwordSuccess');
    }
}

function showError(elementId, message) {
    const errorElement = document.getElementById(elementId);
    if (errorElement) {
        errorElement.textContent = message;
        errorElement.classList.remove('d-none');
    }
}

function hideError(elementId) {
    const errorElement = document.getElementById(elementId);
    if (errorElement) {
        errorElement.classList.add('d-none');
    }
}

function showSuccess(elementId, message) {
    const successElement = document.getElementById(elementId);
    if (successElement) {
        successElement.textContent = message;
        successElement.classList.remove('d-none');
    }
}

function hideSuccess(elementId) {
    const successElement = document.getElementById(elementId);
    if (successElement) {
        successElement.classList.add('d-none');
    }
}

