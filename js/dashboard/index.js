const modal = document.getElementById('newProjectModal');
const newProjectBtn = document.querySelector('.n-btn');
const closeBtn = document.querySelector('.modal-close');
const createProjectBtn = document.getElementById('createProjectBtn');
const projectNameInput = document.getElementById('projectName');
const projectTypeSelect = document.getElementById('projectType');
const projectDescriptionInput = document.getElementById('projectDescription');
const projectsList = document.getElementById('projectsList');
const profileBtn = document.getElementById('profileBtn');
const usernameSpan = document.getElementById('username');

let userData = null;
let currentProjects = [];

document.addEventListener('DOMContentLoaded', async () => {
    await loadUserData();
    await loadProjects();
    await loadDashboardStats();
});

async function loadUserData() {
    try {
        const response = await fetch('/api/v1/getUser');

        if (!response.ok) {
            if (response.status === 401) {
                window.location.href = '/login';
                return;
            }
            throw new Error('Ошибка загрузки данных пользователя');
        }

        const user = await response.json();
        userData = user;

        usernameSpan.textContent = user.username;
        profileBtn.textContent = user.username;

        const myProfileLink = document.getElementById('myProfileLink');
        if (myProfileLink) {
            myProfileLink.href = `/u/${user.username}`;
        }

        return user;
    } catch (error) {
        console.error('Ошибка:', error);
        usernameSpan.textContent = 'Пользователь';
        profileBtn.textContent = 'Profile';
    }
}

async function loadProjects() {
    try {
        const response = await fetch('/api/v1/getProjectsWithStats');
        const data = await response.json();

        if (!response.ok) {
            console.error('Ошибка загрузки проектов:', data.error);
            return;
        }

        currentProjects = data;
        projectsList.innerHTML = '';

        if (data.length === 0) {
            showEmptyState();
        } else {
            data.forEach(project => {
                addProjectToList(project, userData);
            });
        }

        updateStats(data);
    } catch (error) {
        console.error('Ошибка:', error);
    }
}

function updateStats(projects) {
    let totalWorks = 0;
    let totalViews = 0;
    let totalLikes = 0;

    projects.forEach(p => {
        totalWorks += p.works_count || 0;
        totalViews += p.total_views || 0;
        totalLikes += p.total_likes || 0;
    });

    document.getElementById('statProjects').textContent = projects.length;
    document.getElementById('statWorks').textContent = totalWorks;
    document.getElementById('statViews').textContent = totalViews;
    document.getElementById('statLikes').textContent = totalLikes;
}

newProjectBtn.addEventListener('click', () => {
    modal.classList.add('active');
    clearModalFields();
});

closeBtn.addEventListener('click', () => {
    modal.classList.remove('active');
});

modal.addEventListener('click', (e) => {
    if (e.target === modal) {
        modal.classList.remove('active');
    }
});

async function loadDashboardStats() {
    try {
        const response = await fetch('/api/v1/getDashboardStats');
        if (!response.ok) return;

        const data = await response.json();

        displayTopProjects(data.top_projects || []);
        displayRecentLikes(data.recent_likes || []);
    } catch (error) {
        console.error('Ошибка загрузки статистики:', error);
    }
}

function displayTopProjects(projects) {
    const container = document.getElementById('topProjectsList');
    if (!container) return;

    if (projects.length === 0) {
        container.innerHTML = '<div class="widget-empty">Нет данных</div>';
        return;
    }

    container.innerHTML = '';

    projects.forEach((project, index) => {
        const item = document.createElement('a');
        item.className = 'widget-item';
        item.href = `/editor/${project.name}`;

        const medal = index === 0 ? '🥇' : index === 1 ? '🥈' : index === 2 ? '🥉' : `${index + 1}.`;

        item.innerHTML = `
            <span class="widget-rank">${medal}</span>
            <span class="widget-title">${project.name}</span>
            <span class="widget-stats">
                <span>👁 ${project.total_views}</span>
                <span>❤ ${project.total_likes}</span>
            </span>
        `;

        container.appendChild(item);
    });
}

function displayRecentLikes(likes) {
    const container = document.getElementById('recentLikesList');
    if (!container) return;

    if (likes.length === 0) {
        container.innerHTML = '<div class="widget-empty">Пока нет лайков</div>';
        return;
    }

    container.innerHTML = '';

    likes.forEach(like => {
        const item = document.createElement('a');
        item.className = 'widget-item';

        if (like.element_type === 'work' && like.work_slug) {
            item.href = `/editor/${like.project_name}/${like.work_slug}`;
        } else {
            item.href = `/editor/${like.project_name}`;
        }

        const date = new Date(like.created_at);
        const formattedDate = date.toLocaleDateString('ru-RU', {
            day: 'numeric',
            month: 'short'
        });

        const target = like.element_type === 'work' && like.work_title
            ? `→ ${like.work_title}`
            : `→ проект ${like.project_name}`;

        item.innerHTML = `
            <span class="widget-title">${like.liker_name}</span>
            <span class="widget-subtitle">${target}</span>
            <span class="widget-date">${formattedDate}</span>
        `;

        container.appendChild(item);
    });
}

