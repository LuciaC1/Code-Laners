(function() {
    function updateHomeButtons() {
        const token = localStorage.getItem('token');
        const exploreBtn = document.getElementById('exploreBtn');
        const authButtons = document.getElementById('authButtons');
        const featuresSection = document.getElementById('featuresSection');
        
        if (!exploreBtn || !authButtons || !featuresSection) {
            setTimeout(updateHomeButtons, 100);
            return;
        }
        
        const hasToken = token && token.trim() !== '' && token !== 'null' && token !== 'undefined';

        const registerSection = document.getElementById('registerSection');
        const loginSection = document.getElementById('loginSection');

        if (hasToken) {
            exploreBtn.style.display = 'inline-block';
            authButtons.style.display = 'none';
            if (registerSection) registerSection.style.display = 'none';
            if (loginSection) loginSection.style.display = 'none';
            featuresSection.style.display = '';
        } else {
            exploreBtn.style.display = 'none';
            authButtons.style.display = 'flex';
            if (registerSection) registerSection.style.display = '';
            if (loginSection) loginSection.style.display = '';
            featuresSection.style.display = 'none';
        }
    }
    
    function init() {
        updateHomeButtons();
        window.addEventListener('storage', function(e) {
            if (e.key === 'token') {
                updateHomeButtons();
            }
        });
    }
    
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();

