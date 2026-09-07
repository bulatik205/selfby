(function () {
    document.addEventListener('DOMContentLoaded', function () {
        buildBreadcrumbs();

        const baseUrl = window.location.origin;
        const styles = [`${baseUrl}/styles/breadcrumbs.css`];

        styles.forEach(stylePath => {
            const link = document.createElement('link');
            link.rel = 'stylesheet';
            link.href = stylePath;
            document.head.appendChild(link);
        });
    });

    function buildBreadcrumbs() {
        const container = document.getElementById('breadcrumbs');
        if (!container) return;

        const path = window.location.pathname;
        const parts = path.split('/').filter(p => p);

        let html = '';
        let currentPath = '';

        html += `<a href="/" class="breadcrumb-link">Главная</a>`;

        parts.forEach((part, index) => {
            currentPath += '/' + part;

            html += `<span class="breadcrumb-sep">/</span>`;

            if (index === parts.length - 1) {
                html += `<span class="breadcrumb-current">${decodeURIComponent(part)}</span>`;
            } else {
                html += `<a href="${currentPath}" class="breadcrumb-link">${decodeURIComponent(part)}</a>`;
            }
        });

        container.innerHTML = html;
    }
})();