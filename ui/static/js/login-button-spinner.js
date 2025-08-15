const form = document.getElementById("login-form");
const button = document.getElementById("login-button");

form.addEventListener("submit", () => {
	button.setAttribute("aria-busy", "true");
	button.disabled = true;
	button.classList.add("loading");
});

// Listen for HTMX after swap (response received & DOM updated)
form.addEventListener("htmx:afterSwap", (event) => {
	// Check if an error span got updated (meaning login failed)
	if (form.querySelector("span.error")?.textContent.trim()) {
		button.removeAttribute("aria-busy");
		button.disabled = false;
		button.classList.remove("loading");
	}
});
