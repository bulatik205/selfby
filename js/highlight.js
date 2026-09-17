(function () {
    function wrapCodeBlocks(root) {
        const scope = root || document;
        scope.querySelectorAll('pre > code').forEach(code => {
            const pre = code.parentElement;
            if (pre.parentElement && pre.parentElement.classList.contains('code-block')) return;

            const langClass = [...code.classList].find(c => c.startsWith('language-'));
            const lang = langClass ? langClass.replace('language-', '') : 'text';

            const wrapper = document.createElement('div');
            wrapper.className = 'code-block';

            const header = document.createElement('div');
            header.className = 'code-header';

            const langSpan = document.createElement('span');
            langSpan.className = 'code-lang';
            langSpan.textContent = lang;

            const btn = document.createElement('button');
            btn.type = 'button';
            btn.className = 'copy-btn';
            btn.textContent = 'Копировать';
            btn.addEventListener('click', () => {
                const text = code.innerText;
                const done = () => {
                    btn.textContent = 'Скопировано!';
                    setTimeout(() => (btn.textContent = 'Копировать'), 1500);
                };
                if (navigator.clipboard && window.isSecureContext) {
                    navigator.clipboard.writeText(text).then(done);
                } else {
                    const ta = document.createElement('textarea');
                    ta.value = text;
                    ta.style.position = 'fixed';
                    ta.style.opacity = '0';
                    document.body.appendChild(ta);
                    ta.select();
                    document.body.removeChild(ta);
                    done();
                }
            });

            header.append(langSpan, btn);
            pre.parentNode.insertBefore(wrapper, pre);
            wrapper.append(header, pre);
        });
    }

    function highlightAll(root) {
        const scope = root || document;
        scope.querySelectorAll('pre code:not(.hljs)').forEach(block => {
            hljs.highlightElement(block);
        });
    }

    function boot() {
        const link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css';
        document.head.appendChild(link);

        const script = document.createElement('script');
        script.src = 'https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js';

        script.onload = function () {
            highlightAll();
            wrapCodeBlocks();

            const observer = new MutationObserver(mutations => {
                for (const m of mutations) {
                    for (const node of m.addedNodes) {
                        if (node.nodeType !== 1) continue;

                        if (node.matches && node.matches('pre, pre code')) {
                            const scope = node.closest('pre')?.parentElement || node.parentElement || node;
                            highlightAll(scope);
                            wrapCodeBlocks(scope);
                        } else if (node.querySelectorAll) {
                            const codes = node.querySelectorAll('pre code:not(.hljs)');
                            if (codes.length) {
                                highlightAll(node);
                                wrapCodeBlocks(node);
                            }
                        }
                    }
                }
            });

            observer.observe(document.body, { childList: true, subtree: true });
        };

        script.onerror = function () {
            wrapCodeBlocks();
        };

        document.head.appendChild(script);
    }

    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', boot);
    } else {
        boot();
    }
})();