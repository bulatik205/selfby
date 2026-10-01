(function () {
    document.addEventListener('DOMContentLoaded', () => {
        document.querySelectorAll('[data-tabs]').forEach(initTabs);
    });

    function initTabs(container) {
        container.addEventListener('click', e => {
            const btn = e.target.closest('.work-tab');
            if (!btn || !container.contains(btn)) return;

            const targetId = btn.getAttribute('data-tab-target');
            if (!targetId) return;

            switchTab(container, btn, targetId);
        });
    }

    function switchTab(container, activeBtn, targetId) {
        container.querySelectorAll('.work-tab').forEach(t => {
            t.classList.toggle('work-tab--active', t === activeBtn);
            t.setAttribute('aria-selected', t === activeBtn ? 'true' : 'false');
        });

        const card = container.closest('.work') || container.parentElement;
        card.querySelectorAll('.work-content').forEach(c => {
            c.hidden = c.id !== targetId;
        });
    }
})();