createProjectBtn.addEventListener('click', async () => {
    const projectData = {
        name: projectNameInput.value.trim(),
        type: projectTypeSelect.value,
        description: projectDescriptionInput.value.trim()
    };

    if (!projectData.name) {
        showError('Введите название проекта');
        return;
    }

    if (!projectData.type) {
        showError('Выберите тип проекта');
        return;
    }

    try {
        const response = await fetch('/api/v1/newProject', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(projectData)
        });

        const data = await response.json();

        if (response.ok) {
            modal.classList.remove('active');
            clearModalFields();
            removeEmptyState();

            const newProject = {
                id: data.id,
                name: data.name,
                description: data.description,
                type: data.type,
                created_at: data.created_at,
                works_count: 0,
                total_views: 0,
                total_likes: 0
            };

            currentProjects.push(newProject);
            addProjectToList(newProject, userData);
            updateStats(currentProjects);
        } else {
            showError(data.error || 'Ошибка при создании проекта');
        }
    } catch (error) {
        console.error('Ошибка:', error);
        showError('Ошибка соединения с сервером');
    }
});

async function deleteProject(projectName, event) {
    event.preventDefault();
    event.stopPropagation();

    if (!confirm(`Удалить проект "${projectName}"? Все работы будут удалены. Это действие нельзя отменить.`)) {
        return;
    }

    try {
        const response = await fetch(`/api/v1/deleteProject?name=${encodeURIComponent(projectName)}`, {
            method: 'DELETE'
        });

        const data = await response.json();

        if (response.ok) {
            event.target.closest('.project').remove();
            currentProjects = currentProjects.filter(p => p.name !== projectName);

            if (projectsList.children.length === 0) {
                showEmptyState();
            }

            updateStats(currentProjects);
        } else {
            alert(data.error || 'Ошибка удаления');
        }
    } catch (error) {
        console.error('Ошибка удаления:', error);
        alert('Ошибка соединения с сервером');
    }
}

function addProjectToList(project, userData) {
    const projectDiv = document.createElement('div');
    projectDiv.className = 'project';

    const username = userData?.username || '';
    const isPublic = project.type === 'public';
    const lockIcon = isPublic ? '🔓' : '🔒';
    const lockTitle = isPublic ? 'Сделать закрытым' : 'Сделать открытым';

    projectDiv.innerHTML = `
        <a href="/editor/${project.name}" class="project-link">${project.name}</a>
        <a href="/u/${username}/${project.name}" class="project-link icon-btn">
            <img src="../images/view.png" alt="Просмотр">
        </a>
        <button class="project-link icon-btn" onclick="toggleProjectType('${project.name}', this, event)" title="${lockTitle}">
            ${lockIcon}
        </button>
        <button class="project-link icon-btn" onclick="deleteProject('${project.name}', event)" title="Удалить проект">
            🗑️
        </button>
    `;

    projectsList.appendChild(projectDiv);
}

async function toggleProjectType(projectName, button, event) {
    event.preventDefault();
    event.stopPropagation();

    try {
        const response = await fetch('/api/v1/toggleProjectType', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                project_name: projectName
            })
        });

        const data = await response.json();

        if (response.ok) {
            if (data.type === 'public') {
                button.textContent = '🔓';
                button.title = 'Сделать закрытым';
            } else {
                button.textContent = '🔒';
                button.title = 'Сделать открытым';
            }

            const project = currentProjects.find(p => p.name === projectName);
            if (project) {
                project.type = data.type;
            }
        } else {
            alert(data.error || 'Ошибка изменения типа');
        }
    } catch (error) {
        console.error('Ошибка:', error);
        alert('Ошибка соединения с сервером');
    }
}

function clearModalFields() {
    projectNameInput.value = '';
    projectTypeSelect.value = '';
    projectDescriptionInput.value = '';
    removeError();
}

function showError(message) {
    removeError();
    const errorDiv = document.createElement('div');
    errorDiv.className = 'error-message';
    errorDiv.textContent = message;
    errorDiv.style.cssText = `
        color: #e74c3c;
        font-size: 12px;
        margin-top: 10px;
        padding: 8px;
        background: #fde8e8;
        border-radius: 8px;
        text-align: center;
    `;
    document.querySelector('.modal-body').appendChild(errorDiv);
}

function showEmptyState() {
    const emptyDiv = document.createElement('div');
    emptyDiv.className = 'empty-state';
    emptyDiv.textContent = 'У вас пока нет проектов';
    emptyDiv.style.cssText = `
        color: #999;
        font-size: 14px;
        text-align: center;
        padding: 20px;
    `;
    projectsList.appendChild(emptyDiv);
}

function removeEmptyState() {
    const emptyState = document.querySelector('.empty-state');
    if (emptyState) {
        emptyState.remove();
    }
}

function removeError() {
    const existingError = document.querySelector('.error-message');
    if (existingError) {
        existingError.remove();
    }
}

document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && modal.classList.contains('active')) {
        modal.classList.remove('active');
    }
});