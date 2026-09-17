(function () {
    document.addEventListener('DOMContentLoaded', function () {
        const baseUrl = window.location.origin;
        const styles = [
            'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css'
        ];

        styles.forEach(stylePath => {
            const link = document.createElement('link');
            link.rel = 'stylesheet';
            link.href = stylePath;
            document.head.appendChild(link);
        });

        const script = document.createElement('script');
        script.src = 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js';
        script.onload = function () {
            document.querySelectorAll('pre code').forEach(block => {
                hljs.highlightElement(block);
            });
        };
        document.head.appendChild(script);
    });
})();