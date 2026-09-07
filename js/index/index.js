document.addEventListener('DOMContentLoaded', () => {
    initTypewriter();
});

function initTypewriter() {
    const element = document.getElementById('typewriter');
    if (!element) return;
    
    const words = ['SelfBy', 'Проекты', '# My project'];
    let wordIndex = 0;
    let charIndex = 0;
    let isDeleting = false;
    
    function type() {
        const currentWord = words[wordIndex];
        
        if (!isDeleting) {
            if (charIndex < currentWord.length) {
                element.textContent = currentWord.substring(0, charIndex + 1);
                charIndex++;
                const delay = Math.random() * 100 + 100;
                setTimeout(type, delay);
            } else {
                isDeleting = true;
                setTimeout(type, 2500);
            }
        } else {
            if (charIndex > 0) {
                element.textContent = currentWord.substring(0, charIndex - 1);
                charIndex--;
                setTimeout(type, 80);
            } else {
                isDeleting = false;
                wordIndex = (wordIndex + 1) % words.length;
                setTimeout(type, 700);
            }
        }
    }
    
    type();
}