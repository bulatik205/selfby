(function () {
    const loadingEl = document.getElementById('printLoading');
    const articleEl = document.getElementById('printArticle');
    const bodyEl = document.getElementById('printBody');
    const link = document.getElementById('link');
    const toolbarEl = document.getElementById('printToolbar');
    const printBtn = document.getElementById('printActionBtn');
    const closeBtn = document.getElementById('closeBtn');

    const pathParts = window.location.pathname.split('/').filter(Boolean);
    const username = pathParts[1];
    const project = pathParts[2];
    const slug = pathParts[3];

    if (!username || !project || !slug) {
        showError('Не указаны параметры');
        return;
    }

    loadWork();

    async function loadWork() {
        try {
            const res = await fetch(
                `/api/v1/getPublicWork?username=${encodeURIComponent(username)}&project=${encodeURIComponent(project)}&slug=${encodeURIComponent(slug)}`
            );

            if (!res.ok) {
                if (res.status === 403) { showError('Это приватная работа'); return; }
                if (res.status === 404) { showError('Работа не найдена'); return; }
                throw new Error('Ошибка загрузки');
            }

            const work = await res.json();
            renderWork(work);

        } catch (err) {
            console.error(err);
            showError('Ошибка загрузки работы');
        }
    }

    function renderWork(work) {
        document.title = work.title;

        linkFull = "selfby.ru/u/" + username + "/" + project + "/" + slug;
        link.innerHTML = linkFull;
        bodyEl.innerHTML = work.content_html || '<p>Нет контента</p>';

        loadingEl.hidden = true;
        articleEl.hidden = false;
        toolbarEl.hidden = false;

        waitForResources().then(() => {
            window.print();
        });
    }

    function waitForResources() {
        const images = Array.from(document.images).map(img =>
            img.complete ? Promise.resolve() : new Promise(r => {
                img.onload = r;
                img.onerror = r;
            })
        );

        const fonts = document.fonts ? document.fonts.ready : Promise.resolve();

        return Promise.all([...images, fonts]);
    }

    function showError(message) {
        loadingEl.textContent = message;
        loadingEl.hidden = false;
        articleEl.hidden = true;
        toolbarEl.hidden = true;
    }

    function formatDate(iso) {
        return new Date(iso).toLocaleDateString('ru-RU', {
            year: 'numeric', month: 'long', day: 'numeric'
        });
    }

    printBtn.addEventListener('click', () => {
        window.print();
    });

    closeBtn.addEventListener('click', () => {
        if (window.history.length > 1) {
            window.history.back();
        } else {
            window.close();
        }
    });
})();