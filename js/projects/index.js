const projectsList = document.getElementById('projectsList');
const sortTabs = document.querySelectorAll('.sort-tab');
const modal = document.getElementById('newProjectModal');
const newProjectBtn = document.getElementById('newProjectBtn');
const closeBtn = document.querySelector('.modal-close');
const createProjectBtn = document.getElementById('createProjectBtn');
const projectNameInput = document.getElementById('projectName');
const projectTypeSelect = document.getElementById('projectType');
const projectDescriptionInput = document.getElementById('projectDescription');

let currentSort = 'created';

document.addEventListener('DOMContentLoaded', () => {
    loadProjects(currentSort);
    
    sortTabs.forEach(tab => {
        tab.addEventListener('click', function() {
            sortTabs.forEach(t => t.classList.remove('active'));
            this.classList.add('active');
            currentSort = this.getAttribute('data-sort');
            loadProjects(currentSort);
        });
    });
});

async function loadProjects(sortBy = 'created') {
    try {
        const response = await fetch(`/api/v1/getProjectsWithStats?sort=${sortBy}`);
        
        if (!response.ok) {
            if (response.status === 401) {
                window.location.href = '/login';
                return;
            }
            throw new Error('Ошибка загрузки проектов');
        }
        
        const projects = await response.json();
        
        if (projects.length === 0) {
            projectsList.innerHTML = '<div class="empty-state">У вас пока нет проектов</div>';
            return;
        }
        
        projectsList.innerHTML = '';
        
        projects.forEach(project => {
            addProjectToList(project);
        });
        
    } catch (error) {
        console.error('Ошибка загрузки проектов:', error);
        projectsList.innerHTML = '<div class="empty-state">Ошибка загрузки</div>';
    }
}

function addProjectToList(project) {
    const projectCard = document.createElement('a');
    projectCard.className = 'project-card';
    projectCard.href = `/editor/${project.name}`;
    
    const typeLabel = project.type === 'public' ? 'Открытый' : 'Закрытый';
    const createdDate = new Date(project.created_at).toLocaleDateString('ru-RU');
    
    projectCard.innerHTML = `
        <div class="project-info">
            <div class="project-name">${project.name}</div>
            <div class="project-description">${project.description || 'Без описания'}</div>
        </div>
        <div class="project-meta">
            <span class="project-type ${project.type}">${typeLabel}</span>
            <span class="project-date">${createdDate}</span>
        </div>
        <div class="project-stats">
            <span>📄 ${project.works_count} работ</span>
            <span>👁 ${project.total_views} просмотров</span>
            <span>❤ ${project.total_likes} лайков</span>
        </div>
    `;
    
    projectsList.appendChild(projectCard);
}

newProjectBtn.addEventListener('click', () => {
    modal.classList.add('active');
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
        alert('Введите название проекта');
        return;
    }

    if (!projectData.type) {
        alert('Выберите тип проекта');
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
            projectNameInput.value = '';
            projectTypeSelect.value = '';
            projectDescriptionInput.value = '';
            loadProjects(currentSort);
        } else {
            alert(data.error || 'Ошибка при создании проекта');
        }
    } catch (error) {
        console.error('Ошибка создания проекта:', error);
        alert('Ошибка соединения с сервером');
    }
});

document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape' && modal.classList.contains('active')) {
        modal.classList.remove('active');
    }
});