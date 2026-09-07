(function () {
    document.addEventListener('DOMContentLoaded', function () {
        const burgerBtn = document.getElementById('burgerBtn');
        const headerNav = document.getElementById('headerNav');

        if (!burgerBtn || !headerNav) return;

        const overlay = document.createElement('div');
        overlay.className = 'overlay';
        document.body.appendChild(overlay);

        burgerBtn.addEventListener('click', function (e) {
            e.stopPropagation();
            toggleMenu();
        });

        overlay.addEventListener('click', function () {
            closeMenu();
        });

        headerNav.addEventListener('click', function (e) {
            if (e.target.tagName === 'A' || e.target.tagName === 'BUTTON') {
                closeMenu();
            }
        });

        document.addEventListener('keydown', function (e) {
            if (e.key === 'Escape') {
                closeMenu();
            }
        });

        function toggleMenu() {
            burgerBtn.classList.toggle('active');
            headerNav.classList.toggle('active');
            overlay.classList.toggle('active');
        }

        function closeMenu() {
            burgerBtn.classList.remove('active');
            headerNav.classList.remove('active');
            overlay.classList.remove('active');
        }

        const baseUrl = window.location.origin;
        const styles = [`${baseUrl}/styles/menu.css`];

        styles.forEach(stylePath => {
            const link = document.createElement('link');
            link.rel = 'stylesheet';
            link.href = stylePath;
            document.head.appendChild(link);
        });
    });
})();