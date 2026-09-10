const modal = document.getElementById('newWorkModal');
const newWorkBtn = document.querySelector('.n-btn');
const closeBtn = document.querySelector('.modal-close');
const createWorkBtn = document.getElementById('createWorkBtn');
const workTitleInput = document.getElementById('workTitle');
const workSlugDisplay = document.getElementById('work-slug');
const worksList = document.getElementById('worksList');
const profileBtn = document.getElementById('profileBtn');
const projectTitle = document.getElementById('projectTitle');
const worksCountEl = document.getElementById('worksCount');
const totalViewsEl = document.getElementById('totalViews');
const totalLikesEl = document.getElementById('totalLikes');
const projectDateEl = document.getElementById('projectDate');
const viewProjectBtn = document.getElementById('viewProjectBtn');
const deleteProjectBtn = document.getElementById('deleteProjectBtn');
const toggleTypeBtn = document.getElementById('toggleTypeBtn');
const projectTypeEl = document.getElementById('projectType');

const pathParts = window.location.pathname.split('/');
const projectName = pathParts[pathParts.length - 1];

let slugCheckTimeout;
let currentWorks = [];
let currentUser = null;

document.addEventListener('DOMContentLoaded', () => {
    projectTitle.textContent = projectName;
    document.title = `SelfBy: ${projectName}`;
    
    loadUserData();
    loadWorks();
});

async function loadUserData() {
    try {
        const response = await fetch('/api/v1/getUser');
        if (!response.ok) throw new Error('Ошибка загрузки');
        
        const user = await response.json();
        currentUser = user;
        profileBtn.textContent = user.username;
        
        if (viewProjectBtn) {
            viewProjectBtn.href = `/u/${user.username}/${projectName}`;
        }
        
        loadProjectType(user.username);
    } catch (error) {
        profileBtn.textContent = 'Profile';
    }
}

async function loadProjectType(username) {
    try {
        const response = await fetch(`/api/v1/getPublicProject?username=${encodeURIComponent(username)}&project=${encodeURIComponent(projectName)}`);
        if (!response.ok) return;
        
        const project = await response.json();
        
        if (toggleTypeBtn) {
            if (project.type === 'public') {
                toggleTypeBtn.textContent = '🔓';
                toggleTypeBtn.title = 'Сделать закрытым';
                toggleTypeBtn.classList.add('public');
                toggleTypeBtn.classList.remove('private');
            } else {
                toggleTypeBtn.textContent = '🔒';
                toggleTypeBtn.title = 'Сделать открытым';
                toggleTypeBtn.classList.add('private');
                toggleTypeBtn.classList.remove('public');
            }
        }
        
        if (projectTypeEl) {
            if (project.type === 'public') {
                projectTypeEl.textContent = 'Открытый';
                projectTypeEl.className = 'topic-value public';
            } else {
                projectTypeEl.textContent = 'Закрытый';
                projectTypeEl.className = 'topic-value private';
            }
        }
    } catch (error) {
        console.error('Ошибка загрузки типа проекта:', error);
    }
}

async function loadWorks() {
    try {
        const response = await fetch(`/api/v1/getWorks?project=${encodeURIComponent(projectName)}`);
        if (!response.ok) throw new Error('Ошибка загрузки работ');
        
        const works = await response.json();
        currentWorks = works;
        
        worksList.innerHTML = '';
        
        if (works.length === 0) {
            showEmptyState();
        } else {
            works.forEach(addWorkToList);
        }
        
        updateProjectStats(works);
        
    } catch (error) {
        worksList.innerHTML = '<div class="empty-state">Ошибка загрузки работ</div>';
    }
}

function updateProjectStats(works) {
    worksCountEl.textContent = works.length;
    
    const totalViews = works.reduce((sum, work) => sum + (work.views || 0), 0);
    totalViewsEl.textContent = totalViews;
    
    const totalLikes = works.reduce((sum, work) => sum + (work.likes || 0), 0);
    totalLikesEl.textContent = totalLikes;
    
    if (works.length > 0 && works[0].created_at) {
        const firstWorkDate = new Date(works[0].created_at);
        projectDateEl.textContent = firstWorkDate.toLocaleDateString('ru-RU');
    } else {
        projectDateEl.textContent = '-';
    }
}

newWorkBtn.addEventListener('click', () => {
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

workTitleInput.addEventListener('input', () => {
    const title = workTitleInput.value.trim();
    
    if (title) {
        const suggestedSlug = generateSlugClient(title);
        workSlugDisplay.textContent = suggestedSlug;
        workSlugDisplay.className = '';
        
        clearTimeout(slugCheckTimeout);
        slugCheckTimeout = setTimeout(() => {
            checkSlugAvailability(title);
        }, 500);
    } else {
        workSlugDisplay.textContent = 'Введите название...';
        workSlugDisplay.className = '';
    }
});

async function checkSlugAvailability(title) {
    try {
        const requestBody = {
            title: title,
            project_name: projectName
        };
        
        const response = await fetch('/api/v1/checkSlug', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody)
        });
        
        const data = await response.json();
        
        if (data.available) {
            workSlugDisplay.textContent = `${data.slug} ✓`;
            workSlugDisplay.className = 'available';
        } else {
            workSlugDisplay.textContent = `${data.slug} - ${data.message}`;
            workSlugDisplay.className = 'unavailable';
        }
    } catch (error) {
        console.error('Ошибка проверки slug:', error);
        workSlugDisplay.textContent = `Ошибка: ${error.message}`;
        workSlugDisplay.className = 'error';
    }
}

