(function () {
    let viewSent = false;
    
    document.addEventListener('DOMContentLoaded', function () {
        initViews();
    });

    function initViews() {
        const viewElements = document.querySelectorAll('[data-view]');
        
        viewElements.forEach(element => {
            const elementType = element.getAttribute('data-element-type') || 'work';
            const existingId = element.getAttribute('data-element-id');

            if (existingId && existingId !== '') {
                addView(existingId, elementType);
            } else {
                getElementId(elementType).then(elementId => {
                    if (elementId) {
                        element.setAttribute('data-element-id', elementId);
                        addView(elementId, elementType);
                    }
                });
            }
        });
    }

    async function getElementId(elementType) {
        const path = window.location.pathname;
        const parts = path.split('/').filter(p => p);

        if (parts.length >= 4 && parts[0] === 'u') {
            const username = parts[1];
            const projectName = parts[2];
            const workSlug = parts[3];

            try {
                const response = await fetch(`/api/v1/getPublicWork?username=${username}&project=${projectName}&slug=${workSlug}`);
                if (response.ok) {
                    const work = await response.json();
                    return work.id;
                }
            } catch (error) {
                console.error('Ошибка получения ID работы:', error);
            }
        }

        if (parts.length >= 3 && parts[0] === 'u') {
            const username = parts[1];
            const projectName = parts[2];

            try {
                const response = await fetch(`/api/v1/getPublicProject?username=${username}&project=${projectName}`);
                if (response.ok) {
                    const project = await response.json();
                    return project.id;
                }
            } catch (error) {
                console.error('Ошибка получения ID проекта:', error);
            }
        }

        if (parts.length >= 2 && parts[0] === 'u') {
            const username = parts[1];

            try {
                const response = await fetch(`/api/v1/getProfile?username=${username}`);
                if (response.ok) {
                    const profile = await response.json();
                    return profile.id;
                }
            } catch (error) {
                console.error('Ошибка получения ID пользователя:', error);
            }
        }

        return null;
    }

    async function addView(elementId, elementType) {
        if (viewSent) return; 
        
        try {
            viewSent = true; 
            
            const response = await fetch('/api/v1/addView', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    element_id: parseInt(elementId),
                    element_type: elementType
                })
            });

            if (response.ok) {
                const data = await response.json();
                console.log('Просмотр:', data.viewed ? 'учтен' : 'пропущен (меньше часа)');
                
                updateViewCount(data.count);
            }
        } catch (error) {
            console.error('Ошибка добавления просмотра:', error);
            viewSent = false; 
        }
    }

    function updateViewCount(count) {
        const countElements = document.querySelectorAll('[data-view-count]');
        countElements.forEach(el => {
            el.textContent = count;
        });
    }
})();