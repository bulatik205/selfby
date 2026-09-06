function initTips(pageName) {
    if (!getCookie(`tips_${pageName}_closed`)) {
        const tipsBlock = document.getElementById('tipsBlock');
        if (tipsBlock) {
            tipsBlock.style.display = 'flex';
        }
    }

    window.closeTips = function () {
        const expires = new Date();
        expires.setTime(expires.getTime() + (30 * 24 * 60 * 60 * 1000));
        document.cookie = `tips_${pageName}_closed=true; expires=${expires.toUTCString()}; path=/`;

        const tipsBlock = document.getElementById('tipsBlock');
        if (tipsBlock) {
            tipsBlock.style.display = 'none';
        }
    };
}

function getCookie(name) {
    const nameEQ = name + "=";
    const ca = document.cookie.split(';');
    for (let i = 0; i < ca.length; i++) {
        let c = ca[i];
        while (c.charAt(0) === ' ') c = c.substring(1, c.length);
        if (c.indexOf(nameEQ) === 0) return c.substring(nameEQ.length, c.length);
    }
    return null;
}

document.addEventListener('DOMContentLoaded', function () {
    if (typeof pageName !== 'undefined') {
        initTips(pageName);
    }

    const baseUrl = window.location.origin;

    const styles = [`${baseUrl}/styles/tips.css`];

    styles.forEach(stylePath => {
        const link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = stylePath;
        document.head.appendChild(link);
    });
});