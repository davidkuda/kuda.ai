const form = document.getElementById("login-form");
const button = document.getElementById("login-button");

form.addEventListener("submit", () => {
	button.setAttribute("aria-busy", "true");
	button.disabled = true;
	button.classList.add("loading");
});
