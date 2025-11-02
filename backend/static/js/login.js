// Login form handler
document.addEventListener('DOMContentLoaded', function() {
    const loginForm = document.getElementById('loginForm');
    const errorMessage = document.getElementById('errorMessage');

    if (loginForm) {
        loginForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            // Hide previous messages
            errorMessage.classList.add('d-none');

            // Get form data
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
                    // Save token to localStorage
                    localStorage.setItem('token', data.accessToken);
                    
                    // Save user info if available
                    if (data.user) {
                        localStorage.setItem('user', JSON.stringify(data.user));
                    }
                    
                    // Redirect to home
                    window.location.href = '/';
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

