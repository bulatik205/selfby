(function () {
    const STORAGE_KEY = 'theme';
    const BW = 'bw';

    function applyTheme(theme) {
        if (theme === BW) {
            document.documentElement.dataset.theme = BW;
        } else {
            delete document.documentElement.dataset.theme;
        }
    }

    function initTheme() {
        const saved = localStorage.getItem(STORAGE_KEY);
        applyTheme(saved === BW ? BW : 'default');
    }

    window.toggleTheme = function () {
        const isBw = document.documentElement.dataset.theme === BW;
        const next = isBw ? 'default' : BW;
        applyTheme(next);
        localStorage.setItem(STORAGE_KEY, next);
        updateThemeButtons(next);
    };

    function updateThemeButtons(theme) {
        document.querySelectorAll('[data-theme-toggle]').forEach(btn => {
            btn.textContent = theme === BW ? '☀' : '🌙';
            btn.setAttribute('aria-label', theme === BW ? 'Светлая тема' : 'Чёрно-белая тема');
        });
    }

    function bindButtons() {
        document.querySelectorAll('[data-theme-toggle]').forEach(btn => {
            btn.addEventListener('click', window.toggleTheme);
        });
        updateThemeButtons(document.documentElement.dataset.theme === BW ? BW : 'default');
    }

    initTheme();

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', bindButtons);
    } else {
        bindButtons();
    }
})();