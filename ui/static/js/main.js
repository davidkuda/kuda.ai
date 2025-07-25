console.log("Hello from the go server");

// ------------------------------------------------------------
// DOM references
// ------------------------------------------------------------
const themeBtn = document.getElementById("themeToggle");

// ------------------------------------------------------------
// event listeners
// ------------------------------------------------------------
themeBtn.addEventListener("click", async () => {
	const cur = document.documentElement.getAttribute("data-theme") || "dark";
	const next = cur === "dark" ? "light" : "dark";
	applyTheme(next);
	await localStorage.setItem("theme", next);
});

// ------------------------------------------------------------
// functions
// ------------------------------------------------------------
function applyTheme(mode) {
	document.documentElement.setAttribute("data-theme", mode);
	themeBtn.textContent = mode === "dark" ? "light mode" : "dark mode";
}

// ------------------------------------------------------------
// init
// ------------------------------------------------------------
(async () => {
	const theme = await localStorage.getItem("theme");
	applyTheme(theme);
})();
