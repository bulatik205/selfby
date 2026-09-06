const workTitle = document.getElementById('workTitle');
const projectNameEl = document.getElementById('projectName');
const workSlugEl = document.getElementById('workSlug');
const workViews = document.getElementById('workViews');
const workLikes = document.getElementById('workLikes');
const workCreatedAt = document.getElementById('workCreatedAt');
const contentMD = document.getElementById('contentMD');
const saveWorkBtn = document.getElementById('saveWorkBtn');
const viewWorkBtn = document.getElementById('viewWorkBtn');
const profileBtn = document.getElementById('profileBtn');

const pathParts = window.location.pathname.split('/');
const workSlug = pathParts[pathParts.length - 1];
const projectName = pathParts[pathParts.length - 2];

let currentWork = null;

document.addEventListener('DOMContentLoaded', () => {
    loadWork();
});

async function loadWork() {
    try {
        const response = await fetch(`/api/v1/getWork?project=${encodeURIComponent(projectName)}&slug=${encodeURIComponent(workSlug)}`);
        
        if (!response.ok) {
            if (response.status === 401) {
                window.location.href = '/login';
                return;
            }
            throw new Error('Ошибка загрузки работы');
        }
        
        const work = await response.json();
        currentWork = work;
        
        document.title = `${work.owner_name}: ${work.project_name} - ${work.title}`;
        workTitle.textContent = work.title;
        projectNameEl.textContent = work.project_name;
        workSlugEl.textContent = work.slug;
        workViews.textContent = work.views;
        workLikes.textContent = work.likes;
        workCreatedAt.textContent = new Date(work.created_at).toLocaleDateString('ru-RU');
        contentMD.value = work.content_md || '';
        
        viewWorkBtn.href = `/u/${work.owner_name}/${work.project_name}/${work.slug}`;
        
    } catch (error) {
        console.error('Ошибка загрузки работы:', error);
        workTitle.textContent = 'Ошибка загрузки';
    }
}

saveWorkBtn.addEventListener('click', async () => {
    if (!currentWork) return;
    
    const content = contentMD.value;
    
    try {
        const response = await fetch('/api/v1/updateWork', {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                project_name: projectName,
                slug: workSlug,
                content_md: content
            })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            showSaveSuccess();
        } else {
            showSaveError(data.error || 'Ошибка при сохранении');
        }
    } catch (error) {
        console.error('Ошибка сохранения:', error);
        showSaveError('Ошибка соединения с сервером');
    }
});

function showSaveSuccess() {
    saveWorkBtn.textContent = '✓ Сохранено';
    saveWorkBtn.style.background = '#18A64A';
    setTimeout(() => {
        saveWorkBtn.textContent = 'Сохранить';
        saveWorkBtn.style.background = '';
    }, 2000);
}

function showSaveError(message) {
    saveWorkBtn.textContent = '✗ Ошибка';
    saveWorkBtn.style.background = '#e74c3c';
    alert(message);
    setTimeout(() => {
        saveWorkBtn.textContent = 'Сохранить';
        saveWorkBtn.style.background = '';
    }, 2000);
}

document.addEventListener('keydown', (e) => {
    if ((e.ctrlKey || e.metaKey) && e.key === 's') {
        e.preventDefault();
        saveWorkBtn.click();
    }
});