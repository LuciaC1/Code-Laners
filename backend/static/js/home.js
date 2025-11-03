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
        
        if (token && token.trim() !== '' && token !== 'null' && token !== 'undefined') {
            exploreBtn.style.display = 'inline-block';
            authButtons.style.display = 'none';
            featuresSection.style.display = '';
        } else {
            exploreBtn.style.display = 'none';
            authButtons.style.display = 'flex';
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

