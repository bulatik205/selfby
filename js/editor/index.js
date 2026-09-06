const modal = document.getElementById('newWorkModal');
const newWorkBtn = document.querySelector('.n-btn');
const closeBtn = document.querySelector('.modal-close');
const createWorkBtn = document.getElementById('createWorkBtn');
const workTitleInput = document.getElementById('workTitle');
const workSlugDisplay = document.getElementById('work-slug');
const worksList = document.getElementById('worksList');
const profileBtn = document.getElementById('profileBtn');

const pathParts = window.location.pathname.split('/');
const projectName = pathParts[pathParts.length - 1];

let slugCheckTimeout;

document.addEventListener('DOMContentLoaded', () => {
    loadUserData();
    loadWorks();
});

async function loadUserData() {
    try {
        const response = await fetch('/api/v1/getUser');
        if (!response.ok) throw new Error('Ошибка загрузки');
        
        const user = await response.json();
        profileBtn.textContent = user.username;
    } catch (error) {
        profileBtn.textContent = 'Profile';
    }
}

async function loadWorks() {
    try {
        const response = await fetch(`/api/v1/getWorks?project=${encodeURIComponent(projectName)}`);
        if (!response.ok) throw new Error('Ошибка загрузки работ');
        
        const works = await response.json();
        worksList.innerHTML = '';
        
        if (works.length === 0) {
            showEmptyState();
        } else {
            works.forEach(addWorkToList);
        }
    } catch (error) {
        worksList.innerHTML = '<div class="empty-state">Ошибка загрузки работ</div>';
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
    `;
    
    worksList.appendChild(workDiv);
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