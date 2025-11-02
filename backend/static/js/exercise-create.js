function addStep() {
    const container = document.getElementById('stepsContainer');
    const stepCount = container.children.length + 1;
    const stepDiv = document.createElement('div');
    stepDiv.className = 'input-group mb-2';
    stepDiv.innerHTML = `
        <input type="text" class="form-control step-input" placeholder="Paso ${stepCount}">
        <button type="button" class="btn btn-outline-danger" onclick="removeStep(this)">
            <i class="bi bi-trash"></i>
        </button>
    `;
    container.appendChild(stepDiv);
}
function removeStep(button) {
    button.parentElement.remove();
}
function showError(message) {
    const errorDiv = document.getElementById('errorMessage');
    errorDiv.textContent = message;
    errorDiv.classList.remove('d-none');
}
document.getElementById('exerciseForm').addEventListener('submit', async (e) => {
    e.preventDefault();
    const steps = Array.from(document.querySelectorAll('.step-input'))
        .map(input => input.value)
        .filter(value => value.trim() !== '');
    const exerciseData = {
        name: document.getElementById('name').value,
        description: document.getElementById('description').value,
        category: document.getElementById('category').value,
        muscle_group: document.getElementById('muscle_group').value,
        difficulty: document.getElementById('difficulty').value,
        media_url: document.getElementById('media_url').value,
        steps: steps
    };
    try {
        const token = localStorage.getItem('token');
        const response = await fetch('/api/exercises', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
            body: JSON.stringify(exerciseData)
        });
        if (response.ok) {
            window.location.href = '/exercises';
        } else {
            const error = await response.json();
            showError(error.error || 'Error al crear el ejercicio');
        }
    } catch (error) {
        showError('Error al crear el ejercicio');
    }
});
