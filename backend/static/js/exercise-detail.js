

let currentExerciseId = null;
let currentExercise = null;
let isAdmin = false;

async function loadExercise() {
    const exerciseId = window.location.pathname.split('/').pop();
    currentExerciseId = exerciseId;
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/exercises/${exerciseId}`, {
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });

        if (!response.ok) {
            throw new Error('Ejercicio no encontrado');
        }

        const exercise = await response.json();
        currentExercise = exercise;
        
        
        checkUserRole();
        
        renderExercise(exercise);
    } catch (error) {
        document.getElementById('loadingSpinner').classList.add('d-none');
        document.getElementById('errorMessage').textContent = error.message || 'Error al cargar el ejercicio';
        document.getElementById('errorMessage').classList.remove('d-none');
    }
}

function checkUserRole() {
    
    const userStr = localStorage.getItem('user');
    if (userStr) {
        try {
            const user = JSON.parse(userStr);
            isAdmin = user.role === 'admin';
            if (isAdmin) {
                showAdminControls();
            }
        } catch (e) {
            
        }
    }
}

function showAdminControls() {
    const actionButtons = document.getElementById('exerciseActions');
    if (actionButtons) {
        actionButtons.innerHTML = `
            <button class="btn btn-warning" onclick="editExercise()">
                <i class="bi bi-pencil"></i> Editar Ejercicio
            </button>
            <button class="btn btn-outline-danger" onclick="deleteExercise()">
                <i class="bi bi-trash"></i> Eliminar Ejercicio
            </button>
        `;
    }
}

function renderExercise(exercise) {
    document.getElementById('loadingSpinner').classList.add('d-none');
    document.getElementById('exerciseContent').classList.remove('d-none');

    document.getElementById('exerciseName').textContent = exercise.name || 'Sin nombre';
    document.getElementById('exerciseDescription').textContent = exercise.description || 'Sin descripción';

    const badges = document.getElementById('exerciseBadges');
    badges.innerHTML = `
        <span class="badge bg-primary fs-6">${exercise.category || 'N/A'}</span>
        <span class="badge bg-success fs-6">${exercise.muscle_group || 'N/A'}</span>
        <span class="badge bg-warning text-dark fs-6">${exercise.difficulty || 'N/A'}</span>
    `;

    const info = document.getElementById('exerciseInfo');
    info.innerHTML = `
        <li class="mb-2"><strong>Categoría:</strong> ${exercise.category || 'N/A'}</li>
        <li class="mb-2"><strong>Grupo Muscular:</strong> ${exercise.muscle_group || 'N/A'}</li>
        <li class="mb-2"><strong>Dificultad:</strong> ${exercise.difficulty || 'N/A'}</li>
    `;

    if (exercise.steps && exercise.steps.length > 0) {
        const stepsList = document.getElementById('stepsList');
        stepsList.innerHTML = exercise.steps.map(step => `<li>${step}</li>`).join('');
        document.getElementById('exerciseSteps').classList.remove('d-none');
    }

    if (exercise.media_url) {
        renderMedia(exercise.media_url, exercise.name);
    }
}

function isVideoUrl(url) {
    if (!url) return false;
    
    
    const videoExtensions = ['.mp4', '.webm', '.ogg', '.mov', '.avi', '.mkv'];
    const urlLower = url.toLowerCase();
    
    
    if (videoExtensions.some(ext => urlLower.includes(ext))) {
        return true;
    }
    
    
    if (urlLower.includes('youtube.com') || urlLower.includes('youtu.be')) {
        return true;
    }
    
    
    if (urlLower.includes('vimeo.com')) {
        return true;
    }
    
    return false;
}

function getYouTubeEmbedUrl(url) {
    let videoId = '';
    
    
    if (url.includes('youtu.be/')) {
        videoId = url.split('youtu.be/')[1].split('?')[0].split('&')[0];
    }
    
    else if (url.includes('youtube.com')) {
        if (url.includes('v=')) {
            videoId = url.split('v=')[1].split('&')[0];
        } else if (url.includes('/embed/')) {
            videoId = url.split('/embed/')[1].split('?')[0];
        }
    }
    
    return videoId ? `https://www.youtube.com/embed/${videoId}` : null;
}

function getVimeoEmbedUrl(url) {
    
    const match = url.match(/(?:vimeo\.com\/)(\d+)/);
    return match ? `https://player.vimeo.com/video/${match[1]}` : null;
}

function renderMedia(mediaUrl, exerciseName) {
    const mediaContainer = document.getElementById('exerciseMedia');
    const imageContainer = document.getElementById('exerciseImageContainer');
    const videoContainer = document.getElementById('exerciseVideoContainer');
    
    
    imageContainer.innerHTML = '';
    videoContainer.innerHTML = '';
    
    if (isVideoUrl(mediaUrl)) {
        
        const youtubeEmbed = getYouTubeEmbedUrl(mediaUrl);
        if (youtubeEmbed) {
            videoContainer.innerHTML = `
                <div class="ratio ratio-16x9">
                    <iframe 
                        src="${youtubeEmbed}" 
                        frameborder="0" 
                        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture" 
                        allowfullscreen
                        title="${exerciseName || 'Video del ejercicio'}">
                    </iframe>
                </div>
            `;
        }
        
        else {
            videoContainer.innerHTML = `
                <video class="w-100 rounded" controls>
                    <source src="${mediaUrl}" type="video/mp4">
                    <source src="${mediaUrl}" type="video/webm">
                    <source src="${mediaUrl}" type="video/ogg">
                    Tu navegador no soporta la reproducción de videos.
                </video>
            `;
        }
        imageContainer.classList.add('d-none');
        videoContainer.classList.remove('d-none');
    } else {
        
        imageContainer.innerHTML = `
            <img src="${mediaUrl}" class="img-fluid rounded" alt="${exerciseName || 'Imagen del ejercicio'}" style="max-height: 500px; object-fit: contain;">
        `;
        videoContainer.classList.add('d-none');
        imageContainer.classList.remove('d-none');
    }
    
    mediaContainer.classList.remove('d-none');
}

function addToRoutine() {
    if (currentExerciseId) {
        window.location.href = `/routines/create?exercise_id=${currentExerciseId}`;
    }
}

async function editExercise() {
    if (currentExerciseId) {
        window.location.href = `/exercises/${currentExerciseId}/edit`;
    }
}

async function deleteExercise() {
    if (!confirm('¿Estás seguro de eliminar este ejercicio? Esta acción no se puede deshacer.')) {
        return;
    }
    
    try {
        const token = localStorage.getItem('token');
        const response = await fetch(`/api/exercises/${currentExerciseId}`, {
            method: 'DELETE',
            headers: {
                'Authorization': 'Bearer ' + token
            }
        });

        if (response.ok) {
            alert('Ejercicio eliminado correctamente');
            window.location.href = '/exercises';
        } else {
            const data = await response.json();
            alert(data.error || 'Error al eliminar el ejercicio');
        }
    } catch (error) {
        alert('Error al eliminar el ejercicio');
        console.error('Error:', error);
    }
}

document.addEventListener('DOMContentLoaded', loadExercise);

