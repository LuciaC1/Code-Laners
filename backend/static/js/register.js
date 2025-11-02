
document.addEventListener('DOMContentLoaded', function() {
    const registerForm = document.getElementById('registerForm');
    const errorMessage = document.getElementById('errorMessage');
    const successMessage = document.getElementById('successMessage');

    if (registerForm) {
        registerForm.addEventListener('submit', async function(e) {
            e.preventDefault();
            
            
            errorMessage.classList.add('d-none');
            successMessage.classList.add('d-none');

            
            const formData = new FormData(registerForm);
            
            
            const goals = [];
            const goalCheckboxes = registerForm.querySelectorAll('input[name="goals"]:checked');
            goalCheckboxes.forEach(checkbox => {
                goals.push(checkbox.value);
            });

            
            const requestData = {
                name: formData.get('name'),
                email: formData.get('email'),
                password: formData.get('password'),
                date_of_birth: formData.get('date_of_birth')
            };

            
            const weight = formData.get('weight');
            const height = formData.get('height');
            const level = formData.get('level');

            if (weight && weight.trim() !== '') {
                requestData.weight = parseFloat(weight);
            }
            if (height && height.trim() !== '') {
                requestData.height = parseFloat(height);
            }
            if (level && level.trim() !== '') {
                requestData.level = level;
            }
            if (goals.length > 0) {
                requestData.goals = goals;
            }

            try {
                const response = await fetch('/api/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(requestData)
                });

                const data = await response.json();

                if (response.ok) {
                    successMessage.textContent = 'Registro exitoso. Redirigiendo al login...';
                    successMessage.classList.remove('d-none');
                    registerForm.reset();
                    
                    
                    setTimeout(() => {
                        window.location.href = '/login';
                    }, 2000);
                } else {
                    errorMessage.textContent = data.error || 'Error al registrar usuario';
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

