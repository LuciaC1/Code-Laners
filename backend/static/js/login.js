document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    const errorMessage = document.getElementById('errorMessage');
    if (loginForm) {
        loginForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            e.stopPropagation();
            errorMessage.classList.add('d-none');
            const formData = new FormData(loginForm);
            const requestData = {
                email: formData.get('email'),
                password: formData.get('password')
            };
            try {
                const response = await fetch('/api/login', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(requestData)
                });
                const data = await response.json();
                if (response.ok && data.accessToken) {
                    localStorage.setItem('token', data.accessToken);
                    if (data.user) {
                        localStorage.setItem('user', JSON.stringify(data.user));
                    }
                    if (typeof updateNavbarAfterLogin === 'function') {
                        updateNavbarAfterLogin(data.user);
                    }
                    window.location.href = '/profile?login=success';
                } else {
                    errorMessage.textContent = data.error || 'Email o contraseña incorrectos';
                    errorMessage.classList.remove('d-none');
                }
            } catch (error) {
                errorMessage.textContent = 'Error de conexión. Por favor, intenta nuevamente.';
                errorMessage.classList.remove('d-none');
                console.error('Error:', error);
            }
        });
    }
});
