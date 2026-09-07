(function () {
    document.addEventListener('DOMContentLoaded', function () {
        initLikeButtons();

        const baseUrl = window.location.origin;
        const styles = [`${baseUrl}/styles/likes.css`];

        styles.forEach(stylePath => {
            const link = document.createElement('link');
            link.rel = 'stylesheet';
            link.href = stylePath;
            document.head.appendChild(link);
        });
    });

    function initLikeButtons() {
        const likeButtons = document.querySelectorAll('[data-like-btn]');

        likeButtons.forEach(button => {
            const elementType = button.getAttribute('data-element-type') || 'work';
            const existingId = button.getAttribute('data-element-id');

            if (existingId && existingId !== '') {
                setupLikeButton(button, existingId, elementType);
            } else {
                getElementId(elementType).then(elementId => {
                    if (elementId) {
                        button.setAttribute('data-element-id', elementId);
                        setupLikeButton(button, elementId, elementType);
                    }
                });
            }
        });
    }

    function setupLikeButton(button, elementId, elementType) {
        checkLikeState(button, elementId, elementType);

        button.addEventListener('click', function () {
            toggleLike(button, elementId, elementType);
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

    async function checkLikeState(button, elementId, elementType) {
        try {
            const response = await fetch(`/api/v1/checkLike?element_id=${elementId}&element_type=${elementType}`);

            if (response.ok) {
                const data = await response.json();
                updateButton(button, data);
            }
        } catch (error) {
            console.error('Ошибка проверки лайка:', error);
        }
    }

    async function toggleLike(button, elementId, elementType) {
        try {
            const response = await fetch('/api/v1/toggleLike', {
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
                updateButton(button, data);
            } else if (response.status === 401) {
                window.location.href = '/login';
            }
        } catch (error) {
            console.error('Ошибка переключения лайка:', error);
        }
    }

    function updateButton(button, data) {
        const countEl = button.querySelector('[data-like-count]');
        const iconEl = button.querySelector('[data-like-icon]');

        if (countEl) {
            countEl.textContent = data.count;
        }

        if (data.liked) {
            button.classList.add('liked');
            if (iconEl) {
                iconEl.textContent = '❤️';
            }
        } else {
            button.classList.remove('liked');
            if (iconEl) {
                iconEl.textContent = '🤍';
            }
        }
    }
})();