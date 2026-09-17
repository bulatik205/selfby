(function () {
    function init() {
        const link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css';
        document.head.appendChild(link);

        const script = document.createElement('script');
        script.src = 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js';
        script.async = false;

        script.onload = function () {
            function highlightAll(root) {
                const scope = root || document;
                scope.querySelectorAll('pre code:not(.hljs)').forEach(block => {
                    hljs.highlightElement(block);
                });
            }

            highlightAll();

            const observer = new MutationObserver(mutations => {
                for (const m of mutations) {
                    for (const node of m.addedNodes) {
                        if (node.nodeType !== 1) continue;

                        if (node.matches && node.matches('pre code, pre')) {
                            highlightAll(node.parentElement || node);
                        }
                        else if (node.querySelectorAll) {
                            const codes = node.querySelectorAll('pre code:not(.hljs)');
                            if (codes.length) highlightAll(node);
                        }
                    }
                }
            });

            observer.observe(document.body, {
                childList: true,
                subtree: true,
            });
        };

        document.head.appendChild(script);
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', init);
    } else {
        init();
    }
})();