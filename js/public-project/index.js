const projectTitle = document.getElementById('projectTitle');
const projectDescription = document.getElementById('projectDescription');
const ownerLink = document.getElementById('ownerLink');
const ownerAvatar = document.getElementById('ownerAvatar');
const ownerName = document.getElementById('ownerName');
const projectType = document.getElementById('projectType');
const projectDate = document.getElementById('projectDate');
const worksList = document.getElementById('worksList');
const worksSection = document.querySelector('.works-section');
const indexContent = document.getElementById('indexContent');

const pathParts = window.location.pathname.split('/');
const projectName = pathParts[pathParts.length - 1];
const username = pathParts[pathParts.length - 2];

document.addEventListener('DOMContentLoaded', () => {
    loadProject();
});

async function loadProject() {
    try {
        const response = await fetch(`/api/v1/getPublicProject?username=${encodeURIComponent(username)}&project=${encodeURIComponent(projectName)}`);
        
        if (!response.ok) {
            if (response.status === 403) {
                showError('Это приватный проект');
                return;
            }
            if (response.status === 404) {
                showError('Проект не найден');
                return;
            }
            throw new Error('Ошибка загрузки проекта');
        }
        
        const project = await response.json();
        
        document.title = `SelfBy: ${project.owner_name} - ${project.name}`;
        projectTitle.textContent = project.name;
        projectDescription.textContent = project.description || 'Без описания';
        
        ownerLink.href = `/u/${project.owner_name}`;
        ownerAvatar.textContent = project.owner_name.charAt(0).toUpperCase();
        ownerName.textContent = project.owner_name;
        
        const typeLabel = project.type === 'public' ? 'Открытый' : 'Закрытый';
        projectType.textContent = typeLabel;
        projectType.className = `project-type ${project.type}`;
        
        const createdDate = new Date(project.created_at);
        projectDate.textContent = `Создан: ${createdDate.toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        })}`;
        
        if (project.index_work && project.index_work.content_html) {
            if (indexContent) {
                indexContent.innerHTML = project.index_work.content_html;
                indexContent.style.display = 'block';
            }
        }
        
        displayWorks(project.works);
        
    } catch (error) {
        console.error('Ошибка загрузки проекта:', error);
        showError('Ошибка загрузки проекта');
    }
}

function displayWorks(works) {
    if (!works || works.length === 0) {
        worksList.innerHTML = '<div class="empty-state">В этом проекте пока нет работ</div>';
        return;
    }
    
    worksList.innerHTML = '';
    
    works.forEach(work => {
        const workLink = document.createElement('a');
        workLink.className = 'work-card';
        workLink.href = `/u/${username}/${projectName}/${work.slug}`;
        
        const createdDate = new Date(work.created_at);
        const formattedDate = createdDate.toLocaleDateString('ru-RU', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        });
        
        workLink.innerHTML = `
            <span class="work-title">${work.title}</span>
            <span class="work-stats">
                <span>👁 ${work.views}</span>
                <span>❤ ${work.likes}</span>
                <span>📅 ${formattedDate}</span>
            </span>
        `;
        
        worksList.appendChild(workLink);
    });
}

function showError(message) {
    projectTitle.textContent = message;
    projectDescription.textContent = '';
    ownerLink.style.display = 'none';
    projectType.style.display = 'none';
    projectDate.style.display = 'none';
    worksList.innerHTML = '';
    if (indexContent) indexContent.style.display = 'none';
    worksSection.classList.add('hidden');
}