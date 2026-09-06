(function() {
    if (!getCookie('tips_work_closed')) {
        const tipsBlock = document.getElementById('tipsBlock');
        if (tipsBlock) {
            tipsBlock.style.display = 'flex';
        }
    }
})();

function closeTips(page) {
    const expires = new Date();
    expires.setTime(expires.getTime() + (30 * 24 * 60 * 60 * 1000));
    document.cookie = `tips_${page}_closed=true; expires=${expires.toUTCString()}; path=/`;
    
    const tipsBlock = document.getElementById('tipsBlock');
    if (tipsBlock) {
        tipsBlock.style.display = 'none';
    }
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