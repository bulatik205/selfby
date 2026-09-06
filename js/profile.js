(function() {
    const baseUrl = window.location.origin;
    
    const profileBtn = document.getElementById('profileBtn');
    
    if (!profileBtn) return;
    
    async function loadUserData() {
        try {
            const response = await fetch('/api/v1/getUser');
            
            if (!response.ok) {
                profileBtn.textContent = 'Войти';
                profileBtn.onclick = () => {
                    window.location.href = '/login';
                };
                return;
            }
            
            const user = await response.json();
            profileBtn.textContent = user.username;
            
            profileBtn.onclick = () => {
                window.location.href = `${baseUrl}/u/${user.username}`;
            };
            
            return user;
        } catch (error) {
            console.error('Ошибка загрузки пользователя:', error);
            profileBtn.textContent = 'Profile';
            profileBtn.onclick = () => {
                window.location.href = '/login';
            };
        }
    }
    
    loadUserData();
})();