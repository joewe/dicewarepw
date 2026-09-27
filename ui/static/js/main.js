// Navigation active state
document.addEventListener("DOMContentLoaded", () => {
    const navLinks = document.querySelectorAll("nav a");
    for (let i = 0; i < navLinks.length; i++) {
        const link = navLinks[i];
        if (link.getAttribute('href') === window.location.pathname) {
            link.classList.add("live");
            break;
        }
    }

    // Initialize entropy meter bar
    initEntropyMeter();
});

// Copy to clipboard with modern feedback
function copyToClipboard() {
    const p = document.getElementById('passphrase');
    if (!p) return;

    const text = p.innerText.trim();
    if (!text) return;

    const btn = document.getElementById('btn-copy') || document.querySelector('button[onclick="copyToClipboard()"]');

    navigator.clipboard.writeText(text).then(() => {
        if (!btn) return;
        
        btn.classList.add('copied');
        const textSpan = btn.querySelector('.action-text');
        const originalText = textSpan ? textSpan.innerText : btn.innerText;

        if (textSpan) {
            textSpan.innerText = "Kopiert!";
        } else {
            btn.innerText = "Kopiert!";
        }

        // Pulse the passphrase box
        const display = p.closest('.passphrase-display');
        if (display) {
            display.classList.add('copy-flash');
            setTimeout(() => display.classList.remove('copy-flash'), 600);
        }

        setTimeout(() => {
            btn.classList.remove('copied');
            if (textSpan) {
                textSpan.innerText = originalText;
            } else {
                btn.innerText = originalText;
            }
        }, 2200);
    }).catch(err => {
        console.error('Kopieren fehlgeschlagen: ', err);
    });
}

// Add or remove special character "Spice"
function addSpice() {
    const p = document.getElementById('passphrase');
    const btn = document.getElementById('btn-spice') || document.querySelector('button[onclick="addSpice()"]');
    if (!p) return;

    // Check if already spiced
    if (p.dataset.spiced === "true") {
        // Revert to original
        p.innerText = p.dataset.original || p.innerText;
        p.dataset.spiced = "false";
        
        if (btn) {
            btn.classList.remove('active');
            const textSpan = btn.querySelector('.action-text');
            if (textSpan) {
                textSpan.innerText = "+ Sonderzeichen";
            } else {
                btn.innerText = "+ 🌶️";
            }
        }
        return;
    }

    // Save original text if not saved yet
    if (!p.dataset.original) {
        p.dataset.original = p.innerText.trim();
    }

    // Add spice
    const specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?";
    const char = specialChars[Math.floor(Math.random() * specialChars.length)];
    let text = p.innerText.trim();

    if (Math.random() < 0.5) {
        text = char + text;
    } else {
        text = text + char;
    }

    p.innerText = text;
    p.dataset.spiced = "true";

    if (btn) {
        btn.classList.add('active');
        const textSpan = btn.querySelector('.action-text');
        if (textSpan) {
            textSpan.innerText = "- Sonderzeichen";
        } else {
            btn.innerText = "- 🌶️";
        }
    }
}

// Compute & render dynamic entropy meter bar
function initEntropyMeter() {
    const pwbits = document.getElementById('pwbits');
    const meterBar = document.getElementById('entropy-meter-bar');
    if (!pwbits || !meterBar) return;

    const text = pwbits.innerText;
    const match = text.match(/([0-9]+(\.[0-9]+)?)/);

    if (match) {
        const bits = parseFloat(match[1]);
        // 64.5 bits = ~60%, 77.4 bits = ~75%, 90.3 bits = ~88%, 103.2 bits = 100%
        let percent = Math.min(100, Math.max(20, (bits / 103.2) * 100));

        // Smooth width animation
        requestAnimationFrame(() => {
            meterBar.style.width = percent + '%';
        });

        if (bits >= 90) {
            meterBar.style.background = 'linear-gradient(90deg, #10b981 0%, #06b6d4 100%)';
            meterBar.style.boxShadow = '0 0 12px rgba(16, 185, 129, 0.5)';
        } else if (bits >= 75) {
            meterBar.style.background = 'linear-gradient(90deg, #6366f1 0%, #10b981 100%)';
            meterBar.style.boxShadow = '0 0 12px rgba(99, 102, 241, 0.4)';
        } else {
            meterBar.style.background = 'linear-gradient(90deg, #818cf8 0%, #6366f1 100%)';
            meterBar.style.boxShadow = '0 0 10px rgba(129, 140, 248, 0.3)';
        }
    }
}