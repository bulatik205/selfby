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

document.addEventListener('DOMContentLoaded', async () => {
    await loadUserData();
    await loadProjects();
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
        
        return user;
    } catch (error) {
        console.error('Ошибка:', error);
        usernameSpan.textContent = 'Пользователь';
        profileBtn.textContent = 'Profile';
    }
}

async function loadProjects() {
    try {
        const response = await fetch('/api/v1/getProjects');
        const data = await response.json();

        if (response.ok) {
            projectsList.innerHTML = '';
            
            if (data.length === 0) {
                showEmptyState();
            } else {
                data.forEach(project => {
                    addProjectToList(project, userData);
                });
            }
        } else {
            console.error('Ошибка загрузки проектов:', data.error);
        }
    } catch (error) {
        console.error('Ошибка:', error);
    }
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
            addProjectToList(data, userData);
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
            
            if (projectsList.children.length === 0) {
                showEmptyState();
            }
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
    
    projectDiv.innerHTML = `
        <a href="/editor/${project.name}" class="project-link">${project.name}</a>
        <a href="/u/${username}/${project.name}" class="project-link icon-btn">
            <img src="../images/view.png" alt="Просмотр">
        </a>
        <button class="project-link icon-btn" onclick="deleteProject('${project.name}', event)" title="Удалить проект">
            🗑️
        </button>
    `;
    
    projectsList.appendChild(projectDiv);
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