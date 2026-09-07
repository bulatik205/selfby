const usernameEl = document.getElementById('username');
const avatarEl = document.getElementById('avatar');
const userIdEl = document.getElementById('userId');
const createdAtEl = document.getElementById('createdAt');
const projectsCountEl = document.getElementById('projectsCount');
const worksCountEl = document.getElementById('worksCount');
const totalLikesEl = document.getElementById('totalLikes');
const totalViewsEl = document.getElementById('totalViews');
const profileBtn = document.getElementById('profileBtn');
const projectsListEl = document.getElementById('projectsList');

// Получаем username из URL
const pathParts = window.location.pathname.split('/');
const username = pathParts[pathParts.length - 1];

let profile = null;

document.addEventListener('DOMContentLoaded', () => {
    loadProfile();
    loadCurrentUser();
});

function displayProjects(projects) {
    if (!projects || projects.length === 0) {
        projectsListEl.innerHTML = '<div class="empty-state">Нет проектов</div>';
        return;
    }
    
    projectsListEl.innerHTML = '';
    
    projects.forEach(project => {
        const projectCard = document.createElement('div');
        projectCard.className = 'project-card';
        
        const typeLabel = project.type === 'public' ? 'Открытый' : 'Закрытый';
        
        const createdDate = new Date(project.created_at);
        const formattedDate = createdDate.toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
        
        projectCard.innerHTML = `
            <div class="project-info">
                <a href="/u/${username}/${project.name}" class="project-name">${project.name}</a>
                <span class="project-description">${project.description || 'Без описания'}</span>
            </div>
            <div class="project-meta">
                <span class="project-type ${project.type}">${typeLabel}</span>
                <span class="project-works">📄 ${project.works_count} работ</span>
            </div>
            <div class="project-date">
                Создан: ${formattedDate}
            </div>
        `;
        
        projectsListEl.appendChild(projectCard);
    });
}

async function loadProfile() {
    try {
        const response = await fetch(`/api/v1/getProfile?username=${encodeURIComponent(username)}`);
        
        if (!response.ok) {
            throw new Error('Ошибка загрузки профиля');
        }
        
        profile = await response.json();
        
        // Заполняем данные
        document.title = `SelfBy: ${profile.username}`;
        usernameEl.textContent = profile.username;
        avatarEl.textContent = profile.username.charAt(0).toUpperCase();
        userIdEl.textContent = profile.id;
        
        // Форматируем дату
        const createdDate = new Date(profile.created_at);
        createdAtEl.textContent = createdDate.toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });
        
        // Статистика (без разделения на public/private)
        projectsCountEl.textContent = profile.projects_count;
        worksCountEl.textContent = profile.works_count;
        totalLikesEl.textContent = profile.total_likes;
        totalViewsEl.textContent = profile.total_views;
        
        // Отображаем проекты
        displayProjects(profile.projects);
        
    } catch (error) {
        console.error('Ошибка загрузки профиля:', error);
        usernameEl.textContent = 'Пользователь не найден';
    }
}

async function loadCurrentUser() {
    try {
        const response = await fetch('/api/v1/getUser');
        if (!response.ok) return;
        
        const user = await response.json();
        profileBtn.textContent = user.username;
        
        if (user.username === username) {
            profileBtn.textContent = 'Dashboard';
            profileBtn.onclick = () => {
                window.location.href = '/dashboard';
            };
        } else {
            profileBtn.onclick = () => {
                window.location.href = `/u/${user.username}`;
            };
        }
    } catch (error) {
        profileBtn.textContent = 'Profile';
    }
}