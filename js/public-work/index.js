const ownerLink = document.getElementById('ownerLink');
const ownerName = document.getElementById('ownerName');
const projectLink = document.getElementById('projectLink');
const projectName = document.getElementById('projectName');
const workTitle = document.getElementById('workTitle');
const workTitleMain = document.getElementById('workTitleMain');
const workViews = document.getElementById('workViews');
const workLikes = document.getElementById('workLikes');
const workDate = document.getElementById('workDate');
const workContent = document.getElementById('workContent');
const editBtn = document.getElementById('editBtn');

const pathParts = window.location.pathname.split('/');
const workSlug = pathParts[pathParts.length - 1];
const projectNameFromUrl = pathParts[pathParts.length - 2];
const username = pathParts[pathParts.length - 3];

let currentWork = null;

document.addEventListener('DOMContentLoaded', () => {
    loadWork();
    checkEditAccess();
});

async function loadWork() {
    try {
        const response = await fetch(`/api/v1/getPublicWork?username=${encodeURIComponent(username)}&project=${encodeURIComponent(projectNameFromUrl)}&slug=${encodeURIComponent(workSlug)}`);

        if (!response.ok) {
            if (response.status === 403) {
                showError('Это приватная работа');
                return;
            }
            if (response.status === 404) {
                showError('Работа не найдена');
                return;
            }
            throw new Error('Ошибка загрузки работы');
        }

        const work = await response.json();
        currentWork = work;

        document.title = `SelfBy: ${work.owner_name} - ${work.project_name} - ${work.title}`;

        ownerLink.href = `/u/${work.owner_name}`;
        ownerName.textContent = work.owner_name;
        projectLink.href = `/u/${work.owner_name}/${work.project_name}`;
        projectName.textContent = work.project_name;
        workTitle.textContent = work.title;

        workTitleMain.textContent = work.title;

        workViews.textContent = work.views;
        workLikes.textContent = work.likes;

        const createdDate = new Date(work.created_at);
        workDate.textContent = createdDate.toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        });

        if (work.content_html) {
            workContent.innerHTML = work.content_html;
        } else {
            workContent.innerHTML = '<p>Нет контента</p>';
        }

    } catch (error) {
        console.error('Ошибка загрузки работы:', error);
        showError('Ошибка загрузки работы');
    }
}

async function checkEditAccess() {
    try {
        const response = await fetch('/api/v1/getUser');
        if (!response.ok) return;

        const user = await response.json();

        if (user.username === username) {
            editBtn.style.display = 'block';
            editBtn.href = `/editor/${projectNameFromUrl}/${workSlug}`;
        }
    } catch (error) {}
}

function showError(message) {
    workTitleMain.textContent = message;
    workTitle.textContent = message;
    ownerName.textContent = '-';
    projectName.textContent = '-';
    workViews.textContent = '0';
    workLikes.textContent = '0';
    workDate.textContent = '-';
    workContent.innerHTML = '';
    document.querySelector('.work-meta').style.display = 'none';
    editBtn.style.display = 'none';
}