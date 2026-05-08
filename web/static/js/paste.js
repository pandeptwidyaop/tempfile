(function () {
    const tabs = document.querySelectorAll('[data-tab-target]');
    if (tabs.length > 0) {
        tabs.forEach((btn) => {
            btn.addEventListener('click', () => {
                const target = btn.getAttribute('data-tab-target');
                document.querySelectorAll('[data-tab-target]').forEach((b) =>
                    b.classList.toggle('active', b === btn)
                );
                document.querySelectorAll('[data-tab-panel]').forEach((p) =>
                    p.classList.toggle('active', p.getAttribute('data-tab-panel') === target)
                );
            });
        });
    }

    const pasteContent = document.getElementById('pasteContent');
    if (pasteContent) {
        const code = pasteContent.querySelector('code');
        if (code) {
            const text = code.textContent;
            const lines = text.split('\n');
            if (lines.length > 1 && lines[lines.length - 1] === '') {
                lines.pop();
            }
            code.innerHTML = '';
            lines.forEach((line) => {
                const span = document.createElement('span');
                span.className = 'paste-line';
                span.textContent = line.length === 0 ? ' ' : line;
                code.appendChild(span);
            });
        }
    }

    const copyBtn = document.getElementById('pasteCopyBtn');
    if (copyBtn) {
        copyBtn.addEventListener('click', async () => {
            const targetId = copyBtn.getAttribute('data-content-id');
            const target = document.getElementById(targetId);
            if (!target) return;
            const text = (target.querySelector('code') || target).innerText;
            try {
                await navigator.clipboard.writeText(text);
                const original = copyBtn.textContent;
                copyBtn.textContent = 'Copied!';
                copyBtn.disabled = true;
                setTimeout(() => {
                    copyBtn.textContent = original;
                    copyBtn.disabled = false;
                }, 1200);
            } catch (e) {
                copyBtn.textContent = 'Copy failed';
            }
        });
    }
})();
