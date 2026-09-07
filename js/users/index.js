const usersList = document.getElementById('usersList');
const sortTabs = document.querySelectorAll('.sort-tab');

let currentSort = 'views';

document.addEventListener('DOMContentLoaded', () => {
    loadTopUsers(currentSort);
    
    sortTabs.forEach(tab => {
        tab.addEventListener('click', function() {
            sortTabs.forEach(t => t.classList.remove('active'));
            this.classList.add('active');
            
            currentSort = this.getAttribute('data-sort');
            
            loadTopUsers(currentSort);
        });
    });
});

async function loadTopUsers(sortBy = 'views') {
    try {
        const response = await fetch(`/api/v1/getUsers?sort=${sortBy}`);
        
        if (!response.ok) {
            throw new Error('Ошибка загрузки пользователей');
        }
        
        const users = await response.json();
        
        if (users.length === 0) {
            usersList.innerHTML = '<div class="empty-state">Пока нет пользователей</div>';
            return;
        }
        
        usersList.innerHTML = '';
        
        users.forEach((user, index) => {
            addUserToList(user, index);
        });
        
    } catch (error) {
        console.error('Ошибка загрузки пользователей:', error);
        usersList.innerHTML = '<div class="empty-state">Ошибка загрузки</div>';
    }
}

function addUserToList(user, index) {
    const userCard = document.createElement('a');
    userCard.className = 'user-card';
    userCard.href = `/u/${user.username}`;
    
    const rank = index + 1;
    const rankClass = rank <= 3 ? `top-${rank}` : '';
    
    userCard.innerHTML = `
        <span class="user-rank ${rankClass}">#${rank}</span>
        <span class="user-avatar">${user.username.charAt(0).toUpperCase()}</span>
        <div class="user-info">
            <div class="user-name">${user.username}</div>
            <div class="user-stats">
                <span>📂 ${user.projects_count} проектов</span>
                <span>📄 ${user.works_count} работ</span>
                <span>👁 ${user.total_views} просмотров</span>
                <span>❤ ${user.total_likes} лайков</span>
            </div>
        </div>
    `;
    
    usersList.appendChild(userCard);
}