createWorkBtn.addEventListener('click', async () => {
    const title = workTitleInput.value.trim();
    
    if (!title) {
        showError('Введите название работы');
        return;
    }
    
    try {
        const requestBody = {
            title: title,
            project_name: projectName
        };
        
        const response = await fetch('/api/v1/newWork', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody)
        });
        
        const data = await response.json();
        
        if (response.ok) {
            modal.classList.remove('active');
            clearModalFields();
            removeEmptyState();
            addWorkToList(data);
            
            currentWorks.push(data);
            updateProjectStats(currentWorks);
        } else {
            showError(data.error || 'Ошибка при создании работы');
        }
    } catch (error) {
        console.error('Ошибка создания работы:', error);
        showError('Ошибка соединения с сервером');
    }
});

function generateSlugClient(title) {
    const slug = title
        .toLowerCase()
        .replace(/[^a-z0-9\s-]/g, ' ')
        .trim()
        .replace(/\s+/g, '-')
        .replace(/-+/g, '-')
        .replace(/^-|-$/g, '');
    
    return slug || 'work';
}

function addWorkToList(work) {
    const workDiv = document.createElement('div');
    workDiv.className = 'work';
    
    workDiv.innerHTML = `
        <a href="/editor/${projectName}/${work.slug}" class="work-link">
            <span class="work-title">${work.title}</span>
            <span class="work-slug">/${projectName}/${work.slug}</span>
        </a>
        <span class="work-stats">
            <span>👁 ${work.views}</span>
            <span>❤ ${work.likes}</span>
        </span>
        <button class="delete-btn" onclick="deleteWork('${work.slug}', event)" title="Удалить">
            🗑️
        </button>
    `;
    
    worksList.appendChild(workDiv);
}

async function deleteWork(slug, event) {
    event.preventDefault();
    event.stopPropagation();
    
    if (!confirm(`Удалить работу "${slug}"? Это действие нельзя отменить.`)) {
        return;
    }
    
    try {
        const response = await fetch(`/api/v1/deleteWork?project=${encodeURIComponent(projectName)}&slug=${encodeURIComponent(slug)}`, {
            method: 'DELETE'
        });
        
        const data = await response.json();
        
        if (response.ok) {
            currentWorks = currentWorks.filter(w => w.slug !== slug);
            
            worksList.innerHTML = '';
            if (currentWorks.length === 0) {
                showEmptyState();
            } else {
                currentWorks.forEach(addWorkToList);
            }
            
            updateProjectStats(currentWorks);
        } else {
            alert(data.error || 'Ошибка удаления');
        }
    } catch (error) {
        console.error('Ошибка удаления:', error);
        alert('Ошибка соединения с сервером');
    }
}

if (toggleTypeBtn) {
    toggleTypeBtn.addEventListener('click', async () => {
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
                    toggleTypeBtn.textContent = '🔓';
                    toggleTypeBtn.title = 'Сделать закрытым';
                    toggleTypeBtn.classList.add('public');
                    toggleTypeBtn.classList.remove('private');
                } else {
                    toggleTypeBtn.textContent = '🔒';
                    toggleTypeBtn.title = 'Сделать открытым';
                    toggleTypeBtn.classList.add('private');
                    toggleTypeBtn.classList.remove('public');
                }
                
                if (projectTypeEl) {
                    if (data.type === 'public') {
                        projectTypeEl.textContent = 'Открытый';
                        projectTypeEl.className = 'topic-value public';
                    } else {
                        projectTypeEl.textContent = 'Закрытый';
                        projectTypeEl.className = 'topic-value private';
                    }
                }
            } else {
                alert(data.error || 'Ошибка изменения типа');
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Ошибка соединения с сервером');
        }
    });
}

if (toggleTypeBtn) {
    toggleTypeBtn.addEventListener('click', async () => {
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
                    toggleTypeBtn.textContent = '🔓';
                    toggleTypeBtn.title = 'Сделать закрытым';
                    toggleTypeBtn.classList.add('public');
                    toggleTypeBtn.classList.remove('private');
                } else {
                    toggleTypeBtn.textContent = '🔒';
                    toggleTypeBtn.title = 'Сделать открытым';
                    toggleTypeBtn.classList.add('private');
                    toggleTypeBtn.classList.remove('public');
                }
            } else {
                alert(data.error || 'Ошибка изменения типа');
            }
        } catch (error) {
            console.error('Ошибка:', error);
            alert('Ошибка соединения с сервером');
        }
    });
}

function clearModalFields() {
    workTitleInput.value = '';
    workSlugDisplay.textContent = 'Введите название...';
    workSlugDisplay.className = '';
    removeError();
}

function showError(message) {
    removeError();
    const errorDiv = document.createElement('div');
    errorDiv.className = 'error-message';
    errorDiv.textContent = message;
    document.querySelector('.modal-body').appendChild(errorDiv);
}

function showEmptyState() {
    const emptyDiv = document.createElement('div');
    emptyDiv.className = 'empty-state';
    emptyDiv.textContent = 'В этом проекте пока нет работ';
    worksList.appendChild(emptyDiv);
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