var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

function copyToClipboard() {
	const text = document.getElementById('passphrase').innerText;
	navigator.clipboard.writeText(text).then(function () {
		const btn = document.querySelector('button[onclick="copyToClipboard()"]');
		const originalText = btn.innerText;
		btn.innerText = "Copied!";
		setTimeout(() => {
			btn.innerText = originalText;
		}, 2000);
	}, function (err) {
		console.error('Could not copy text: ', err);
	});
}

function addSpice() {
	const p = document.getElementById('passphrase');
	const btn = document.querySelector('button[onclick="addSpice()"]');

	// Check if already spiced
	if (p.dataset.spiced === "true") {
		// Revert to original
		p.innerText = p.dataset.original;
		p.dataset.spiced = "false";
		btn.innerText = "+ 🌶️";
		return;
	}

	// Save original text if not saved yet
	if (!p.dataset.original) {
		p.dataset.original = p.innerText;
	}

	// Add spice
	const specialChars = "!@#$%^&*()_+-=[]{}|;:,.<>?";
	const char = specialChars[Math.floor(Math.random() * specialChars.length)];
	let text = p.innerText;

	if (Math.random() < 0.5) {
		text = char + text;
	} else {
		text = text + char;
	}

	p.innerText = text;
	p.dataset.spiced = "true";
	btn.innerText = "- 🌶️";
}